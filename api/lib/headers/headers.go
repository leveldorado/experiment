package headers

import "net/http"

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

func SetContentType(h http.Header, value string) {
	h.Set(ContentType, ApplicationJSON)
}
