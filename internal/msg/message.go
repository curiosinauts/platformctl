package msg

import (
	"fmt"
	"os"
)

var (
	check = "  \u2713 "
	x     = "  \u2717 "
	point = "  \u2022 "
)

func Info(message string) {
	fmt.Println(point + message)
	fmt.Println()
}

func Formaterr(err error) {
	fmt.Println("   ", err.Error())
	fmt.Println()
}

// Error prints error then exit if err is not nil
func Error(err error) {
	if err != nil {
		Formaterr(err)
		os.Exit(1)
	}
}

// Failure formats failure message
func Failure(message string) {
	fmt.Println(x + message)
	fmt.Println()
}

// Success formats success message
func Success(message string) {
	fmt.Println(check + message)
	fmt.Println()
}
