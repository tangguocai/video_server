package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"my/config"
	"net/http"
	"net/url"
	. "video_server/web/defs"
)

/*
	利用 http.client 作为代理转发请求
*/

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func Request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var (
		resp *http.Response
		err  error
	)

	u, _ := url.Parse(b.Url)
	u.Host = config.GetLbAddr() + ":" + u.Port()
	newUrl := u.String()
	log.Println("newUrl:", newUrl)

	switch b.Method {
	case http.MethodGet:
		req, err := http.NewRequest("GET", newUrl, nil)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		log.Println("req:", req)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		log.Println("resp:", resp)
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = r.Header
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}
}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(res))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}
