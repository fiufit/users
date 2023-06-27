package repositories

import (
	"context"
	"errors"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/contracts/certifications"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Certifications interface {
	Create(ctx context.Context, certification models.Certification) (models.Certification, error)
	Get(ctx context.Context, request certifications.GetCertificationsRequest) (certifications.GetCertificationsResponse, error)
	GetByID(ctx context.Context, id uint) (models.Certification, error)
	Update(ctx context.Context, certification models.Certification) (models.Certification, error)
}

type CertificationRepository struct {
	db       *gorm.DB
	logger   *zap.Logger
	firebase Firebase
}

func NewCertificationRepository(db *gorm.DB, logger *zap.Logger, firebase Firebase) CertificationRepository {
	return CertificationRepository{db: db, logger: logger, firebase: firebase}
}

func (repo CertificationRepository) Create(ctx context.Context, certification models.Certification) (models.Certification, error) {
	db := repo.db.WithContext(ctx)
	result := db.Create(&certification)
	if result.Error != nil {
		repo.logger.Error("Unable to create certification", zap.Any("cert", certification), zap.Error(result.Error))
		return models.Certification{}, result.Error
	}

	repo.fillCertificationVideoUrl(ctx, &certification)
	return certification, nil
}

func (repo CertificationRepository) GetByID(ctx context.Context, id uint) (models.Certification, error) {
	db := repo.db.WithContext(ctx)
	var cert models.Certification

	res := db.Where("id = ?", id).Preload("User").First(&cert)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.Certification{}, contracts.ErrCertificationNotFound
		}

		repo.logger.Error("Unable to get certification by id", zap.Error(res.Error), zap.Any("certID", id))
		return models.Certification{}, res.Error
	}

	repo.fillCertificationVideoUrl(ctx, &cert)
	return cert, nil
}

func (repo CertificationRepository) Get(ctx context.Context, request certifications.GetCertificationsRequest) (certifications.GetCertificationsResponse, error) {
	db := repo.db.WithContext(ctx)
	var res []models.Certification

	if request.Status != "" {
		db = db.Where("status = ?", request.Status)
	}

	if request.UserID != "" {
		db = db.Where("user_id = ?", request.UserID)
	}

	result := db.Scopes(database.Paginate(res, &request.Pagination, db)).Preload("User").Find(&res)
	if result.Error != nil {
		repo.logger.Error("Unable to get certifications", zap.Any("req", request), zap.Error(result.Error))
		return certifications.GetCertificationsResponse{}, result.Error
	}

	for i := range res {
		repo.fillCertificationVideoUrl(ctx, &res[i])
	}

	return certifications.GetCertificationsResponse{
		Certifications: res,
		Pagination:     request.Pagination,
	}, nil
}

func (repo CertificationRepository) Update(ctx context.Context, certification models.Certification) (models.Certification, error) {
	db := repo.db.WithContext(ctx)

	res := db.Save(&certification)
	if res.Error != nil {
		repo.logger.Error("Unable to update certification", zap.Any("cert", certification), zap.Error(res.Error))
		return models.Certification{}, res.Error
	}

	repo.fillCertificationVideoUrl(ctx, &certification)
	return certification, nil
}

func (repo CertificationRepository) fillCertificationVideoUrl(ctx context.Context, cert *models.Certification) {
	if cert == nil {
		return
	}
	cert.VideoUrl = repo.firebase.GetCertificationVideoUrl(ctx, cert.UserID)
}
