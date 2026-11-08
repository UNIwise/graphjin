package qcode

import "testing"

func TestIsIdentityVariable(t *testing.T) {
	reserved := []string{"user_id", "user_id_raw", "user_id_provider"}
	for _, name := range reserved {
		if !isIdentityVariable(name) {
			t.Errorf("%q must be filled from the authenticated identity", name)
		}
	}

	// Ordinary caller-supplied variables. Policies and services pass $userId
	// and similar names as plain values, so they must not be treated as identity.
	ordinary := []string{
		"userId", "userID", "userIdRaw", "userIDRaw", "userIdProvider", "userIDProvider",
		"participantUserId", "user_id2", "USER_ID", "",
	}
	for _, name := range ordinary {
		if isIdentityVariable(name) {
			t.Errorf("%q must be an ordinary variable, not an identity variable", name)
		}
	}
}
