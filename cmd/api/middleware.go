package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"socialone/internal/store"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("No authorization header")
			app.unauthorisedErrorResponse(w, r, errors.New("no authorization header"))
			return
		}

		// Parse the auth header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Invalid authorization header")
			app.unauthorisedErrorResponse(w, r, errors.New("invalid authorization header"))
			return
		}

		// Validate the token
		token := parts[1]
		jwtTokens, err := app.authenticator.ValidateToken(token)
		if err != nil {
			log.Printf("Error parsing token: %v\n", err)
			app.unauthorisedErrorResponse(w, r, fmt.Errorf("error parsing token: %w", err))
			return
		}
		claims, _ := jwtTokens.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			log.Printf("Error parsing user ID: %v\n", err)
			app.unauthorisedErrorResponse(w, r, fmt.Errorf("error parsing user ID: %w", err))
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			log.Printf("Error getting user by ID: %v\n", err)
			app.unauthorisedErrorResponse(w, r, err)
			return
		}
		// Add the user to the context
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("No authorization header")
				app.unauthorisedBasicErrorResponse(w, r, errors.New("no authorization header"))
				return
			}

			// Parse the auth header
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				log.Println("Invalid authorization header")
				app.unauthorisedBasicErrorResponse(w, r, errors.New("invalid authorization header"))
				return
			}

			// Decode the credentials
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				log.Printf("Error decoding credentials: %v\n", err)
				app.unauthorisedErrorResponse(w, r, err)
				return
			}
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				log.Println("Invalid credentials")
				app.unauthorisedBasicErrorResponse(w, r, errors.New("invalid credentials"))
				return
			}

			// Check the credentials
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromContext(r.Context())
		if post.UserId == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			log.Printf("Error checking role precedence: %v\n", err)
			app.internalServerError(w, r, err)
			return
		}
		if !allowed {
			log.Println("User not allowed")
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		log.Printf("Error getting role by name: %v\n", err)
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
