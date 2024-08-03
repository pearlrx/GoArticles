package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"GoArticles/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	DB *sql.DB
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
	return &ArticleHandler{DB: db}
}

func (h *ArticleHandler) GetArticleByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	article, err := h.getArticleByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Article not found"})
	}

	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	id, err := h.createArticle(c.Request().Context(), article)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create article"})
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

func (h *ArticleHandler) UpdateArticle(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	var article models.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.updateArticle(c.Request().Context(), id, article); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	if err := h.deleteArticle(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) AddArticleCategory(c echo.Context) error {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	if err := h.addArticleCategory(c.Request().Context(), articleID, categoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add category to article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) AddArticleTag(c echo.Context) error {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	tagID, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tag ID"})
	}

	if err := h.addArticleTag(c.Request().Context(), articleID, tagID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add tag to article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) LikeArticle(c echo.Context) error {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.likeArticle(c.Request().Context(), articleID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to like article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) getArticleByID(ctx context.Context, id int) (*models.Article, error) {
	query := `SELECT id, author_id, title, content, created_at, updated_at FROM articles WHERE id=$1`
	row := h.DB.QueryRowContext(ctx, query, id)

	var article models.Article
	if err := row.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
		return nil, err
	}

	return &article, nil
}

func (h *ArticleHandler) createArticle(ctx context.Context, article models.Article) (int, error) {
	query := `INSERT INTO articles (author_id, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := h.DB.QueryRowContext(ctx, query, article.AuthorID, article.Title, article.Content, article.CreatedAt, article.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *ArticleHandler) updateArticle(ctx context.Context, id int, article models.Article) error {
	query := `UPDATE articles SET title=$1, content=$2, updated_at=$3 WHERE id=$4`
	_, err := h.DB.ExecContext(ctx, query, article.Title, article.Content, article.UpdatedAt, id)
	return err
}

func (h *ArticleHandler) deleteArticle(ctx context.Context, id int) error {
	query := `DELETE FROM articles WHERE id=$1`
	_, err := h.DB.ExecContext(ctx, query, id)
	return err
}

func (h *ArticleHandler) addArticleCategory(ctx context.Context, articleID, categoryID int) error {
	query := `INSERT INTO article_categories (article_id, category_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, categoryID)
	return err
}

func (h *ArticleHandler) addArticleTag(ctx context.Context, articleID, tagID int) error {
	query := `INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, tagID)
	return err
}

func (h *ArticleHandler) likeArticle(ctx context.Context, articleID, userID int) error {
	query := `INSERT INTO article_likes (article_id, user_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, userID)
	return err
}
