package utils

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/stretchr/testify/mock"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

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

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}
	return sb.String()
}
