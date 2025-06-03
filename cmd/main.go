package main

import (
	"fmt"

	"github.com/SergeyWhite/EXPDoc/internal/psql"
)

func main() {
	fmt.Println("Hello, World!")
	psql.Init_db()
}
