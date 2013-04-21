
package main

import "fmt"

/*
 Print paramater a few times; allows us to see how routines interleave.
*/
	func show(param string, ch chan bool) {
	for i := 0; i < 3; i++ {
		fmt.Println(param, ":", i)
	}
	if ch != nil{
		ch <- true
	}
}

	

func main() {
	ch := make(chan bool)
	show("main", nil)
	go show("routine A", ch)
	go show("routine B", ch)
	count := 0
	for count < 2{
		done := <- ch
		if done{
			count ++
		}
	}
	fmt.Println("done")
}

