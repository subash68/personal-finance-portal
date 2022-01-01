package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/spf13/viper"
)

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

var (
	prefix string

	reqID uint64
)

func init() {
	hostname, err := os.Hostname()

	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var buf [12]byte
	var b64 string

	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "+", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}

	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}

	return false
}

func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myid := atomic.AddUint64(&reqID, 1)
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, fmt.Sprintf("%s-%06d", prefix, myid))

		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}

	return ""
}
