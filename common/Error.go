package common

import "strings"

func CheckExist(err string) bool {
	return strings.Contains(err, "already exists")
}
