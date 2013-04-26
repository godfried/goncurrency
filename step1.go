package main

import(
 "fmt"
"strconv"
)

/*
 Print paramater a few times; allows us to see how routines interleave.
*/
func show(param string) {
	for i := 0; i < 3; i++ {
		fmt.Println(param, ":", i)
	}
}

	

func main() {
	routines := 5
	show("main")
	//launch a few goroutines
	for i := 0; i < routines; i++ {
		go show("routine "+strconv.Itoa(i))
	}
//	var input string
//	fmt.Scanln(&input)
	fmt.Println("done")
}

