package main

import (
	"fmt"
	"strconv"
)

func main() {

	counter := 1
	limit := 100000000
	//thereshold := 25.0
	for counter <= limit {
		fmt.Print("COUNTER AT " + strconv.Itoa(counter) + " ")
		counter++
	}
}
