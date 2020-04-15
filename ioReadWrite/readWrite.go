package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	var b bytes.Buffer
	b.Write([]byte("你好"))

	_, _ = fmt.Fprintln(&b, ",", "zhaokh")

	fmt.Println(b)
	_, _ = b.WriteTo(os.Stdout)

	var p [100]byte
	n, err := b.Read(p[:])
	fmt.Println(n, err, string(p[:]))

	data, err := ioutil.ReadAll(&b)
	fmt.Println(string(data), err)
}
