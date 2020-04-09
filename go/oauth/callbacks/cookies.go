package callbacks

import (
	"html/template"
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/common"
)

//CookieCallbacks set stuff in cookies
type CookieCallbacks struct {
	errorTemplate *template.Template
	path          string
	key           string
	secure        bool
}

// NewCookiesCallbacks creates a new CookieCallbacks instance
func NewCookiesCallbacks(secure bool) *CookieCallbacks {
	return &CookieCallbacks{
		errorTemplate: defaultErrorTemplate,
		path:          "__google_oauth",
		key:           "__google_id",
		secure:        secure,
	}
}

func (c *CookieCallbacks) renderer() common.Renderer {
	return common.NewHTMLRenderer(nil, c.errorTemplate)
}

// WithErrorTemplate allows you to override the default error template
func (c *CookieCallbacks) WithErrorTemplate(template *template.Template) *CookieCallbacks {
	c.errorTemplate = template
	return c
}

// WithKey allows you to override the default cookie key used to persist the token
func (c *CookieCallbacks) WithKey(key string) *CookieCallbacks {
	c.key = key
	return c
}

// WithPath allows you to override the default path used to store the cookie
func (c *CookieCallbacks) WithPath(path string) *CookieCallbacks {
	c.path = path
	return c
}

// OnError return an internal server eror status
func (c *CookieCallbacks) OnError(w http.ResponseWriter, err error) {
	c.renderer().Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// OnSuccess writes the new token to local storage and redirects to a given location
func (c *CookieCallbacks) OnSuccess(w http.ResponseWriter, location, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.key,
		Value:    token,
		Path:     c.path,
		HttpOnly: true,
		Secure:   c.secure,
		MaxAge:   1 * 60 * 60,
	})
	w.WriteHeader(http.StatusFound)
	w.Header().Set("Location", location)
}

// OnInvalidToken returns an invalid token status
func (c *CookieCallbacks) OnInvalidToken(w http.ResponseWriter, err error) {
	c.renderer().Error(w, http.StatusUnauthorized, MessageTokenRejected)
}

// OnRefresh writes the new token to the X-Google-Id header
func (c *CookieCallbacks) OnRefresh(w http.ResponseWriter, token string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     c.key,
		Value:    token,
		Path:     c.path,
		HttpOnly: true,
		Secure:   c.secure,
		MaxAge:   1 * 60 * 60,
	})
	return nil
}
