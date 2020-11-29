package account

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Alkemic/webrss/repository"

	"github.com/google/uuid"
)

type userRepository interface {
	GetByID(ctx context.Context, id int) (repository.User, error)
	GetByEmail(ctx context.Context, email string) (repository.User, error)
	Update(ctx context.Context, user repository.User) error
}

type AuthenticateHandler struct {
	log         *log.Logger
	userRepo    userRepository
	sessionRepo sessionRepository
}

func NewAuthenticateHandler(log *log.Logger, userRepo userRepository, sessionRepo sessionRepository) *AuthenticateHandler {
	return &AuthenticateHandler{
		log:         log,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func clearSession(rw http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1, Secure: req.TLS != nil}
	http.SetCookie(rw, cookie)
}

func (a AuthenticateHandler) setSession(userID int, rw http.ResponseWriter, req *http.Request) error {
	sessionData := map[string]interface{}{
		sessionUserIDName: userID,
	}
	sessionID := uuid.New().String()
	if err := a.sessionRepo.Set(sessionID, sessionData); err != nil {
		return fmt.Errorf("cannot set session: %w", err)
	}
	cookie := &http.Cookie{
		Name:  sessionCookieName,
		Value: sessionID,
		Path:  "/", MaxAge: 60 * 60 * 24 * 31,
		Secure: req.TLS != nil,
	}
	http.SetCookie(rw, cookie)
	return nil
}

func (a *AuthenticateHandler) Login(rw http.ResponseWriter, req *http.Request) {
	var email, password string
	if req.Method == http.MethodPost {
		backURL := req.URL.Query().Get(backParamName)
		if backURL == "" {
			backURL = "/"
		}

		if err := req.ParseForm(); err != nil {
			log.Println("cannot parse form:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		email = req.Form.Get("email")
		user, err := a.userRepo.GetByEmail(req.Context(), email)
		if err != nil {
			log.Println("cannot get user:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		password = req.Form.Get("password")
		if user.ValidatePassword(password) {
			if nil := a.setSession(user.ID, rw, req); err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(rw, req, backURL, http.StatusFound)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmplData := map[string]string{
		"email":    email,
		"password": password,
		"url":      req.URL.String(),
	}
	if err := tmpl.Execute(rw, tmplData); err != nil {
		log.Println("cannot execute template:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *AuthenticateHandler) Logout(rw http.ResponseWriter, req *http.Request) {
	sessionID, err := getSessionID(req)
	if err != nil {
		a.log.Println("cannot get session:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	clearSession(rw, req)
	if err := a.sessionRepo.Delete(sessionID); err != nil {
		a.log.Println("cannot delete session:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, "/", http.StatusFound)
}

func getSessionID(req *http.Request) (string, error) {
	cookie, err := req.Cookie(sessionCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", fmt.Errorf("cannot get cookie: %w", err)
	}
	return cookie.Value, nil
}

func (a *AuthenticateHandler) Edit(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if err := req.ParseForm(); err != nil {
		a.log.Println("cannot parse form:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	email := req.Form.Get("email")
	name := req.Form.Get("name")
	password := req.Form.Get("password")
	a.log.Printf("'%s', '%s', '%s'\n", email, name, password)

	if email == "" || name == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "empty email or user name")
		return
	}

	user := GetUser(req)
	user.Email = email
	user.Name = name
	if password != "" {
		if err := user.SetPassword(password); err != nil {
			a.log.Println("cannot update password:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err := a.userRepo.Update(req.Context(), user); err != nil {
		a.log.Println("cannot save user:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
