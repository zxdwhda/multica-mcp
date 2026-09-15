package oauth

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"os"
)

var ErrMissing = errors.New("record not found")
var ErrUsed = errors.New("record already used")

type Store interface {
	Put(context.Context, string, any) error
	Get(context.Context, string, any) error
	Claim(context.Context, string) error
}
type OSSStore struct {
	Endpoint, Bucket string
	aead             cipher.AEAD
}

func NewOSSStore(endpoint, bucket, secret string) (*OSSStore, error) {
	key := sha256.Sum256([]byte("multica-mcp OAuth storage v1:" + secret))
	b, e := aes.NewCipher(key[:])
	if e != nil {
		return nil, e
	}
	a, e := cipher.NewGCM(b)
	return &OSSStore{endpoint, bucket, a}, e
}
func (s *OSSStore) bucket() (*oss.Bucket, error) {
	id, secret, token := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN")
	if id == "" || secret == "" {
		return nil, fmt.Errorf("FC role credentials unavailable")
	}
	c, e := oss.New(s.Endpoint, id, secret, oss.SecurityToken(token), oss.Timeout(5, 15))
	if e != nil {
		return nil, e
	}
	return c.Bucket(s.Bucket)
}
func hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func (s *OSSStore) Put(ctx context.Context, key string, value any) error {
	b, e := s.bucket()
	if e != nil {
		return e
	}
	raw, e := json.Marshal(value)
	if e != nil {
		return e
	}
	n := make([]byte, s.aead.NonceSize())
	if _, e = rand.Read(n); e != nil {
		return e
	}
	encrypted := s.aead.Seal(n, n, raw, []byte(key))
	return b.PutObject("multica/"+key, bytes.NewReader(encrypted), oss.ForbidOverWrite(true), oss.WithContext(ctx))
}
func (s *OSSStore) Get(ctx context.Context, key string, value any) error {
	b, e := s.bucket()
	if e != nil {
		return e
	}
	r, e := b.GetObject("multica/"+key, oss.WithContext(ctx))
	if e != nil {
		var se oss.ServiceError
		if errors.As(e, &se) && se.StatusCode == 404 {
			return ErrMissing
		}
		return e
	}
	defer r.Close()
	raw, e := io.ReadAll(io.LimitReader(r, 128*1024))
	if e != nil {
		return e
	}
	n := s.aead.NonceSize()
	if len(raw) < n {
		return errors.New("invalid encrypted state")
	}
	raw, e = s.aead.Open(nil, raw[:n], raw[n:], []byte(key))
	if e != nil {
		return e
	}
	return json.Unmarshal(raw, value)
}
func (s *OSSStore) Claim(ctx context.Context, key string) error {
	e := s.Put(ctx, "claims/"+key, map[string]bool{"used": true})
	var se oss.ServiceError
	if errors.As(e, &se) && (se.StatusCode == 409 || se.StatusCode == 412) {
		return ErrUsed
	}
	return e
}
