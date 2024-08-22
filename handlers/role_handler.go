package handlers

import (
	"GoArticles/models"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewRoleHandler(db *sql.DB, logger *logrus.Logger) *RoleHandler {
	return &RoleHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *RoleHandler) GetRoles(c echo.Context) error {
	h.Logger.Info("Fetching all roles")
	roles, err := h.getRoles(c.Request().Context())
	if err != nil {
		h.Logger.WithError(err).Error("Failed to fetch roles from database")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch roles"})
	}
	h.Logger.Infof("Fetched %d roles successfully", len(roles))
	return c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) getRoles(ctx context.Context) ([]models.Role, error) {
	query := `SELECT id, name FROM roles`
	rows, err := h.DB.QueryContext(ctx, query)
	if err != nil {
		h.Logger.WithError(err).Error("Error executing query to fetch roles")
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err = rows.Scan(&role.ID, &role.Name); err != nil {
			h.Logger.WithError(err).Error("Error scanning role from database")
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (h *RoleHandler) AssignRoleToUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid user ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid role ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	if err = h.removeRoleFromUser(c.Request().Context(), userID, 3); err != nil {
		h.Logger.WithError(err).Error("Failed to remove existing roles from user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove existing roles"})
	}

	if err = h.assignRoleToUser(c.Request().Context(), userID, roleID); err != nil {
		h.Logger.WithError(err).Error("Failed to assign new role to user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign new role to user"})
	}

	h.Logger.Infof("Role %d assigned to user %d successfully", roleID, userID)
	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) assignRoleToUser(ctx context.Context, userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		h.Logger.WithError(err).Error("Error executing query to assign role to user")
	}
	return err
}

func (h *RoleHandler) RemoveRoleFromUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid user ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid role ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	if roleID == 3 {
		h.Logger.Warn("Attempted to remove role with ID 3, which is not allowed")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot remove role with ID 3"})
	}

	h.Logger.Infof("Removing role %d from user %d", roleID, userID)

	if err = h.removeRoleFromUser(c.Request().Context(), userID, roleID); err != nil {
		h.Logger.WithError(err).Error("Failed to remove role from user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove role from user"})
	}

	h.Logger.Infof("Role %d removed from user %d successfully", roleID, userID)
	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) removeRoleFromUser(ctx context.Context, userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2`
	_, err := h.DB.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		h.Logger.WithError(err).Error("Error executing query to remove role from user")
	}
	return err
}
