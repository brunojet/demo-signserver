package http_handlers

func contains(s, substr string) bool {
	return s != "" && substr != "" && (len(s) >= len(substr)) && (s == substr || (len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
