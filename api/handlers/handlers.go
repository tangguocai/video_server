package handlers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"video_server/api/auth"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/response"
	"video_server/api/session"
	"video_server/api/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// read body from request
	reqBody, _ := ioutil.ReadAll(r.Body)
	uBody := &defs.UserCredential{}

	if err := json.Unmarshal(reqBody, uBody); err != nil {
		response.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(uBody.UserName, uBody.Pwd); err != nil {
		log.Println("DBError:", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(uBody.UserName)
	su := &defs.SignedUp{
		Success:   true,
		SessionId: id,
	}

	if resp, err := json.Marshal(su); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	/*uname := p.ByName("user_name")
	io.WriteString(w, uname)*/
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		//io.WriteString(w, "wrong")
		response.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	uname := p.ByName("user_name")
	log.Printf("Login url name: %s", uname)
	log.Printf("Login body name: %s", ubody.UserName)
	if uname != ubody.UserName {
		response.SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	log.Printf("%s", ubody.UserName)
	pwd, err := dbops.GetUserCredential(ubody.UserName)
	log.Printf("Login pwd: %s", pwd)
	log.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		response.SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	id := session.GenerateNewSessionId(ubody.UserName)
	si := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	if uname != "" {
		if !auth.ValidateUser(w, r) {
			log.Println("Unathorized user")
			return
		}
	}
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("Error in GetUserInfo: %v", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u, err := dbops.GetAllUsers()
	if err != nil {
		log.Printf("Error in GetAllUser: %v", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := &defs.Users{Users: u}
	if resp, err := json.Marshal(ui); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}

func AddVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !auth.ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	nvBody := &defs.NewVideo{}
	if err := json.Unmarshal(reqBody, nvBody); err != nil {
		log.Printf("%v", err)
		response.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddVideoInfo(nvBody.AuthorId, nvBody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo: %v", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	/*uname := p.ByName("user_name")
	if uname != "" {
		if !auth.ValidateUser(w, r) {
			log.Println("Unathorized user")
			return
		}
	}*/

	vs, err := dbops.ListVideoInfo( /*uname,*/ 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !auth.ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("Error in DeletVideo: %s", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	response.SendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !auth.ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cBody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cBody); err != nil {
		log.Printf("%v", err)
		response.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cBody.AuthorId, cBody.Content); err != nil {
		log.Printf("Error in PostComment: %s", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
	} else {
		response.SendNormalResponse(w, "ok", 201)
	}
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !auth.ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		response.SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		response.SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		response.SendNormalResponse(w, string(resp), 200)
	}
}
