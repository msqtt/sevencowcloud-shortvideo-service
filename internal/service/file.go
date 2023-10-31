package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/oss"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
)

// FileServer is used to handle http multipart file uploading.
type FileServer struct {
	conf  config.Config
	token token.TokenMaker
	kodo  *oss.Kodo
	store db.Store
}

func NewFileService(conf config.Config, token token.TokenMaker, store db.Store) *FileServer {
	kodo := oss.NewKodo(conf.KodoHttps, conf.KodoCDN, conf.QiniuAK, conf.QiniuSK)
	return &FileServer{
		conf:  conf,
		token: token,
		kodo:  kodo,
		store: store,
	}
}

// BadResultMap creates a uni grpc bad result json bytes.
func BadResultMap(code codes.Code, msg string) []byte {
	ret := make(map[string]any)
	ret["code"] = code
	ret["message"] = msg
	ret["details"] = []struct{}{}
	b, _ := json.Marshal(ret)
	return b
}

// authHandleWrap wraps handler with authorization.
func (fs *FileServer) authHandleWrap(w http.ResponseWriter,
	r *http.Request, pathParams map[string]string, handler runtime.HandlerFunc) {
	accessToken := r.Header["Authorization"]
	if len(accessToken) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(BadResultMap(codes.Unauthenticated, "access token is not provied"))
		return
	}
	p, err := fs.token.ValidToken(accessToken[0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(BadResultMap(codes.Unauthenticated,
			fmt.Sprintf("access token is invalid: %v", err)))
		return
	}

	r2 := r.WithContext(context.WithValue(r.Context(), "payload", p))
	handler(w, r2, pathParams)
}

// UploadAvatar overwrites user's avatar using received form multipart file.
func (fs *FileServer) UploadAvatar(w http.ResponseWriter,
	r *http.Request, pathParams map[string]string) {
	fs.authHandleWrap(w, r, pathParams,
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			id := pathParams["user_id"]
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(BadResultMap(codes.InvalidArgument, "wrong farmat of user id"))
				return
			}
			ctx := r.Context()
			payload, ok := ctx.Value("payload").(*token.Payload)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to convert payload"))
				return
			}
			if uid != payload.UserID {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(BadResultMap(codes.Unauthenticated, "cannot upload other user's avatar"))
				return
			}

			// find user
			usr, err := fs.store.GetUserByID(ctx, uid)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to find user"))
			}

			// limit body size 2MB
			var size int64 = 2 << 20
			http.MaxBytesReader(w, r.Body, size)
			// get multipart file
			err = r.ParseMultipartForm(size)
			if err != nil {
				w.WriteHeader(http.StatusNoContent)
				w.Write(BadResultMap(codes.InvalidArgument, "cannot parse multipart form"))
				return
			}
			fh := r.MultipartForm.File["file"]
			if len(fh) == 0 {
				w.WriteHeader(http.StatusNoContent)
				w.Write(BadResultMap(codes.InvalidArgument, "no file uploaded"))
				return
			}
			// get image type
			f, err := fh[0].Open()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to read file"))
				return
			}
			b2, _ := ioutil.ReadAll(f)
			buf1 := bytes.NewBuffer(b2)
			buf2 := bytes.NewBuffer(b2)

			_, imgType, err := image.DecodeConfig(buf1)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(BadResultMap(codes.InvalidArgument, "unknown image type"))
				return
			}

			nanoid, err := gonanoid.New()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to generate random id"))
				return
			}

			// save to oss
			key := fmt.Sprintf("img/avatar/%s.%s", nanoid, imgType)
			_, err2 := fs.kodo.UploadDataByForm(ctx, fs.conf.KodoBucket, key, buf2)
			if err2 != nil {
				log.Println(err2)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "cannot save file"))
				return
			}

			// update database
			args := db.UpdateAvatarParams{
				AvatarLink: sql.NullString{String: key, Valid: true},
				ID:         usr.ProfileID,
			}
			err = fs.store.UpdateAvatar(ctx, args)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to update avatar"))
			}

			retLink := make(map[string]string)
			retLink["link"] = filepath.Join(fs.conf.KodoLink, key)

			b, _ := json.Marshal(retLink)
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		})
}
