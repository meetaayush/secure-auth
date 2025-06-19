package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/meetaayush/secure-auth/store"
)

type UserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload UserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.User{
		Email: payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Users.RegisterUser(ctx, &user)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			app.badRequestResponse(w, r, errors.New("user with this email already exists"))
		default:
			app.badRequestResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// parse the credentials
	var payload UserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// get user from database
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate the password
	err = user.Password.ComparePassword(payload.Password)
	if err != nil {
		app.unauthorizedError(w, r, errors.New("invalid email or password"))
		return
	}

	// generate the session and store in redis db
	ip := r.RemoteAddr
	ua := r.UserAgent()
	session, err := app.store.Session.CreateSession(ctx, user.Id, ip, ua)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// generate JWT token
	claims := jwt.MapClaims{
		"sub": user.Id,
		"sid": session.SessionId,
		"exp": time.Now().Add(app.config.AuthTokenConfig.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.AuthTokenConfig.iss,
		"aud": app.config.AuthTokenConfig.iss,
	}
	token, err := app.store.Auth.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// set the cookie
	tokenTTL := time.Hour * 24
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   app.config.Env == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(tokenTTL),
	})

	// return to the client
	if err := app.jsonResponse(w, http.StatusOK, "logged in"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getLoggedinUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.unauthorizedError(w, r, errors.New("invalid user"))
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromCtx(r)
	if session != nil {
		ctx := r.Context()

		err := app.store.Session.DeleteSession(ctx, session.SessionId)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func getSessionFromCtx(r *http.Request) *store.Session {
	session, _ := r.Context().Value(sessionCtx).(*store.Session)
	return session
}
