package handlers

import (
	"GoArticles/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PermissionHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewPermissionHandler(db *sql.DB, logger *logrus.Logger) *PermissionHandler {
	return &PermissionHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *PermissionHandler) AssignPermissionToRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid role ID provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid role id"})
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid permission ID provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid permission id"})
	}

	h.Logger.Infof("Assigning permission ID %d to role ID %d", permissionID, roleID)

	query := `INSERT INTO role_permissions(role_id, permission_id) VALUES ($1, $2)`
	_, err = h.DB.ExecContext(c.Request().Context(), query, roleID, permissionID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to assign permission ID %d to role ID %d", permissionID, roleID)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to assign permission to role"})
	}

	h.Logger.Infof("Successfully assigned permission ID %d to role ID %d", permissionID, roleID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Role assigned successfully"})
}

func (h *PermissionHandler) RemovePermissionFromRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid role ID provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid role id"})
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid permission ID provided")
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid permission id"})
	}

	h.Logger.Infof("Removing permission ID %d from role ID %d", permissionID, roleID)

	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err = h.DB.ExecContext(c.Request().Context(), query, roleID, permissionID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to remove permission ID %d from role ID %d", permissionID, roleID)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to remove permission from role"})
	}

	h.Logger.Infof("Successfully removed permission ID %d from role ID %d", permissionID, roleID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Permission removed successfully"})
}

func (h *PermissionHandler) GetPermissionsByRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid role ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	h.Logger.Infof("Fetching permissions for role ID %d", roleID)

	query := `SELECT p.id, p.name FROM permissions p
              JOIN role_permissions rp ON p.id = rp.permission_id
              WHERE rp.role_id = $1`
	rows, err := h.DB.QueryContext(c.Request().Context(), query, roleID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to retrieve permissions for role ID %d", roleID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve permissions"})
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err = rows.Scan(&permission.ID, &permission.Name); err != nil {
			h.Logger.WithError(err).Errorf("Failed to scan permission for role ID %d", roleID)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan permission"})
		}
		permissions = append(permissions, permission)
	}

	h.Logger.Infof("Successfully fetched %d permissions for role ID %d", len(permissions), roleID)
	return c.JSON(http.StatusOK, permissions)
}
