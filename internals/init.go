package internals

import (
	"fmt"
)

func GetTerm() string {
	var term string
	fmt.Println("Input for searching...")
	fmt.Scanln(&term)
	return term
}
