package fakefile

import "fmt"

/*
	block of comment
fezfezf
fez
*/
func guillaume() []string {
	var res []string
	fmt.Println("hahaha found you")
	test := []string{"ok", "ok"}
	for _, v := range test {
		res = append(res, v)
	}
	return res
}
