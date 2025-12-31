package helpers

import "strings"

func NormalizeLocation(location string) (string, error) {
	s := strings.ToLower(location)
	s = strings.ReplaceAll(s, ",", " ")
	s = strings.Join(strings.Fields(s), " ")
	return s, nil
}
