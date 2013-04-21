package main

import "fmt"

/*
 Print paramater a few times; allows us to see how routines interleave.
*/
func show(param string) {
	for i := 0; i < 3; i++ {
		fmt.Println(param, ":", i)
	}
}

	

func main() {
	show("main")
	go show("routine A")
	go show("routine B")
//	var input string
//	fmt.Scanln(&input)
	fmt.Println("done")
}

