package auth

import (
	"context"
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

		userInfo, err := a.tokensManager.VerifyJWTToken(bearer)
		if err != nil {
			web.WriteUnauthorized(w, err)
			return
		}

		ctxUserID := context.WithValue(context.Background(), "userID", userInfo.UserID)
		ctxUserRoleID := context.WithValue(ctxUserID, "userRoleID", userInfo.UserRoleID)

		next.ServeHTTP(w, r.WithContext(ctxUserRoleID))
	})
}

func (a *AuthService) AuthOnlyAdminMiddleware(next http.Handler) http.Handler {
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
			web.WriteForbidden(w, nil)
			return
		}

		ctxUserID := context.WithValue(context.Background(), "userID", userInfo.UserID)
		next.ServeHTTP(w, r.WithContext(ctxUserID))
	})
}
