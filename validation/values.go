package validation

// IsValidHostnamePort returns true if the given hostport is a valid host:port
func IsValidHostnamePort(hostport string) bool {
	if err := V.Var(hostport, "required,hostname_port"); err != nil {
		return false
	}
	return true
}
