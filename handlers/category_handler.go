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

type CategoryHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewCategoryHandler(DB *sql.DB, Logger *logrus.Logger) *CategoryHandler {
	return &CategoryHandler{
		DB,
		Logger,
	}
}

func (h *CategoryHandler) AddArticleCategory(c echo.Context) error {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid article ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid category ID provided")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	if err = h.addArticleCategory(c.Request().Context(), articleID, categoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add category to article"})
	}
	
	return c.NoContent(http.StatusNoContent)
}

func (h *CategoryHandler) addArticleCategory(ctx context.Context, articleID, categoryID int) error {
	h.Logger.Infof("Adding category ID %d to article ID %d", categoryID, articleID)

	query := `INSERT INTO article_categories (article_id, category_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, categoryID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to add category ID %d to article ID %d", categoryID, articleID)
		return err
	}

	h.Logger.Infof("Successfully added category ID %d to article ID %d", categoryID, articleID)
	return nil
}

func (h *CategoryHandler) RemoveCategoryFromArticle(c echo.Context) error {
	articleID := c.Param("article_id")
	categoryID := c.Param("category_id")

	query := `DELETE FROM article_categories WHERE article_id = $1 AND category_id = $2`
	result, err := h.DB.Exec(query, articleID, categoryID)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to remove category from article")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove category from article"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found for the article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CategoryHandler) GetCategoriesForArticle(c echo.Context) error {
	articleID := c.Param("article_id")
	var categories []models.Category

	query := `SELECT c.id, c.name
              FROM categories c
              JOIN article_categories ac ON c.id = ac.category_id
              WHERE ac.article_id = $1`
	h.Logger.Infof("Executing query to retrieve articles for category ID: %s", articleID)

	rows, err := h.DB.Query(query, articleID)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to retrieve categories for article")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve categories for article"})
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			h.Logger.WithError(err).Error("Failed to scan category")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process categories"})
		}
		categories = append(categories, category)
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetArticlesForCategory(c echo.Context) error {
	categoryID := c.Param("category_id")
	h.Logger.Infof("Received request to get articles for category ID: %s", categoryID)

	var articles []models.Article

	query := `SELECT a.id, a.title, a.author_id, a.content, a.created_at, a.updated_at
          FROM articles a
          JOIN article_categories ac ON a.id = ac.article_id
          WHERE ac.category_id = $1`
	h.Logger.Infof("Executing query to retrieve articles for category ID: %s", categoryID)

	rows, err := h.DB.Query(query, categoryID)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to retrieve articles for category")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve articles for category"})
	}
	defer func() {
		if err = rows.Close(); err != nil {
			h.Logger.WithError(err).Warn("Failed to close rows")
		}
	}()

	for rows.Next() {
		var article models.Article
		h.Logger.Infof("Processing article for category ID: %s", categoryID)
		if err = rows.Scan(&article.ID, &article.Title, &article.AuthorID, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
			h.Logger.WithError(err).Error("Failed to scan article")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process articles"})
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		h.Logger.WithError(err).Error("Error occurred during rows iteration")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error occurred during rows iteration"})
	}

	h.Logger.Infof("Successfully retrieved %d articles for category ID: %s", len(articles), categoryID)
	return c.JSON(http.StatusOK, articles)
}
