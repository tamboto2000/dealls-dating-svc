package main

import (
	"fmt"
	"time"
)

func main() {
	m := make(map[string]any)
	m["val"] = "hello"

	t, _ := m["val"].(time.Time)

	fmt.Println(t)
}
