package callbacks

import (
	"html/template"
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/common"
)

type successData struct {
	Key      string
	Token    string
	Redirect string
}

const successString = `
<!doctype html>
  <script>try { localStorage.setItem("{{.Key}}", "{{.Token}}") } finally { window.location = "{{.Redirect}}" }</script>
</html>
`

var (
	successTemplate *template.Template
)

func init() {
	successTemplate = template.Must(template.New("__oauth__success").Parse(successString))
}

// LocalStorageCallbacks are used for storing identity tokens
// in local storage, this requires that the oauth handlers be in
// the same subdomain as the frontend that is authenticating
type LocalStorageCallbacks struct {
	errorTemplate *template.Template
	key           string
	headerKey     string
}

// NewLocalStorageCallbacks creates a new LocalStorageCallbacks instance
func NewLocalStorageCallbacks() *LocalStorageCallbacks {
	return &LocalStorageCallbacks{
		errorTemplate: defaultErrorTemplate,
		key:           "__google_id",
		headerKey:     "X-Google-Id",
	}
}

func (c *LocalStorageCallbacks) renderer() common.Renderer {
	return common.NewHTMLRenderer(successTemplate, c.errorTemplate)
}

// WithErrorTemplate allows you to override the default error template
func (c *LocalStorageCallbacks) WithErrorTemplate(template *template.Template) *LocalStorageCallbacks {
	c.errorTemplate = template
	return c
}

// WithKey allows you to override the default localStorage key used to persist the token
func (c *LocalStorageCallbacks) WithKey(key string) *LocalStorageCallbacks {
	c.key = key
	return c
}

// WithHeaderKey allows you to override the default header used to indicate a token refresh
func (c *LocalStorageCallbacks) WithHeaderKey(key string) *LocalStorageCallbacks {
	c.headerKey = key
	return c
}

// OnError return an internal server eror status
func (c *LocalStorageCallbacks) OnError(w http.ResponseWriter, err error) {
	c.renderer().Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// OnSuccess writes the new token to local storage and redirects to a given location
func (c *LocalStorageCallbacks) OnSuccess(w http.ResponseWriter, location, token string) {
	c.renderer().Render(w, http.StatusOK, successData{
		Key:      c.key,
		Token:    token,
		Redirect: location,
	})
}

// OnInvalidToken returns an invalid token status
func (c *LocalStorageCallbacks) OnInvalidToken(w http.ResponseWriter, err error) {
	c.renderer().Error(w, http.StatusUnauthorized, MessageTokenRejected)
}

// OnRefresh writes the new token to the X-Google-Id header
func (c *LocalStorageCallbacks) OnRefresh(w http.ResponseWriter, token string) error {
	w.Header().Add(c.headerKey, token)
	return nil
}
