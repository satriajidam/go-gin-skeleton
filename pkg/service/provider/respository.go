package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

// ProviderSQLModel is a SQL database model for provider.
type ProviderSQLModel struct {
	ID        uint       `gorm:"column:id;PRIMARY_KEY"`
	UUID      string     `gorm:"column:uuid;UNIQUE;UNIQUE_INDEX;NOT NULL"`
	ShortName string     `gorm:"column:short_name;INDEX;NOT NULL"`
	LongName  string     `gorm:"column:long_name;NOT NULL"`
	CreatedAt time.Time  `gorm:"column:created_at;NOT NULL"`
	UpdatedAt time.Time  `gorm:"column:updated_at;NOT NULL"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// TableName sets provider table name.
func (pm *ProviderSQLModel) TableName() string {
	return "provider"
}

func (pm *ProviderSQLModel) toProvider() *domain.Provider {
	return &domain.Provider{
		UUID:      pm.UUID,
		ShortName: pm.ShortName,
		LongName:  pm.LongName,
	}
}

type repository struct {
	conn *sql.Connection
}

// NewRepository creates new provider repository.
func NewRepository(conn *sql.Connection, automigrate bool) domain.ProviderRepository {
	if automigrate {
		conn.DB.AutoMigrate(&ProviderSQLModel{})
	}
	return &repository{conn}
}

// CreateProvider creates new provider in the database.
func (r *repository) CreateProvider(ctx context.Context, p domain.Provider) error {
	if err := r.conn.DB.Create(&ProviderSQLModel{
		UUID:      p.UUID,
		ShortName: p.ShortName,
		LongName:  p.LongName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}).Error; err != nil {
		r.conn.LogError(err, "Failed creating new provider")
		return err
	}
	return nil
}

// UpdateProvider updates the existing provider in the database.
func (r *repository) UpdateProvider(ctx context.Context, p domain.Provider) error {
	if err := r.conn.DB.Model(&ProviderSQLModel{}).
		Where("uuid = ? AND deleted_at IS NULL", p.UUID).Updates(
		map[string]interface{}{
			"short_name": p.ShortName,
			"long_name":  p.LongName,
			"updated_at": time.Now(),
		},
	).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.ErrNotFound
		}
		r.conn.LogError(err, fmt.Sprintf("Failed updating provider with '%s' UUID", p.UUID))
		return err
	}
	return nil
}

// DeleteProviderByUUID deletes existing provider in the database based on its UUID.
func (r *repository) DeleteProviderByUUID(ctx context.Context, uuid string) error {
	if err := r.conn.DB.Where("uuid = ?", uuid).Delete(&ProviderSQLModel{}).Error; err != nil {
		r.conn.LogError(err, fmt.Sprintf("Failed deleting provider with '%s' UUID", uuid))
		return err
	}
	return nil
}

// GetProviderByUUID gets a provider in the database based on its UUID.
func (r *repository) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var pm ProviderSQLModel
	if err := r.conn.DB.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, domain.ErrNotFound
		}
		r.conn.LogError(err, fmt.Sprintf("Failed getting provider with '%s' UUID", uuid))
		return nil, err
	}
	return pm.toProvider(), nil
}

// GetProviderByShortName gets a provider in the database based on its short name.
func (r *repository) GetProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var pm ProviderSQLModel
	if err := r.conn.DB.Where("short_name = ? AND deleted_at IS NULL", shortName).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, domain.ErrNotFound
		}
		r.conn.LogError(err, fmt.Sprintf("Failed getting provider with '%s' short name", shortName))
		return nil, err
	}
	return pm.toProvider(), nil
}

// GetProviders gets all providers in the database.
func (r *repository) GetProviders(ctx context.Context, offset, limit int) ([]domain.Provider, error) {
	var pms []ProviderSQLModel
	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 1
	}
	if err := r.conn.DB.Offset(offset).Limit(limit).Find(&pms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		r.conn.LogError(err, "Failed getting providers")
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
