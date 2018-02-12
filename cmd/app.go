package main

import (
	"fmt"
	"strconv"
)

func main() {
	counter := 1
	limit := 100000000
	for counter <= limit {
		fmt.Println("COUNTER AT " + strconv.Itoa(counter))
		counter++
	}
}
