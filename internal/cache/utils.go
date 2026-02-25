package cache

import "strings"

func Str(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func NilValidate(s *string) *string {
	if s == nil {
		return nil
	}
	if strings.TrimSpace(*s) == "" {
		return nil
	}
	return s
}
