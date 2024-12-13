package middleware

import (
	"bcca_crawler/internal/auth"
	"bcca_crawler/internal/config"	
	"bcca_crawler/internal/caching"
	"bcca_crawler/internal/auth/roles"
	"context"
	"fmt"
	"net/http"
	"time"	
)


// func middlewareLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

//Fix middlewareAuth, middlewarePermission : it shouldn't deny access, merely provide context if user is authenticated or not and their role

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
				fmt.Println("Middleware : Should have set User: ", ctx)
				role, err := caching.GetRoleCache(user_token.UserID)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Error getting role from cache.. This should not happen")
				}
				newCtx := context.WithValue(ctx, auth.UserRoleKey, role)
				fmt.Println("Middleware : Should have set Role and User: ", newCtx)

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
		fmt.Println("Middleware : Should have set Role and User: ", ctx)
		Newctx := context.WithValue(ctx, auth.UserRoleKey, role)
		next.ServeHTTP(w, r.WithContext(Newctx))
	})
}

