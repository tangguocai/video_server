package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	. "video_server/web/handlers"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", HomeHandler)
	router.POST("/", HomeHandler)
	router.GET("/userhome", UserHomeHandler)
	router.POST("/userhome", UserHomeHandler)
	router.POST("/api", ApiHandler)
	router.GET("/videos/:vid-id", ProxyVideoHandler)
	router.POST("/upload/:vid-id", ProxyUploadHandler)
	router.ServeFiles("/statics/*filepath", http.Dir(`D:\GoProject\src\stream_video_server\templates`))

	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8080", r)
}
