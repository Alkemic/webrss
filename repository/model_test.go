package repository

import "testing"

func TestUser_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     []byte
		want     bool
	}{{
		name:     "correct password",
		password: "asd123",
		hash:     []byte("$2a$10$d8F54Ud7kSezcoRHrxakGeMimg3lnxS8eH/5JiN6/LtCp54ewDnHa"),
		want:     true,
	}, {
		name:     "correct password with different hash",
		password: "asd123",
		hash:     []byte("$2a$10$UOOX91bR9nRj9EjoqUceL.jg6W3EFuqtMGrnxB7T3qNSZE.u5v0WC"),
		want:     true,
	}, {
		name:     "wrong password",
		password: "asd124",
		hash:     []byte("$2a$10$UOOX91bR9nRj9EjoqUceL.jg6W3EFuqtMGrnxB7T3qNSZE.u5v0WC"),
		want:     false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				Password: tt.hash,
			}
			if got := u.ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
