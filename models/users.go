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
	UserID       int `json:"user_id"`
	SettingKey   int `json:"setting_key"`
	SettingValue int `json:"setting_value"`
}
