package strconcat2

import (
	"bytes"
	"fmt"
	"strings"
)

const BLOG = "http://www.flysnow.org/"

func StringPlus(p []string) string {
	var s string
	l := len(p)
	for i := 0; i < l; i++ {
		s += p[i]
	}
	return s
}

func StringFmt(p []interface{}) string {
	return fmt.Sprint(p...)
}

func StringJoin(p []string) string {
	return strings.Join(p, "")
}
func StringBuffer(p []string) string {
	var b bytes.Buffer
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringBuilder(p []string) string {
	var b strings.Builder
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringBuilder1(p []string, cap int) string {
	var b strings.Builder
	l := len(p)
	b.Grow(cap)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func initStrings(N int) []string {
	s := make([]string, N)
	for i := 0; i < N; i++ {
		s[i] = BLOG
	}
	return s
}

func initStringi(N int) []interface {
} {
	s := make([]interface{}, N)
	for i := 0; i < N; i++ {
		s[i] = BLOG
	}
	return s
}
