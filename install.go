package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("./install/sql.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	queries := string(b) // convert content to a 'string'
	requests := strings.Split(queries, "-----------")
	fmt.Println(requests)
}
