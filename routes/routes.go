package routes

import (
	"GoArticles/handlers"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func InitRoutes(e *echo.Echo, db *sql.DB) {

	logger := logrus.New()

	userHandler := handlers.NewUserHandler(db, logger)
	articleHandler := handlers.NewArticleHandler(db, logger)
	roleHandler := handlers.NewRoleHandler(db, logger)
	permissionHandler := handlers.NewPermissionHandler(db, logger)
	userSettingsHandler := handlers.NewUserSettingsHandler(db, logger)

	// Get user by id
	e.GET("/users/:id", userHandler.GetUserByID)
	// Get user role by id
	e.GET("/users/:id/roles", userHandler.GetUserRoles)
	// Creating new user
	e.POST("/users", userHandler.CreateUser)
	// Removing a user by id
	e.DELETE("/users/:id", userHandler.DeleteUser)

	// Get user settings by ID
	e.GET("/users/:user_id/settings", userSettingsHandler.GetUserSettings)
	e.PUT("/users/:user_id/settings", userSettingsHandler.UpdateUserSettings)

	// Get article by ID
	e.GET("/articles/:id", articleHandler.GetArticleByID)
	// Create a new article
	e.POST("/articles", articleHandler.CreateArticle)
	// Update existing article by ID
	e.PUT("/articles/:id", articleHandler.UpdateArticle)
	// Delete article by ID
	e.DELETE("/articles/:id", articleHandler.DeleteArticle)
	// Add a category to an article
	e.POST("/articles/:article_id/categories/:category_id", articleHandler.AddArticleCategory)
	// Add a tag to an article
	e.POST("/articles/:article_id/tags/:tag_id", articleHandler.AddArticleTag)
	// Like an article by user ID
	e.POST("/articles/:article_id/like/:user_id", articleHandler.LikeArticle)

	// Get role
	e.GET("/roles", roleHandler.GetRoles)
	// Add role to user
	e.PUT("/users/:user_id/roles/:role_id", roleHandler.AssignRoleToUser)
	// Delete role from user
	e.DELETE("/users/:user_id/roles/:role_id", roleHandler.RemoveRoleFromUser)

	// Add permission to role
	e.PUT("/roles/:role_id/permissions/:permission_id", permissionHandler.AssignPermissionToRole)
	// Delete permission from role
	e.DELETE("/roles/:role_id/permissions/:permission_id", permissionHandler.RemovePermissionFromRole)
	// Get permission
	e.GET("/roles/:role_id/permissions", permissionHandler.GetPermissionsByRole)

	//TODO: add setting handler
}
