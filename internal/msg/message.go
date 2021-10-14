package msg

import (
	"fmt"
)

var (
	check = "  \u2713 "
	x     = "  \u2717 "
	point = "  \u2022 "
)

// Info formats info message
func Info(message string) {
	fmt.Println()
	fmt.Println(point + message)
}

// Failure formats failure message
func Failure(message string) {
	fmt.Println()
	fmt.Println(x + message)
}

// Success formats success message
func Success(message string) {
	fmt.Println()
	fmt.Println(check + message)
}
