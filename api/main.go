package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"video_server/api/auth"
	"video_server/api/handlers"
	"video_server/api/session"
)

type middlewareHandler struct {
	r *httprouter.Router
}

func NewMiddlewareHandler(r *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.r = r
	return m
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth.ValidateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	log.Println("preparing to post requests")
	router := httprouter.New()

	router.POST("/user", handlers.CreateUser)
	router.POST("/user/:user_name", handlers.Login)
	router.GET("/user/:user_name", handlers.GetUserInfo)
	router.GET("/user/", handlers.GetAllUsers)
	router.POST("/user/:user_name/videos", handlers.AddVideoInfo)
	router.GET("/user/:user_name/videos", handlers.ListAllVideos)
	router.DELETE("/user/:user_name/videos/:vid-id", handlers.DeleteVideo)
	router.POST("/videos/:vid-id/comments", handlers.PostComment)
	router.GET("/videos/:vid-id/comments", handlers.ShowComments)

	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegisterHandlers()
	mh := NewMiddlewareHandler(r)
	http.ListenAndServe(":8000", mh)
}
