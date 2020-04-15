package common

import (
	"fmt"
)

//func Routes() {
//	http.HandleFunc("/sendjson", SendJSON)
//}
//
//func SendJSON(rw http.ResponseWriter, r *http.Request) {
//	u := struct {
//		Name string
//	}{
//		"张三",
//	}
//
//	rw.Header().Set("Content-Type", "application/json")
//	rw.WriteHeader(http.StatusOK)
//	_ = json.NewEncoder(rw).Encode(u)
//}

func Tag(tag int) {
	switch tag {
	case 1:
		fmt.Println("Android")
	case 2:
		fmt.Println("Go")
	case 3:
		fmt.Println("Java")
	default:
		fmt.Println("C")
	}
}
