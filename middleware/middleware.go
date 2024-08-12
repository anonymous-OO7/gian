package middleware

import (
	"encoding/json"
	"fmt"
	"gian/utils"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(secretKey string, next http.HandlerFunc, isClientAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split the Authorization header to get the token part
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}
		// Get the token string
		tokenString := splitToken[1]
		var err error
		if isClientAuth {
			encryptionKey := []byte("thisisaverysecureencryptionkey12")
			tokenString, err = utils.DecryptAES(encryptionKey, tokenString)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Token error!! Malformed", "reason": err.Error()})
				return
			}
			// Update the Authorization header with the decrypted token
			r.Header.Set("Authorization", "Bearer "+tokenString)
		}
		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid signing method")
			}
			return []byte(secretKey), nil
		})
		// Check for parsing errors
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		excludedPathsClient := []string{"/get-clientinfo"}
		isExcluded := false
		for _, path := range excludedPathsClient {
			if r.URL.Path == path {
				isExcluded = true
				break
			}
		}

		// Check if the token is valid and not expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
			if time.Now().After(expirationTime) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			// Check if the URL param is equal to the token client code if the request comes from the client
			if isClientAuth && !isExcluded {
				clientCode := chi.URLParam(r, "code")
				if clientCode == "" {
					clientCode = r.Header.Get("code")
				}
				tokenClientCode, ok := claims["client_code"].(string)
				if !ok || tokenClientCode != clientCode {
					http.Error(w, "Unauthorized: client code mismatch with token code", http.StatusUnauthorized)
					return
				} else {
					r.Header.Set("Companycode", tokenClientCode)
				}
			}

			// Call the next handler in the chain if authorization is successful
			next(w, r)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	}
}

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// List of paths that don't require token validation
		excludedPaths := []string{"/login", "/signup", "/device", "/update-password", "/master-update-password", "/request-otp", "/verify-otp"}

		// Check if the requested path is in the excludedPaths
		isExcluded := false
		for _, path := range excludedPaths {
			if r.URL.Path == path {
				isExcluded = true
				break
			}
		}

		// If the route is excluded or is for "connect_bench", call the next handler without token validation
		if isExcluded {
			next(w, r)
		} else {
			// Otherwise, validate the token using AuthMiddleware with your secret key
			AuthMiddleware("your-secret-key", next, false)(w, r)
		}
	}
}
