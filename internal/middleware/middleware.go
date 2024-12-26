package middleware

import (
	"bcca_crawler/internal/auth"
	"bcca_crawler/internal/config"	
	"bcca_crawler/internal/caching"
	"bcca_crawler/internal/auth/roles"
	"bcca_crawler/internal/json_utils"
	"context"
	"fmt"
	"net/http"
	"time"	
)


func MiddlewareAuth(c *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth_cookie, err := auth.GetCookieToken(r, "auth_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		refresh_cookie, err := auth.GetCookieToken(r, "refresh_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return			
		}
		
		user_id, err := auth.ValidateJWT(auth_cookie, c.Secret)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "ValidateJWT Function: token has invalid claims: token is expired" {
				user_token, err := auth.ValidateRefreshToken(refresh_cookie, c)
				if err != nil {
					fmt.Println(err)
					next.ServeHTTP(w, r)
					return
				}
				auth.SetAuthCookies(w, user_token)
				ctx := context.WithValue(r.Context(), auth.UserIDKey, user_token.UserID)
				
				role, err := caching.GetRoleCache(user_token.UserID)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Error getting role from cache.. This should not happen")
				}
				newCtx := context.WithValue(ctx, auth.UserRoleKey, role)
				
				next.ServeHTTP(w, r.WithContext(newCtx))
				return

			} else {
				next.ServeHTTP(w, r)
				return
			}
		}
		ctx := context.WithValue(r.Context(), auth.UserIDKey, user_id)
		role, err := caching.GetRoleCache(user_id)
		if err != nil {
			role_string,err := c.Db.GetUserRoleByID(ctx, user_id)
			if err != nil {
				fmt.Println(err)
				next.ServeHTTP(w, r)
				return
			}
			role,err = roles.RoleFromString(role_string)
			if err != nil {
				fmt.Println(err)
				next.ServeHTTP(w, r)
				return
			}

		}
		caching.SetRoleCache(user_id, role, time.Now().Add(time.Minute * 60))		
		Newctx := context.WithValue(ctx, auth.UserRoleKey, role)
		next.ServeHTTP(w, r.WithContext(Newctx))
	})
}


func WithAuthAndRole(role roles.Role, handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := auth.GetUserFromContext(r)
        if err != nil {
            json_utils.RespondWithError(w, http.StatusForbidden, "Invalid token")
            return
        }		

        if user.Role >= role {
            json_utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to use this function.")
            return
        }
        handler(w, r) // Call the actual handler
    }
}

