package provider

import (
	"github.com/jinzhu/gorm"
)

// ProviderModel is a database model for provider table.
type ProviderModel struct {
	UUID      string `gorm:"uuid"`
	ShortName string `gorm:"short_name"`
	LongName  string `gorm:"long_name"`
}

// TableName sets provider table name.
func (pm *ProviderModel) TableName() string {
	return "provider"
}

func (pm *ProviderModel) toProvider() Provider {
	return Provider{
		UUID:      pm.UUID,
		ShortName: pm.ShortName,
		LongName:  pm.LongName,
	}
}

// Repository provides methods for interacting with provider repository.
type Repository interface {
	CreateOrUpdateProvider(pm ProviderModel) error
	GetProviderByUUID(uuid string) (*ProviderModel, error)
	ListProviders(limit int) ([]ProviderModel, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates new provider repository.
func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&ProviderModel{})
	return &repository{db}
}

// CreateProvider creates new provider or updates if it's already exist.
func (r *repository) CreateOrUpdateProvider(pm ProviderModel) error {
	if err := r.db.Model(&pm).Where("uuid = ?", pm.UUID).Updates(
		map[string]interface{}{"short_name": pm.ShortName, "long_name": pm.LongName},
	).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return r.db.Create(&pm).Error
		}
		return err
	}
	return nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (r *repository) GetProviderByUUID(uuid string) (*ProviderModel, error) {
	var result ProviderModel
	if err := r.db.Where("uuid = ?", uuid).First(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ListProviders lists all providers.
func (r *repository) ListProviders(limit int) ([]ProviderModel, error) {
	var result []ProviderModel
	query := r.db
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
