package main

import "fmt"

func panic(v interface{}) {
	fmt.Printf("%s\n", v)
}

func main() {
	panic("not an actual panic")
}
