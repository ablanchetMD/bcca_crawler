package roles

import (
	"fmt"
	"strings"
)

// Define a custom type
type Role int

const (
	Guest Role = iota
	User
	Editor
	Admin 	
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
	case Guest:
		return "Guest"
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
	case "guest":
		return Guest, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", roleStr)
	}
}