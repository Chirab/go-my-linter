package fakefiles

import "fmt"

// tcd his comment is not good
func fakefile(s string) {
	fmt.Println(s)
	var firstSlice []string
	secondeSlice := []string{"i am a slice"}
	_ = copy(firstSlice, secondeSlice)
	fmt.Println("hello world!")
	for a := 0; a < 10; a++ {
		fmt.Println("okkkk")
	}

}

// This is a fake comment.
func fakeParseFile() {
	testSlice := []string{"a", "b", "c", "e"}
	for _, l := range testSlice {
		fmt.Println(l)
	}
	fmt.Println("good by world")
	fmt.Printf("good by world")

	for i := 0; i < 10; i++ {
		fmt.Println("okkkk")
	}

}
