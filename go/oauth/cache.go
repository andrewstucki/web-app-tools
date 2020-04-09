package oauth

import (
	"net/http"
	"time"
)

var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func disableCaching(w http.ResponseWriter) {
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}
}
