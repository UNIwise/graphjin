package qcode

// isIdentityVariable reports whether a query variable is filled from the
// authenticated identity instead of from the caller. Only the snake_case names
// are reserved. The camelCase spellings (userId, userID, ...) are ordinary
// variables: consumers pass $userId as a plain value and expect it honoured.
func isIdentityVariable(name string) bool {
	switch name {
	case "user_id", "user_id_raw", "user_id_provider":
		return true
	}

	return false
}
