package main

import (
	"net/http"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"encoding/json"
	"os"
	"io"
	"time"
	"strings"
)

const HEADER_SYNC_TIME = "SyncTime"

type RouterWebServer struct {

}

func NewRouterWebServer() *RouterWebServer {

	return &RouterWebServer{}
}

func (rws *RouterWebServer) Serve() {

	http.HandleFunc("/sync/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok && rws.passAuth(u, p) {
			rws.sync(w, u, r.Header.Get(HEADER_SYNC_TIME))
		} else {
			glog.Error("router web server sync error when no auth")
		}
	})

	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(config.FileRootDir))))
	http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok && rws.passAuth(u, p) {
			if strings.EqualFold(r.Method,"POST") {
				rws.upload(u, w, r)
			}else if strings.EqualFold(r.Method,"GET") {
				rws.download(w,r)
			}
		} else {
			glog.Error("router web server file error when no auth")
		}
	})

	addr := fmt.Sprintf("%s:%d", config.WebDomain, config.WebPort)
	http.ListenAndServe(addr, nil)
}

func (rws *RouterWebServer) passAuth(u string, p string) bool {
	user := RedisQueryUser(u)
	if user != nil && user.Uuid == u && user.Token == p {
		return true;
	} else {
		return false;
	}
}

func (rws *RouterWebServer) download(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/file/",
		http.FileServer(http.Dir(config.FileRootDir))).ServeHTTP(w, r)
}

func (rws *RouterWebServer) upload(u string, w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file");
	if err != nil {
		glog.Error(err)
		return
	}
	defer file.Close()
	tf := fmt.Sprintf("%d_%s", time.Now().Unix(), handle.Filename)
	f, err := os.OpenFile(config.FileRootDir + u + "/" + tf, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		glog.Error(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		glog.Error(err)

	}

	var r UploadResponse
	r.Code = 200
	r.FileUrl = fmt.Sprintf("%s:%d/file/%s/%s",config.WebDomain,config.WebPort,
				u,tf);
	ob,err := json.Marshal(r)
	if err == nil {
		_,err = w.Write(ob)
		if err != nil {
			glog.Error(err)
		}
	}else {
		glog.Error(err)
	}
}

func (rws *RouterWebServer) sync(w http.ResponseWriter, u string, time string) {
	glog.Infof("%s sync from time %s", u, time)
	t, err := strconv.ParseInt(time, 10, 64)
	if err == nil {
		var o []MessageBody
		err := storeManager.QueryMessage(u, t, &o)
		if err == nil {
			ob, err := json.Marshal(o)
			if err == nil {
				_, err := w.Write(ob)
				if err != nil {
					glog.Error(err)
				}
			} else {
				glog.Error(err)
			}
		} else {
			glog.Error(err)
		}
	} else {
		glog.Error(err)
	}
}




