package account

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Alkemic/webrss/repository"
)

const (
	sessionCookieName = "session"
	sessionUserIDName = "userID"
	backParamName     = "back"

	LoginPageURL  = "/login"
	LogoutPageURL = "/logout"
)

var ErrMissingUserID = errors.New("missing user id in session data")

type sessionRepository interface {
	Get(sessionID string) (map[string]interface{}, error)
	Set(sessionID string, data map[string]interface{}) error
	Delete(sessionID string) error
}

func buildLoginUrl(loginURL string, req *http.Request) string {
	return loginURL + "?" + backParamName + "=" + url.PathEscape(req.RequestURI)
}

type Middleware struct {
	log         *log.Logger
	userRepo    userRepository
	sessionRepo sessionRepository
}

func NewAuthenticateMiddleware(log *log.Logger, userRepo userRepository, sessionRepo sessionRepository) *Middleware {
	return &Middleware{
		log:         log,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (m *Middleware) getUserID(sessionID string) (int, error) {
	sessionData, err := m.sessionRepo.Get(sessionID)
	if err != nil {
		return 0, fmt.Errorf("cannot get session: %w", err)
	}
	rawUserID, ok := sessionData[sessionUserIDName]
	if !ok || rawUserID == "" {
		return 0, ErrMissingUserID
	}
	userID, ok := rawUserID.(int)
	if !ok {
		return 0, ErrMissingUserID
	}
	return userID, nil
}

func (m *Middleware) LoginRequiredMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == LoginPageURL || req.URL.Path == LogoutPageURL {
			f(rw, req)
			return
		}
		sessionID, err := getSessionID(req)
		if err != nil {
			m.log.Println("cannot get session:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		userID, err := m.getUserID(sessionID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) || errors.Is(err, ErrMissingUserID) {
				if strings.Contains(req.Header.Get("Accept"), "application/json") {
					http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
				http.Redirect(rw, req, buildLoginUrl(LoginPageURL, req), http.StatusFound)
				return
			}

			m.log.Println("cannot get user id:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		user, err := m.userRepo.GetByID(req.Context(), userID)
		if err != nil {
			m.log.Println("cannot get user:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		SetUser(req, user)

		f(rw, req)
	}
}
