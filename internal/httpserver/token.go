package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type ctxKey string

func (s *server) validateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authString := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authString) != 2 {
			s.error(
				w,
				r,
				http.StatusUnauthorized,
				errors.New("invalid authorization header"),
			)
			return
		}
		token, err := jwt.Parse(
			authString[1],
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(
						"unexpected signing method: %v",
						token.Header["alg"],
					)
				}
				return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
			},
		)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if !claims["admin"].(bool) {
			s.error(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) getRefreshToken(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return token, nil
}
