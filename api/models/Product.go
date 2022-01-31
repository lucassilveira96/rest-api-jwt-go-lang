package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID                uint64          `gorm:"primary_key;auto_increment" json:"id"`
	Description       string          `gorm:"size:255;not null;unique" json:"description"`
	Price             float32         `gorm:"size:255;not null;" json:"price"`
	ProductCategory   ProductCategory `json:"productCategories"`
	ProductCategoryId uint32          `sql:"type:int REFERENCES product_categories(id)" json:"product_category_id"`
	CreatedAt         time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.ID = 0
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.ProductCategory = ProductCategory{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {
	fmt.Println(p.ProductCategoryId)
	if p.Description == "" {
		return errors.New("Required Description")
	}
	if p.Price < 0 {
		return errors.New("Invalid Price")
	}
	return nil
}

func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&ProductCategory{}).Where("id = ?", p.ProductCategoryId).Take(&p.ProductCategory).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i := range products {
			err := db.Debug().Model(&User{}).Where("id = ?", products[i].ProductCategoryId).Take(&products[i].ProductCategory).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.ProductCategoryId).Take(&p.ProductCategory).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) UpdateAProduct(db *gorm.DB) (*Product, error) {

	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{Description: p.Description, Price: p.Price, ProductCategoryId: p.ProductCategoryId, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&ProductCategory{}).Where("id = ?", p.ProductCategoryId).Take(&p.ProductCategory).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) DeleteAProduct(db *gorm.DB, pid uint64) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
