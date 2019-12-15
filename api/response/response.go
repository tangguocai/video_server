package response

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/defs"
)

func SendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	res, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(res))
}

func SendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
