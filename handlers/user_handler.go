package handlers

import (
	"GoArticles/models"
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"net/http"
	"strconv"
)

type UserHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewUserHandler(db *sql.DB, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		Logger: logger,
		DB:     db,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		h.Logger.WithError(err).Warn("Failed to bind user data")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	h.Logger.WithFields(logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
	}).Info("Creating a new user")

	if user.Username == "" || user.Email == "" || user.Password == "" {
		h.Logger.Warn("Empty inputs")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty inputs"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to hash password")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	user.PasswordHash = string(passwordHash)

	tx, err := h.DB.BeginTx(c.Request().Context(), nil)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to begin transaction")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to begin transaction"})
	}

	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
              VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	err = tx.QueryRowContext(c.Request().Context(), query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		tx.Rollback() // Откат транзакции при ошибке
		h.Logger.WithError(err).Error("Failed to create user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	defaultRoleID := 3
	query = `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(c.Request().Context(), query, user.ID, defaultRoleID)
	if err != nil {
		tx.Rollback() // Откат транзакции при ошибке
		h.Logger.WithError(err).Error("Failed to assign default role")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign default role"})
	}

	// Установка дефолтных настроек
	query = `INSERT INTO user_settings (user_id, setting_key, setting_value)
             SELECT $1, s.setting_key, s.default_value
             FROM settings s
             ON CONFLICT (user_id, setting_key) DO NOTHING`
	_, err = tx.ExecContext(c.Request().Context(), query, user.ID)
	if err != nil {
		tx.Rollback() // Откат транзакции при ошибке
		h.Logger.WithError(err).Error("Failed to set default user settings")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set default user settings"})
	}

	if err = tx.Commit(); err != nil {
		h.Logger.WithError(err).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	h.Logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User created successfully")
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid user ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err = h.deleteUserDependencies(c.Request().Context(), id); err != nil {
		h.Logger.WithError(err).Error("Failed to delete user dependencies")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user dependencies"})
	}

	query := `DELETE FROM users WHERE id=$1`
	result, err := h.DB.ExecContext(c.Request().Context(), query, id)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to delete user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		h.Logger.WithError(err).Error("Failed to check affected rows")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check affected rows"})
	}

	if rowsAffected == 0 {
		h.Logger.Warn("User not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	h.Logger.WithField("user_id", id).Info("User deleted successfully")
	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) deleteUserDependencies(ctx context.Context, userID int) error {
	// Пример удаления зависимостей. Обновите запросы в зависимости от вашей схемы БД.
	queries := []string{
		`DELETE FROM article_likes WHERE user_id=$1`,
		`DELETE FROM article_tags WHERE article_id IN (SELECT id FROM articles WHERE author_id=$1)`,
		`DELETE FROM article_categories WHERE article_id IN (SELECT id FROM articles WHERE author_id=$1)`,
		`DELETE FROM articles WHERE author_id=$1`,
		// Добавьте дополнительные запросы для других зависимостей
	}

	for _, query := range queries {
		_, err := h.DB.ExecContext(ctx, query, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.WithError(err).WithField("user_id", id).Warn("Invalid user ID format")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	h.Logger.WithField("user_id", id).Info("Fetching user by ID")

	user, err := h.getUserByID(c.Request().Context(), id)
	if err != nil {
		h.Logger.WithError(err).WithField("user_id", id).Warn("User not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	h.Logger.WithField("user_id", id).Info("User found successfully")
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserRoles(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.WithError(err).WithField("user_id", id).Warn("Invalid user ID format")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roles, err := h.getUserRoles(c.Request().Context(), id)
	if err != nil {
		h.Logger.WithError(err).WithField("user_id", id).Warn("Failed to retrieve roles")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve roles"})
	}

	h.Logger.WithField("user_id", id).Info("User role found successfully")
	return c.JSON(http.StatusOK, roles)
}

func (h *UserHandler) getUserByID(ctx context.Context, id int) (*models.User, error) {
	h.Logger.Infof("Fetching user with ID %d", id)

	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id=$1`
	row := h.DB.QueryRowContext(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		h.Logger.WithError(err).Errorf("Failed to fetch user with ID %d", id)
		return nil, err
	}

	h.Logger.Infof("Successfully fetched user with ID %d", id)
	return &user, nil
}

func (h *UserHandler) getUserRoles(ctx context.Context, userID int) ([]models.UserRole, error) {
	h.Logger.Infof("Fetching roles for user with ID %d", userID)

	query := `SELECT user_id, role_id FROM user_roles WHERE user_id=$1`
	rows, err := h.DB.QueryContext(ctx, query, userID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to fetch roles for user with ID %d", userID)
		return nil, err
	}
	defer rows.Close()

	var roles []models.UserRole
	for rows.Next() {
		var role models.UserRole
		if err = rows.Scan(&role.UserID, &role.RoleID); err != nil {
			h.Logger.WithError(err).Errorf("Error scanning role for user with ID %d", userID)
			return nil, err
		}
		roles = append(roles, role)
	}

	h.Logger.Infof("Successfully fetched %d roles for user with ID %d", len(roles), userID)
	return roles, nil
}
