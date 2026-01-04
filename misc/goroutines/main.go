package main

import "fmt"

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go func() { b <- 1 }()
	//go func() { c <- 2 }()
	go func() { <-a }()

	select {
	case a <- <-b:
		fmt.Println("1")
	case a <- <-c:
		fmt.Println("2")
	default:
		fmt.Println("3")
	}
}
