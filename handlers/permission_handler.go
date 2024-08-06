package handlers

import (
	"GoArticles/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type PermissionHandler struct {
	DB *sql.DB
}

func NewPermissionHandler(db *sql.DB) *PermissionHandler {
	return &PermissionHandler{DB: db}
}

func (h *PermissionHandler) AssignPermissionToRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid role id"})
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid permission id"})
	}

	query := `INSERT INTO role_permissions(role_id, permission_id) VALUES ($1, $2)`
	_, err = h.DB.ExecContext(c.Request().Context(), query, roleID, permissionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to assigned permission to role"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role assigned successfully"})
}

func (h *PermissionHandler) RemovePermissionFromRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid role id"})
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid permission id"})
	}

	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err = h.DB.ExecContext(c.Request().Context(), query, roleID, permissionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to remove permission from role"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Permission remove successfully"})
}

func (h *PermissionHandler) GetPermissionsByRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	query := `SELECT p.id, p.name FROM permissions p
              JOIN role_permissions rp ON p.id = rp.permission_id
              WHERE rp.role_id = $1`
	rows, err := h.DB.QueryContext(c.Request().Context(), query, roleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve permissions"})
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err = rows.Scan(&permission.ID, &permission.Name); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan permission"})
		}
		permissions = append(permissions, permission)
	}

	return c.JSON(http.StatusOK, permissions)
}
