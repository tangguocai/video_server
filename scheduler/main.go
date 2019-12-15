package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"video_server/scheduler/handlers"
	"video_server/scheduler/task_runner"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", handlers.VidDelRecHandler)

	return router
}

func main() {
	go task_runner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
