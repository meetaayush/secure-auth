package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := extractTokenFromCookie(r)
		if err != nil {
			app.unauthorizedError(w, r, errors.New("invalid user"))
			return
		}

		jwtToken, err := app.store.Auth.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, errors.New("invalid token"))
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		// get session ID
		sessionId, err := uuid.Parse(claims["sid"].(string))
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		session, err := app.store.Session.GetSession(ctx, sessionId)
		if err != nil {
			app.unauthorizedError(w, r, errors.New("invalid session"))
			return
		}

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			app.unauthorizedError(w, r, errors.New("invalid user"))
			return
		}

		ctx1 := context.WithValue(ctx, userCtx, user)
		ctx2 := context.WithValue(ctx1, sessionCtx, session)

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
