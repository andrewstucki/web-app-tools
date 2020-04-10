package common

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/unrolled/render"
)

const (
	internalErrorHTMLString = `<!doctype html><body>Internal Server Error</body></html>`
	internalErrorJSONString = `{"reason":"Internal Server Error"}`
)

// Renderer renders HTTP responses
type Renderer interface {
	Render(w http.ResponseWriter, status int, data interface{})
	Error(w http.ResponseWriter, status int, message string)
	InternalError(w http.ResponseWriter)
}

type htmlRenderer struct {
	renderTemplate *template.Template
	errorTemplate  *template.Template
}

// NewHTMLRenderer makes a renderer that speaks HTML
func NewHTMLRenderer(renderTemplate, errorTemplate *template.Template) Renderer {
	return &htmlRenderer{
		renderTemplate: renderTemplate,
		errorTemplate:  errorTemplate,
	}
}

func (r *htmlRenderer) Render(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(status)
	if err := r.renderTemplate.Execute(w, data); err != nil {
		r.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func (r *htmlRenderer) Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(status)
	if err := r.errorTemplate.Execute(w, message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, internalErrorHTMLString)
	}
}

func (r *htmlRenderer) InternalError(w http.ResponseWriter) {
	r.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

type jsonRenderer struct {
	renderer *render.Render
}

// NewJSONRenderer makes a renderer that speaks JSON
func NewJSONRenderer() Renderer {
	return &jsonRenderer{
		renderer: render.New(),
	}
}

func (r *jsonRenderer) Render(w http.ResponseWriter, status int, data interface{}) {
	if err := r.renderer.JSON(w, status, data); err != nil {
		r.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

type errorResponse struct {
	Reason string `json:"reason"`
}

func (r *jsonRenderer) Error(w http.ResponseWriter, status int, message string) {
	if err := r.renderer.JSON(w, status, &errorResponse{
		Reason: message,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, internalErrorJSONString)
	}
}

func (r *jsonRenderer) InternalError(w http.ResponseWriter) {
	r.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
