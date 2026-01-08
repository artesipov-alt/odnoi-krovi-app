package dto

import "github.com/artesipov-alt/odnoi-krovi-app/ent"

// DevResponse represents a generic dev response
type DevResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// GetDeletedUsersResponse represents the response for getting deleted users
type GetDeletedUsersResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Users   []*ent.User `json:"users"`
}
