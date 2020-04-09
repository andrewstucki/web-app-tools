package callbacks

import "html/template"

const (
	// MessageTokenRejected is displayed when a token handed back from Google has been rejected
	// for some reason, often due to an Audience or Domain mismatch
	MessageTokenRejected = "The token received was rejected, make sure you signed in with the right account."
	errorString          = `
<!doctype html>
	<body>
	{{.}}
	</body>
</html>
`
)

var (
	defaultErrorTemplate *template.Template
)

func init() {
	defaultErrorTemplate = template.Must(template.New("__oauth__error").Parse(errorString))
}
