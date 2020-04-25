package main

import "fmt"

var (
	Red   = Color("\033[0;31m%s\033[0m")
	Green = Color("\033[0;32m%s\033[0m")
	Cyan  = Color("\033[0;36m%s\033[0m")
	White = Color("\033[0;37m%s\033[0m")
	Faint = Color("\033[2;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
