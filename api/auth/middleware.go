package auth

import (
	"context"
	"memo/pkg/response"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Mux   *http.ServeMux
	Store AuthStore
}

func (m *AuthMiddleware) Handle(pattern string, handler http.HandlerFunc) {
	m.Mux.Handle(pattern, m.use(handler))
}

func (m *AuthMiddleware) use(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		parts := strings.Split(auth, " ")
		if len(parts) < 2 {
			response.RespondErr(w, response.Unauthorized())
			return
		}

		token := parts[1]
		if userId, err := m.Store.FindToken(token, r.Context()); err != nil {
			response.RespondErr(w, response.Unauthorized())
		} else {
			r = r.WithContext(context.WithValue(r.Context(), "user", userId))

			handler(w, r)
		}
	}
}
