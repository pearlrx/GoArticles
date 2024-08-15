package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"` // Используется только для ввода
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRole struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

type UserSettings struct {
	UserID       int    `json:"user_id"`
	SettingKey   string `json:"setting_key"`
	SettingValue string `json:"setting_value"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RolePermission struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}
