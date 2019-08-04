package main

import (
	"fmt"
	"os"
)

func main() {
	x := "C:\\Users\\admin\\python\\de_stats\\match24.dem "
	fmt.Println(x)
	file, err := os.Open(x)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}
