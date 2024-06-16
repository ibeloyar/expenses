package auth

import (
	"fmt"
	"net/http"

	"github.com/B-Dmitriy/expenses/pgk/web"
)

func (a *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, err := a.passManager.GetAuthorizationHeader(r)
		if err != nil {
			web.WriteUnauthorized(w, err)
			return
		}

		_, err = a.tokensManager.VerifyJWTToken(bearer)
		if err != nil {
			web.WriteUnauthorized(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *AuthService) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, err := a.passManager.GetAuthorizationHeader(r)
		if err != nil {
			web.WriteUnauthorized(w, err)
			return
		}

		userInfo, err := a.tokensManager.VerifyJWTToken(bearer)
		if err != nil {
			web.WriteUnauthorized(w, err)
			return
		}

		if userInfo.UserRoleID != 1 {
			web.WriteForbidden(w, fmt.Errorf("forbidden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
