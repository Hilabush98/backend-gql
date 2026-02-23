package cache

func Str(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
