package main

import (
	"doc/internal/psql"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	psql.Init_db()
}
