package main

import "fmt"

func main() {
	type Book struct {
		id int
	}
	bookshelf := make(chan Book, 3)

	for i := 0; i < cap(bookshelf)*2; i++ {
		select {
		case bookshelf <- Book{id: i}:
			fmt.Println("成功放上书架", i)
		default:
			fmt.Println("书架已占满")
		}
	}
	for i := 0; i < cap(bookshelf)*2; i++ {
		select {
		case book := <-bookshelf:
			fmt.Println("=取下", book)
		default:
			fmt.Println("空了")
		}
	}

}
