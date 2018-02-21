package chi_middlewares

import (
	"context"
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"github.com/akrylysov/algnhsa"
)

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

func AWSRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var reqID string

		awsRequest, ok := algnhsa.ProxyRequestFromContext(r.Context())

		if !ok {
			uuidID, err := uuid.NewV4()
			if err != nil {
				rand.Seed(time.Now().UnixNano())
				randId := rand.Int63()
				reqID = base64.URLEncoding.EncodeToString([]byte(strconv.FormatInt(randId, 10)))
			} else {
				reqID = uuidID.String()
			}
		} else {
			reqID = awsRequest.RequestContext.RequestID
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}

	return ""
}
