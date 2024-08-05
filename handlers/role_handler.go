package handlers

import (
	"GoArticles/models"
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	DB *sql.DB
}

func NewRoleHandler(db *sql.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

func (h *RoleHandler) GetRoles(c echo.Context) error {
	roles, err := h.getRoles(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch roles"})
	}
	return c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) getRoles(ctx context.Context) ([]models.Role, error) {
	query := `SELECT id, name FROM roles`
	rows, err := h.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err = rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (h *RoleHandler) AssignRoleToUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	// Удаляем все роли у пользователя, кроме роли с ID 3
	if err = h.removeRoleFromUser(c.Request().Context(), userID, 3); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove existing roles"})
	}

	// Назначаем новую роль пользователю
	if err = h.assignRoleToUser(c.Request().Context(), userID, roleID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign new role to user"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) assignRoleToUser(ctx context.Context, userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, userID, roleID)
	return err
}

func (h *RoleHandler) RemoveRoleFromUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	if roleID == 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot remove role with ID 3"})
	}

	if err = h.removeRoleFromUser(c.Request().Context(), userID, roleID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove role from user"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) removeRoleFromUser(ctx context.Context, userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2`
	_, err := h.DB.ExecContext(ctx, query, userID, roleID)
	return err
}
