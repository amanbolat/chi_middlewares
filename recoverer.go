package chi_middlewares

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"github.com/rs/zerolog"
)

func Recoverer(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {

					logEntry := logger.With().Logger()
					logEntry.Error().Timestamp().Interface("recover_info", rvr).Bytes("debug_stack", debug.Stack()).Msg("panic_on_request")

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
