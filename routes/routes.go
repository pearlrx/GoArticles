package routes

import (
	"GoArticles/handlers"
	"database/sql"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *sql.DB) {
	userHandler := handlers.NewUserHandler(db)
	articleHandler := handlers.NewArticleHandler(db)
	roleHandler := handlers.NewRoleHandler(db)
	permissionHandler := handlers.NewPermissionHandler(db)

	// Get user by id
	e.GET("/users/:id", userHandler.GetUserByID)
	// Get user role by id
	e.GET("/users/:id/roles", userHandler.GetUserRoles)
	// Get user settings by ID
	e.GET("/users/:id/settings", userHandler.GetUserSettings)
	// Creating new user
	e.POST("/users", userHandler.CreateUser)
	// Removing a user by id
	e.DELETE("/users/:id", userHandler.DeleteUser)
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

	//TODO: creating logger with logrus
}
