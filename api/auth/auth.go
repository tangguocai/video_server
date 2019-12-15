package auth

import (
	"log"
	"net/http"
	"video_server/api/defs"
	"video_server/api/response"
	"video_server/api/session"
)

var (
	HEADER_FIELD_SESSION = "X-Session-Id"
	HEADER_FIELD_UNAME   = "X-User-Name"
)

func ValidateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		log.Println("ErrorNotAuthUser")
		response.SendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}

	return true
}
