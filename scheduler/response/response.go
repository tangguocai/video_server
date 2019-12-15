package response

import (
	"net/http"
	"io"
)

func SendResponse(w http.ResponseWriter, sc int, resp string) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
