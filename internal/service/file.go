package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"

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

func NewFileService(conf config.Config, token token.TokenMaker, store db.Store,
kodo *oss.Kodo) *FileServer {
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

func (fs *FileServer) UploadVideo(w http.ResponseWriter,
	r *http.Request, pathParams map[string]string) {
	fs.authHandleWrap(w, r, pathParams,
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			ctx := r.Context()
			payload := ctx.Value("payload").(*token.Payload)

			u, err := fs.store.GetUserByID(ctx, payload.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(BadResultMap(codes.Unauthenticated, "user not found"))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to find user"))
				return
			}

			var allSize int64 = fs.conf.VideoLimit << 20
			var coverSize int64 = fs.conf.ImageLimit << 20

			// 2G limit
			r.Body = http.MaxBytesReader(w, r.Body, allSize)

			// get multipart file
			err = r.ParseMultipartForm(allSize)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				maxErr := &http.MaxBytesError{}
				if errors.As(err, &maxErr) {
					w.Write(BadResultMap(codes.OutOfRange, "total file size exceeds limit"))
					return
				}
				w.Write(BadResultMap(codes.InvalidArgument, "cannot parse multipart  file"))
				return
			}

			fh := r.MultipartForm.File["cover"]
			var coverKey string
			if len(fh) > 0 {
				if fh[0].Size > coverSize {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(BadResultMap(codes.OutOfRange, "cover size exceeds limit"))
					return
				}
				coverFile, err := fh[0].Open()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(BadResultMap(codes.Internal, "failed to read cover file"))
					return
				}
				defer coverFile.Close()

				coverKey, err = fs.UploadImg(ctx, fmt.Sprintf("tmp/cover/%d-", u.ID), coverFile)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(BadResultMap(codes.Internal, "failed to upload cover file"))
					return
				}
			}

			fh2 := r.MultipartForm.File["video"]
			if len(fh2) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(BadResultMap(codes.InvalidArgument, "no file uploaded"))
				return
			}

			f, err2 := fh2[0].Open()
			if err2 != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to open video file"))
				return
			}
			defer f.Close()

			headByte := make([]byte, 512)
			if _, err4 := f.Read(headByte); err4 != nil {
				log.Println(err4)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to read video file"))
				return
			}
			exts, err3 := mime.ExtensionsByType(http.DetectContentType(headByte))
			if err3  != nil{
				log.Println(err3)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to get extension"))
			}
			// todo
			extension := exts[1]

			nanoid, err3 := gonanoid.New()
			if err3 != nil {
				log.Println(err3)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to read file"))
				return
			}

			_, err4 := f.Seek(0, 0)
			if err4 != nil {
				log.Println(err4)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to seek video file"))
				return
			}

			key := fmt.Sprintf("tmp/video/%d-%s%s",u.ID, nanoid, extension)

			// log.Println(fs.conf.KodoBucket, key, extension, f)
			// upload video file shit
			_, err = fs.kodo.UploadDataByForm(ctx, fs.conf.KodoBucket, key, f)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to upload video file"))
				return
			}

			if len(coverKey) == 0 {
				coverKey = fmt.Sprintf("tmp/cover/%d-%s.jpg", u.ID, nanoid)
			}
			// insert video table and get video cover
			param := db.AddVideoParams{
				CoverLink:   coverKey,
				SrcLink:     key,
			}

			result, err := fs.store.AddVideo(ctx, param)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to add video"))
				return
			}
			id, err := result.LastInsertId()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to get video id"))
				return
			}

			clink := uriPath(fs.conf.KodoLink, coverKey)
			slink := uriPath(fs.conf.KodoLink, key)
			ret := make(map[string]any)
			ret["videoId"] = id
			ret["coverLink"] = clink
			ret["srcLink"] = slink
			ret["updatedAt"] = strconv.FormatInt(time.Now().Unix(), 10)
			b, _ := json.Marshal(ret)
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		},
	)
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
				w.Write(BadResultMap(codes.InvalidArgument, "wrong format of user id"))
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
			var size int64 = fs.conf.ImageLimit << 20

			r.Body = http.MaxBytesReader(w, r.Body, size)

			// get multipart file
			err = r.ParseMultipartForm(size)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				maxErr := &http.MaxBytesError{}
				if errors.As(err, &maxErr) {
					w.Write(BadResultMap(codes.OutOfRange, "image file size exceeds limit"))
					return
				}
				w.Write(BadResultMap(codes.InvalidArgument, "cannot parse multipart  file"))
				return
			}
			fh := r.MultipartForm.File["avatar"]
			if len(fh) == 0 {
				w.WriteHeader(http.StatusBadRequest)
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
			defer f.Close()

			key, err2 := fs.UploadImg(ctx, "img/avatar/", f)
			if err2 != nil {
				log.Println(err2)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(BadResultMap(codes.Internal, "failed to upload avatar"))
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

func (fs *FileServer) UploadImg(ctx context.Context, path string, file multipart.File) (string, error) {
	_, imgType, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return "", nil
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s%s.%s", path, nanoid, imgType)
	_, err2 := fs.kodo.UploadDataByForm(ctx, fs.conf.KodoBucket, key, file)
	if err2 != nil {
		return "", err
	}
	return key, nil
}
