package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Alkemic/forms"
	"github.com/google/uuid"

	"github.com/Alkemic/webrss/repository"
)

type settingsRepository interface {
	GetUser(ctx context.Context) (repository.User, error)
	SaveUser(ctx context.Context, user repository.User) error
}

type AuthenticateHandler struct {
	log          *log.Logger
	settingsRepo settingsRepository
	sessionRepo  sessionRepository
}

func NewAuthenticateHandler(log *log.Logger, settingsRepo settingsRepository, sessionRepo sessionRepository) *AuthenticateHandler {
	return &AuthenticateHandler{
		log:          log,
		settingsRepo: settingsRepo,
		sessionRepo:  sessionRepo,
	}
}

func clearSession(rw http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1, Secure: req.TLS != nil}
	http.SetCookie(rw, cookie)
}

func (a AuthenticateHandler) setUserSession(username string, rw http.ResponseWriter, req *http.Request) error {
	sessionID, err := getSessionID(req)
	if err != nil || sessionID == "" {
		sessionID = uuid.New().String()
	}
	sessionData := map[string]interface{}{
		sessionUsernameKey: username,
	}
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

func getBackURL(req *http.Request) string {
	backURL := req.URL.Query().Get(backParamName)
	if backURL == "" {
		return "/"
	}
	return backURL
}

func (a *AuthenticateHandler) Login(rw http.ResponseWriter, req *http.Request) {
	loginForm := newLoginForm()

	if err := req.ParseForm(); err != nil {
		a.log.Println("cannot parse form:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodPost && loginForm.IsValid(req.PostForm) {
		username := loginForm.CleanedData["username"].(string)
		password := loginForm.CleanedData["password"].(string)
		user, err := a.settingsRepo.GetUser(req.Context())
		if errors.Is(err, sql.ErrNoRows) {
			loginForm.AddError("Incorrect email or password")
		} else if err != nil {
			a.log.Println("cannot get user:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			if user.Validate(username, password) {
				if err := a.setUserSession(user.Name, rw, req); err != nil {
					a.log.Println("cannot set session:", err)
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				http.Redirect(rw, req, getBackURL(req), http.StatusFound)
				return
			} else {
				loginForm.AddError("Incorrect email or password")
			}
		}
	}

	tmplData := struct {
		Form *forms.Form
		URL  string
	}{
		Form: loginForm,
		URL:  req.URL.String(),
	}
	if err := template.Must(template.ParseFiles("templates/login.html")).Execute(rw, tmplData); err != nil {
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
	username := req.Form.Get("name")
	password := req.Form.Get("password")

	if username == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "empty username")
		return
	}

	user := GetUser(req)
	user.Name = username
	if password != "" {
		if err := user.SetPassword(password); err != nil {
			a.log.Println("cannot update password:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err := a.settingsRepo.SaveUser(req.Context(), user); err != nil {
		a.log.Println("cannot save user:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := a.setUserSession(user.Name, rw, req); err != nil {
		a.log.Println("cannot set session:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
