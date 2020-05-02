package utils

import "fmt"

func AddAndFilter(filter, extra string) string {
	return fmt.Sprintf("(&(%s)(%s))", filter, extra)
}

func AddOrFilter(filter, extra string) string {
	return fmt.Sprintf("(|(%s)(%s)", filter, extra)
}