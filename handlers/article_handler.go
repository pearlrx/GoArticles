package handlers

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"

	"GoArticles/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewArticleHandler(db *sql.DB, logger *logrus.Logger) *ArticleHandler {
	return &ArticleHandler{
		DB:     db,
		Logger: logger,
	}
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
	if err = c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err = h.updateArticle(c.Request().Context(), id, article); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}

	if err = h.deleteArticle(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete article"})
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

	if err = h.addArticleTag(c.Request().Context(), articleID, tagID); err != nil {
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

	h.Logger.Infof("User ID %d likes article ID %d", userID, articleID)
	if err = h.likeArticle(c.Request().Context(), articleID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to like article"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) getArticleByID(ctx context.Context, id int) (*models.Article, error) {
	h.Logger.Infof("Fetching article with ID %d", id)

	query := `SELECT id, author_id, title, content, created_at, updated_at FROM articles WHERE id=$1`
	row := h.DB.QueryRowContext(ctx, query, id)

	var article models.Article
	if err := row.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
		h.Logger.WithError(err).Errorf("Failed to fetch article with ID %d", id)
		return nil, err
	}

	h.Logger.Infof("Successfully fetched article with ID %d", id)
	return &article, nil
}

func (h *ArticleHandler) createArticle(ctx context.Context, article models.Article) (int, error) {
	h.Logger.Infof("Creating a new article with title %s", article.Title)

	query := `INSERT INTO articles (author_id, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := h.DB.QueryRowContext(ctx, query, article.AuthorID, article.Title, article.Content, article.CreatedAt, article.UpdatedAt).Scan(&id)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to create article")
		return 0, err
	}

	h.Logger.Infof("Successfully created article with ID %d", id)
	return id, nil
}

func (h *ArticleHandler) updateArticle(ctx context.Context, id int, article models.Article) error {
	h.Logger.Infof("Updating article with ID %d", id)

	query := `UPDATE articles SET title=$1, content=$2, updated_at=$3 WHERE id=$4`
	_, err := h.DB.ExecContext(ctx, query, article.Title, article.Content, article.UpdatedAt, id)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to update article with ID %d", id)
		return err
	}

	h.Logger.Infof("Successfully updated article with ID %d", id)
	return nil
}

func (h *ArticleHandler) deleteArticle(ctx context.Context, id int) error {
	h.Logger.Infof("Deleting article with ID %d", id)

	query := `DELETE FROM articles WHERE id=$1`
	_, err := h.DB.ExecContext(ctx, query, id)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to delete article with ID %d", id)
		return err
	}

	h.Logger.Infof("Successfully deleted article with ID %d", id)
	return nil
}

func (h *ArticleHandler) addArticleTag(ctx context.Context, articleID, tagID int) error {
	h.Logger.Infof("Adding tag ID %d to article ID %d", tagID, articleID)

	query := `INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, tagID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to add tag ID %d to article ID %d", tagID, articleID)
		return err
	}

	h.Logger.Infof("Successfully added tag ID %d to article ID %d", tagID, articleID)
	return nil
}

func (h *ArticleHandler) likeArticle(ctx context.Context, articleID, userID int) error {
	h.Logger.Infof("User ID %d likes article ID %d", userID, articleID)

	query := `INSERT INTO article_likes (article_id, user_id) VALUES ($1, $2)`
	_, err := h.DB.ExecContext(ctx, query, articleID, userID)
	if err != nil {
		h.Logger.WithError(err).Errorf("Failed to like article ID %d by user ID %d", articleID, userID)
		return err
	}

	h.Logger.Infof("User ID %d successfully liked article ID %d", userID, articleID)
	return nil
}
