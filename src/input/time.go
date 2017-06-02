package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("%X", time.Now().Nanosecond())
}
