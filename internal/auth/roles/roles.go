package roles

import (
	"fmt"
	"strings"
)

// Define a custom type
type Role int

const (
	Admin Role = iota
	Editor
	User
)

// Map to string values
func (r Role) String() string {
	switch r {
	case Admin:
		return "Admin"
	case Editor:
		return "Editor"
	case User:
		return "User"
	default:
		return "Unknown"
	}
}

// Function to map string to Role
func RoleFromString(roleStr string) (Role, error) {
	// Normalize input string to handle case insensitivity
	switch strings.ToLower(roleStr) {
	case "admin":
		return Admin, nil
	case "editor":
		return Editor, nil
	case "user":
		return User, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", roleStr)
	}
}