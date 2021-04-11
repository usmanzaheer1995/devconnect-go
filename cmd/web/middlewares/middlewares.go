package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"net/http"
	"os"
	"strings"
)

// HTTP middleware setting a value on the request context
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		fmt.Println("Inside middleware")
		ctx := context.WithValue(r.Context(), "user", "123")

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SecondMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		fmt.Println(user)
		next.ServeHTTP(w, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header string
		if r.Header.Get("Authorization") != "" {
			header = r.Header.Get("Authorization")
		} else if r.Header.Get("x-auth-token") != "" {
			header = r.Header.Get("x-auth-token")
		} else {
			utils.ERROR(w, http.StatusForbidden, errors.New("invalid token"), nil)
			return
		}
		header = strings.TrimSpace(header)

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			utils.ERROR(w, http.StatusForbidden, errors.New("invalid token"), nil)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "userID", claims["userID"]) // adding the user ID to the context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
