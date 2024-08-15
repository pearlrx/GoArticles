package handlers

import (
	"GoArticles/models"
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserSettingsHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewUserSettingsHandler(db *sql.DB, logger *logrus.Logger) *UserSettingsHandler {
	return &UserSettingsHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *UserSettingsHandler) GetUserSettings(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.Logger.WithError(err).Error("Invalid user ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	settings, err := h.getUserSettings(c.Request().Context(), userID)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to fetch user settings")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user settings"})
	}

	return c.JSON(http.StatusOK, settings)
}

func (h *UserSettingsHandler) UpdateUserSettings(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid user ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var settings []models.UserSettings
	if err = c.Bind(&settings); err != nil {
		h.Logger.WithError(err).Warn("Failed to bind settings data")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	tx, err := h.DB.BeginTx(c.Request().Context(), nil)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to begin transaction")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to begin transaction"})
	}

	for _, setting := range settings {
		query := `INSERT INTO user_settings (user_id, setting_key, setting_value)
                  VALUES ($1, $2, $3)
                  ON CONFLICT (user_id, setting_key) DO UPDATE
                  SET setting_value = EXCLUDED.setting_value`
		_, err = tx.ExecContext(c.Request().Context(), query, userID, setting.SettingKey, setting.SettingValue)
		if err != nil {
			tx.Rollback() // Откат транзакции при ошибке
			h.Logger.WithError(err).Error("Failed to update user settings")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user settings"})
		}
	}

	if err = tx.Commit(); err != nil {
		h.Logger.WithError(err).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	h.Logger.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("User settings updated successfully")
	return c.NoContent(http.StatusNoContent)
}

func (h *UserSettingsHandler) getUserSettings(ctx context.Context, userID int) ([]models.UserSettings, error) {
	query := `SELECT user_id, setting_key, setting_value FROM user_settings WHERE user_id = $1`
	rows, err := h.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.UserSettings
	for rows.Next() {
		var setting models.UserSettings
		if err = rows.Scan(&setting.UserID, &setting.SettingKey, &setting.SettingValue); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

func (h *UserSettingsHandler) updateUserSettings(ctx context.Context, settings models.UserSettings) error {
	query := `UPDATE user_settings SET setting_value = $1 WHERE user_id = $2 AND setting_key = $3`
	_, err := h.DB.ExecContext(ctx, query, settings.SettingValue, settings.UserID, settings.SettingKey)
	return err
}
