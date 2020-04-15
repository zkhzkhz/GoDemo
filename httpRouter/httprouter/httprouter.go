package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//_, _ = fmt.Fprintf(w, "Blog:%s \nWechat:%s", "www.flysnow.org", p.ByName("name"))
	panic("dafafferf异常")
}
func main() {
	router := httprouter.New()
	router.GET("/user/*name", Index)

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error:%s", v)
	}

	log.Fatal(http.ListenAndServe(":8081", router))
}
