package webapi

func ValidateObject(obj map[string]string, fields []string) bool {
	for _, key := range fields {
		if _, ok := obj[key]; !ok {
			return false
		}
	}

	return true
}
