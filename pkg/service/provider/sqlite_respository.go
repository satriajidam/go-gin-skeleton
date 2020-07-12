package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

// ProviderSQLiteModel is a SQLite database model for provider.
type ProviderSQLiteModel struct {
	ID        uint       `gorm:"column:id;PRIMARY_KEY"`
	UUID      string     `gorm:"column:uuid;UNIQUE;UNIQUE_INDEX;NOT NULL"`
	ShortName string     `gorm:"column:short_name;INDEX;NOT NULL"`
	LongName  string     `gorm:"column:long_name;NOT NULL"`
	CreatedAt time.Time  `gorm:"column:created_at;NOT NULL"`
	UpdatedAt time.Time  `gorm:"column:updated_at;NOT NULL"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// TableName sets provider table name.
func (pm *ProviderSQLiteModel) TableName() string {
	return "provider"
}

func (pm *ProviderSQLiteModel) toProvider() *domain.Provider {
	return &domain.Provider{
		UUID:      pm.UUID,
		ShortName: pm.ShortName,
		LongName:  pm.LongName,
	}
}

func providerToProviderSQLiteModel(p domain.Provider) *ProviderSQLiteModel {
	return &ProviderSQLiteModel{
		UUID:      p.UUID,
		ShortName: p.ShortName,
		LongName:  p.LongName,
	}
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates new provider repository.
func NewRepository(db *gorm.DB) domain.ProviderRepository {
	db.AutoMigrate(&ProviderSQLiteModel{})
	return &repository{db}
}

// CreateOrUpdateProvider creates new provider or updates if the existing one.
func (r *repository) CreateOrUpdateProvider(ctx context.Context, p domain.Provider) error {
	if p.UUID == "" {
		if err := r.createProvider(ctx, p); err != nil {
			log.Error(err, "Failed creating new provider")
			return err
		}
	} else {
		if err := r.updateProvider(ctx, p); err != nil {
			if gorm.IsRecordNotFoundError(err) {
				log.Warn(fmt.Sprintf("Provider with '%s' UUID doesn't exist", p.UUID))
				return domain.ErrNotFound
			}
			log.Error(err, fmt.Sprintf("Failed updating provider with '%s' UUID", p.UUID))
			return err
		}
	}
	return nil
}

func (r *repository) createProvider(ctx context.Context, p domain.Provider) error {
	return r.db.Create(&ProviderSQLiteModel{
		UUID:      uuid.New().String(),
		ShortName: p.ShortName,
		LongName:  p.LongName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}).Error
}

func (r *repository) updateProvider(ctx context.Context, p domain.Provider) error {
	return r.db.Model(&ProviderSQLiteModel{}).
		Where("uuid = ? AND deleted_at IS NULL", p.UUID).
		Updates(
			map[string]interface{}{
				"short_name": p.ShortName,
				"long_name":  p.LongName,
				"updated_at": time.Now(),
			},
		).Error
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (r *repository) DeleteProviderByUUID(ctx context.Context, uuid string) error {
	if err := r.db.Where("uuid = ?", uuid).Delete(&ProviderSQLiteModel{}).Error; err != nil {
		log.Error(err, fmt.Sprintf("Failed deleting provider with '%s' UUID", uuid))
		return err
	}
	return nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (r *repository) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var pm ProviderSQLiteModel
	if err := r.db.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warn(fmt.Sprintf("Provider with '%s' UUID doesn't exist", uuid))
			return nil, domain.ErrNotFound
		}
		log.Error(err, fmt.Sprintf("Failed getting provider with '%s' UUID", uuid))
		return nil, err
	}
	return pm.toProvider(), nil
}

// GetProviderByShortName gets a provider based on its short name.
func (r *repository) GetProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var pm ProviderSQLiteModel
	if err := r.db.Where("short_name = ? AND deleted_at IS NULL", shortName).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warn(fmt.Sprintf("Provider with '%s' short name doesn't exist", shortName))
			return nil, domain.ErrNotFound
		}
		log.Error(err, fmt.Sprintf("Failed getting provider with '%s' short name", shortName))
		return nil, err
	}
	return pm.toProvider(), nil
}

// ListProviders lists all providers.
func (r *repository) ListProviders(ctx context.Context, limit int) ([]domain.Provider, error) {
	var pms []ProviderSQLiteModel
	query := r.db
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&pms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.Error(err, "Failed getting providers")
		return nil, err
	}
	results := []domain.Provider{}
	for _, pm := range pms {
		results = append(results, domain.Provider{
			UUID:      pm.UUID,
			ShortName: pm.ShortName,
			LongName:  pm.LongName,
		})
	}
	return results, nil
}
