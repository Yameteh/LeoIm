package main

import (
	"net/http"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"encoding/json"
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

func (rws *RouterWebServer) sync(w http.ResponseWriter, u string, time string) {
	glog.Infof("%s sync from time %s", u, time)
	t, err := strconv.ParseInt(time, 10, 64)
	if err == nil {
		var o []MessageBody
		err := storeManager.QueryMessage(u, t, o)
		if err == nil {
			ob, err := json.Marshal(o)
			if err == nil {
				w.Write(ob)
			} else {
				glog.Error(err)
			}
		}else {
			glog.Error(err)
		}
	} else {
		glog.Error(err)
	}
}




