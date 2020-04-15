package common

import (
	"fmt"
	"strconv"
	"testing"
)

//func init() {
//	Routes()
//}
//
//// 运行这个单元测试，就可以看到我们访问/sendjsonAPI的结果里，并且我们没有启动任何HTTP服务就达到了目的。
//// 这个主要利用httptest.NewRecorder()创建一个http.ResponseWriter，模拟了真实服务端的响应，
//// 这种响应时通过调用http.DefaultServeMux.ServeHTTP方法触发的。
//func TestSendJSON(t *testing.T) {
//	//req, err := http.NewRequest(http.MethodGet, "/sendjson", nil)
//	//if err != nil {
//	//	t.Fatal("创建Request失败")
//	//}
//	//
//	//rw := httptest.NewRecorder()
//	//http.DefaultServeMux.ServeHTTP(rw, req)
//	//log.Println("code: ", rw.Code)
//	//log.Println("body: ", rw.Body.String())
//
//	server := mockServer()
//	defer server.Close()
//	resq, err := http.Get(server.URL)
//	if err != nil {
//		t.Fatal("创建Get失败")
//	}
//	defer resq.Body.Close()
//
//	log.Println("code", resq.StatusCode)
//	bytes, err := ioutil.ReadAll(resq.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("body:%s\n", bytes)
//}
//func mockServer() *httptest.Server {
//	//API调用处理函数
//	sendJson := func(rw http.ResponseWriter, r *http.Request) {
//		u := struct {
//			Name string
//		}{
//			Name: "张三",
//		}
//
//		rw.Header().Set("Content-Type", "application/json")
//		rw.WriteHeader(http.StatusOK)
//		_ = json.NewEncoder(rw).Encode(u)
//	}
//	//适配器转换
//	return httptest.NewServer(http.HandlerFunc(sendJson))
//}

func TestTag(t *testing.T) {
	Tag(1)
	Tag(2)
	Tag(3)
	Tag(6)
}
func BenchmarkSprintf(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%d", num)
	}
}

func BenchmarkFormat(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(int64(num), 10)
	}
}

func BenchmarkItoa(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(num)
	}
}
