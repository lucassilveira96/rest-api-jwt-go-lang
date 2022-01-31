package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ProductCategory struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Description string    `gorm:"size:255;not null;unique" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *ProductCategory) Prepare() {
	p.ID = 0
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ProductCategory) Validate() error {

	if p.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

func (p *ProductCategory) SaveProductCategory(db *gorm.DB) (*ProductCategory, error) {
	var err error
	err = db.Debug().Model(&ProductCategory{}).Create(&p).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	return p, nil
}

func (p *ProductCategory) FindAllProductCategory(db *gorm.DB) (*[]ProductCategory, error) {
	var err error
	productCategories := []ProductCategory{}
	err = db.Debug().Model(&ProductCategory{}).Limit(100).Find(&productCategories).Error
	if err != nil {
		return &[]ProductCategory{}, err
	}

	return &productCategories, nil
}

func (p *ProductCategory) FindProductCategoryByID(db *gorm.DB, pid uint64) (*ProductCategory, error) {
	var err error
	err = db.Debug().Model(&ProductCategory{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	return p, nil
}

func (p *ProductCategory) UpdateAProductCategory(db *gorm.DB) (*ProductCategory, error) {

	var err error
	err = db.Debug().Model(&ProductCategory{}).Where("id = ?", p.ID).Updates(ProductCategory{Description: p.Description, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	return p, nil
}

func (p *ProductCategory) DeleteAProductCategory(db *gorm.DB, pid uint64) (int64, error) {

	db = db.Debug().Model(&ProductCategory{}).Where("id = ?", pid).Take(&ProductCategory{}).Delete(&ProductCategory{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product Category not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
