package transport

import (
	"net/http"
)

type Transport interface {
	HandleRequest(http.ResponseWriter, *http.Request)
}
