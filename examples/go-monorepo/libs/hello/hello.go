package hello

import (
	"fmt"
)

func Greet(audience string) string {
	return fmt.Sprintf("Hello, %s!", audience)
}
