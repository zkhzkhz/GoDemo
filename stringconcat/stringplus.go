package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {

}

func StringPlus() string {
	//运行 go build -gcflags="-m -m" XXX.go查看编译器对代码的优化
	//.\stringplus.go:13:6: can inline StringPlus as: func() string { var s string; s = <N>; s += "昵称:飞雪无情\n"; s += "博客:www.flysnow.org/\n"; s += "微信公众号:flysnow_org"; return s }
	//编译器帮我们把字符串做了优化，只剩下3个s+=
	var s string
	s += "昵称" + ":" + "飞雪无情" + "\n"
	s += "博客" + ":" + "www.flysnow.org/" + "\n"
	s += "微信公众号" + ":" + "flysnow_org"
	return s
}

func StringFmt() string {
	return fmt.Sprint("昵称", ":", "飞雪无情", "\n", "博客", ":", "http://www.flysnow.org/", "\n", "微信公众号", ":", "flysnow_org")
}

func StringJoin() string {
	s := []string{"昵称", ":", "飞雪无情", "\n", "博客", ":", "http://www.flysnow.org/", "\n", "微信公众号", ":", "flysnow_org"}
	return strings.Join(s, "")
}

func StringBuffer() string {
	var b bytes.Buffer
	b.WriteString("昵称")
	b.WriteString(":")
	b.WriteString("飞雪无情")
	b.WriteString("\n")
	b.WriteString("博客")
	b.WriteString(":")
	b.WriteString("http://www.flysnow.org/")
	b.WriteString("\n")
	b.WriteString("微信公众号")
	b.WriteString(":")
	b.WriteString("flysnow_org")
	return b.String()
}

func StringBuilder() string {
	var b strings.Builder
	b.WriteString("昵称")
	b.WriteString(":")
	b.WriteString("飞雪无情")
	b.WriteString("\n")
	b.WriteString("博客")
	b.WriteString(":")
	b.WriteString("http://www.flysnow.org/")
	b.WriteString("\n")
	b.WriteString("微信公众号")
	b.WriteString(":")
	b.WriteString("flysnow_org")
	return b.String()
}
