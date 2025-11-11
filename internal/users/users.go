// internal/users/users.go
package users

var currentUserEmail = "student@example.com"

func GetCurrentUser() string {
	return currentUserEmail
}