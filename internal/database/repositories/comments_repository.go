package repositories

import (
	"github.com/KlintonLee/go-rest-api-course/internal/database/models"
	"github.com/jinzhu/gorm"
)

type CommentsRepository interface {
	PostComment(comment *models.Comment) (*models.Comment, error)
	GetComment(id string) (*models.Comment, error)
	GetCommentsBySlug(slug string) ([]models.Comment, error)
	GetAllComments() ([]models.Comment, error)
	UpdateComment(id string, newComment *models.Comment) (*models.Comment, error)
	DeleteComment(id string) error
}

type CommentsRepositoryDB struct {
	Db *gorm.DB
}

func (repo *CommentsRepositoryDB) PostComment(comment *models.Comment) (*models.Comment, error) {
	err := repo.Db.Create(comment).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentsRepositoryDB) GetComment(id string) (*models.Comment, error) {
	var comment models.Comment
	err := repo.Db.First(&comment, "id = ?", id).Error
	if err != nil {
		return &models.Comment{}, err
	}

	return &comment, nil
}

func (repo *CommentsRepositoryDB) GetCommentsBySlug(slug string) ([]models.Comment, error) {
	var comments []models.Comment
	err := repo.Db.Find(&comments, "slug = ?", slug).Error
	if err != nil {
		return []models.Comment{}, nil
	}

	return comments, nil
}

func (repo *CommentsRepositoryDB) GetAllComments() ([]models.Comment, error) {
	var comments []models.Comment
	err := repo.Db.Find(&comments).Error
	if err != nil {
		return []models.Comment{}, nil
	}

	return comments, nil
}

func (repo *CommentsRepositoryDB) UpdateComment(id string, newComment *models.Comment) (*models.Comment, error) {
	comment, err := repo.GetComment(id)
	if err != nil {
		return &models.Comment{}, err
	}

	err = repo.Db.Model(comment).Updates(newComment).Error
	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func (repo *CommentsRepositoryDB) DeleteComment(id string) error {
	comment, err := repo.GetComment(id)
	if err != nil {
		return err
	}

	err = repo.Db.Delete(comment).Error
	if err != nil {
		return err
	}

	return nil
}
