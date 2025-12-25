package ct

// Returns false if control chars are present on 's'.
func controlCharsFree(s string) bool {
	for _, r := range s {
		switch r {
		case '\n', '\r', '\t':
			continue // allowed control chars
		default:
			if r < 32 {
				return false // reject other control chars
			}
		}
	}
	return true
}
