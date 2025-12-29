package main

import "fmt"

type something struct {
	i int
}

func (s *something) Do() int {
	return s.i
}

type doer interface {
	Do() int
}

func main() {
	var s *something
	if s == nil {
		fmt.Println("s is nil")
	}

	var d doer = s
	if d == nil {
		fmt.Println("d is nil")
		return
	}

	res := d.Do()
	fmt.Println(res)
}
