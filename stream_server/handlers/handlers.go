package handlers

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	. "video_server/stream_server/defs"
	"video_server/stream_server/ossops"
	"video_server/stream_server/response"
)

func TestPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}

// 播放
func StreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	/*vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, "Error of open video: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()*/
	log.Println("Enter the streamHandler")
	tartgetUrl := "http://tgcvideos.oss-cn-beijing.aliyuncs.com/videos/" + p.ByName("vid-id")
	http.Redirect(w, r, tartgetUrl, 301)
}

// 上传
func UploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "File is too large")
		return
	}

	// FormFile => <form name="file"
	file, _, err := r.FormFile("file")
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ossFn := "videos/" + fn
	path := VIDEO_DIR + fn
	bn := "tgcvideos"
	ret := ossops.UploadToOss(ossFn, path, bn)
	if !ret {
		response.SendErrorResponse(w, http.StatusInternalServerError, "Internal error")
		return
	}

	os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully.")
}
