
package main

import (
"fmt"
"strconv"
)

/*
 Print paramater a few times; allows us to see how routines interleave.
*/
func show(param string, ch chan bool) {
	for i := 0; i < 3; i++ {
		fmt.Println(param, ":", i)
	}
	if ch != nil{
		//Tell receiver we are done
		ch <- true
	}
}

	

func main() {
	routines := 5
	ch := make(chan bool)
	show("main", nil)
	//launch a few goroutines
	for i := 0; i < routines; i++ {
		go show("routine "+strconv.Itoa(i), ch)
	}
	count := 0
	//Wait unil all routines have completed
	for count < routines - 1{
		<- ch
		count ++
	}
	fmt.Println("done")
}

