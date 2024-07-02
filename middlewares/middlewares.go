package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type contextKey string

const UserIDKey contextKey = "userID"

func CreateToken(data string) (string, error) {
	claim := jwt.MapClaims{
		"sub": data,
		"iss": "taskscheduler",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := tokenClaim.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiredAt, _ := claims["exp"].(float64)
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("expiredAt claim is missing or not a float64")
		}

		if float64(time.Now().Unix()) > expiredAt {
			return "", errors.New("token is expired")
		}
		return userID, nil
	} else {
		return "", errors.New("invalid token")
	}

	// return nil
}

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("x-access-token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Status unauthorized", http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		token := cookie.Value
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := VerifyToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
