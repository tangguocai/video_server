package handlers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"my/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"video_server/web/client"
	. "video_server/web/defs"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func HomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		pg := &HomePage{Name: "Tangguocai"}
		t, err := template.ParseFiles(`D:\GoProject\src\stream_video_server\templates\home.html`) //"./templates/home.html"
		if err != nil {
			log.Printf("Parsing template home.html error: %v", err)
			return
		}

		t.Execute(w, pg)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) == 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func UserHomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")
	var pg *UserPage
	if len(cname.Value) != 0 {
		pg = &UserPage{Name: cname.Value} // 从session中获取
	} else if len(fname) != 0 {
		pg = &UserPage{Name: fname}
	}
	//"./templates/userhome.html"D:\GoProject\src\video_server\templates\userhome.html
	t, err := template.ParseFiles(`D:\GoProject\src\stream_video_server\templates\userhome.html`)
	if err != nil {
		log.Printf("Parsing template userhome.html error: %v", err)
		return
	}
	t.Execute(w, pg)
}

func ApiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//在处理请求之前进行预处理
	if r.Method != http.MethodPost {
		res, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(res))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		res, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(res))
		return
	}

	client.Request(apiBody, w, r)
	defer r.Body.Close()
}

func ProxyVideoHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func ProxyUploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
