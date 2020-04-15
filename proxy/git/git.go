package main

import (
	"fmt"
	"strings"
)

type Git interface {
	Clone(url string) bool
}

type GitHub struct{}

func (p GitHub) Clone(url string) bool {
	if strings.HasPrefix(url, "https") {
		fmt.Println("clone from " + url)
		return true
	}

	fmt.Println("failed to clone from " + url)
	return false
}

type GitBash struct {
	Gitcmd Git
}

func (p GitBash) Clone(url string) bool {
	return p.Gitcmd.Clone(url)
}

type Coder struct {
}

func (c Coder) GetCode(url string) {
	gitBash := GetGit(1)
	if gitBash.Clone(url) {
		fmt.Println("success")
	} else {
		fmt.Println("failed")
	}
}

func GetGit(t int) Git {
	if t == 1 {
		return GitBash{Gitcmd: GitHub{}}
	}
	return nil
}

func main() {
	coder := Coder{}
	coder.GetCode("https://github.com/qibin0506/go-designpattern")
	coder.GetCode("http://github.com/qibin0506/go-designpattern")
}
