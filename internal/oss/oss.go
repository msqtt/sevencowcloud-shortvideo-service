package oss

import (
	"context"
	"fmt"
	"io"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Kodo is used to upload file to qiniu oss.
type Kodo struct {
	conf      *storage.Config
	accesskey string
	secretKey string
}

func NewKodo(https, cdn bool, ak, sk string) *Kodo {
	conf := &storage.Config{
		Region:        &storage.ZoneHuanan,
		UseHTTPS:      https,
		UseCdnDomains: cdn,
	}
	return &Kodo{
		conf:      conf,
		accesskey: ak,
		secretKey: sk,
	}
}

// UploadDataByForm uses form to upload little file from bytes data.
func (kd *Kodo) UploadDataByForm(ctx context.Context,
	bucket, fileName string, data io.Reader) (*storage.PutRet, error) {
	pp := storage.PutPolicy{Scope: fmt.Sprintf("%s:%s", bucket, fileName)}
	mac := qbox.NewMac(kd.accesskey, kd.secretKey)
	token := pp.UploadToken(mac)
	fUploader := storage.NewFormUploader(kd.conf)
	// return format
	ret := new(storage.PutRet)

	err := fUploader.Put(ctx, ret, token, fileName, data, -1, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot upload data: %w", err)
	}
	return ret, nil
}
