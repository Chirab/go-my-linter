package fakefile

import "fmt"

// tcd his comment is not good
func fakefile() {
	var firstSlice []string
	secondeSlice := []string{"i am a slice"}
	_ = copy(firstSlice, secondeSlice)
	fmt.Println("hello world!")
}

// This is a fake comment.
func fakeParseFile() {
	fmt.Println("good by world")
	fmt.Printf("good by world")

}
