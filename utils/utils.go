package utils

import (
	"fmt"
	"strings"

	"github.com/stretchr/testify/mock"
)

func GeneratePlaceHolders(n int) string {
	s := ``
	sep := ","
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		s += fmt.Sprintf("$%d", i) + sep
	}
	strings.Split(s, sep)
	return s
}

func GenerateUpdatePlaceHolders(fields []string) string {
	s := ``
	sep := ", "
	for i, field := range fields {
		if i == len(fields)-1 {
			sep = ""
		}
		s += fmt.Sprintf(field+"=$%d", i+1) + sep
	}

	return s
}

func GenerateMockArguments(n int, args ...interface{}) []interface{} {
	for i := 0; i < n; i++ {
		args = append(args, mock.Anything)
	}
	return args
}
