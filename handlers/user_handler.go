package handlers

import (
	"GoArticles/models"
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Log received user data
	log.Printf("Received User Data: %+v", user)

	// Validate input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty inputs"})
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	user.PasswordHash = string(passwordHash)

	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
              VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	// Log query and data being sent to the database
	log.Printf("Executing query: %s with values: %s, %s, %s", query, user.Username, user.Email, user.PasswordHash)

	err = h.DB.QueryRowContext(c.Request().Context(), query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)

	if err != nil {
		log.Printf("Error executing query: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	query := `DELETE FROM users WHERE id=$1`
	result, err := h.DB.ExecContext(c.Request().Context(), query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check affected rows"})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.getUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserRoles(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roles, err := h.getUserRoles(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve roles"})
	}

	return c.JSON(http.StatusOK, roles)
}

func (h *UserHandler) GetUserSettings(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	settings, err := h.getUserSettings(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve settings"})
	}

	return c.JSON(http.StatusOK, settings)
}

func (h *UserHandler) getUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id=$1`
	row := h.DB.QueryRowContext(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *UserHandler) getUserRoles(ctx context.Context, userID int) ([]models.UserRole, error) {
	query := `SELECT user_id, role_id FROM user_roles WHERE user_id=$1`
	rows, err := h.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.UserRole
	for rows.Next() {
		var role models.UserRole
		if err := rows.Scan(&role.UserID, &role.RoleID); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (h *UserHandler) getUserSettings(ctx context.Context, userID int) ([]models.UserSettings, error) {
	query := `SELECT user_id, setting_key, setting_value FROM user_settings WHERE user_id=$1`
	rows, err := h.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.UserSettings
	for rows.Next() {
		var setting models.UserSettings
		if err := rows.Scan(&setting.UserID, &setting.SettingKey, &setting.SettingValue); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}
