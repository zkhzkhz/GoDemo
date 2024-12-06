package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	//http.HandleFunc("/", Index)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServe(":8080", &CustomMux{}))
}

//func Index(w http.ResponseWriter, r *http.Request) {
//	_, _ = fmt.Fprint(w, "Blog:www.flysnow.org\nwechat:flysnow_org")
//}

type CustomMux struct {
}

func (cm *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, _ = fmt.Fprint(w, "Blog:www.flysnow.org\nwechat:flysnow_org")
	} else {
		_, _ = fmt.Fprint(w, "bad http method request")
	}
}
