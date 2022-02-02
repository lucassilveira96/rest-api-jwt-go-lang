package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Stock struct {
	ID                uint64  `gorm:"primary_key;auto_increment" json:"id"`
	MaximumQuantity   int32   `gorm:"size:255;not null;" json:"maximum_quantity"`
	MinimumQuantity   int32   `gorm:"size:255;not null;" json:"minimum_quantity"`
	AvailableQuantity int32   `gorm:"size:255;not null;" json:"avaliable_quantity"`
	UsedQuantity      int32   `gorm:"size:255;not null;" json:"used_quantity"`
	Product           Product `json:"products"`
	ProductId         uint32  `sql:"type:int REFERENCES products(id)" json:"product_id"`
	Deleted           int32
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         time.Time `gorm:"default:NULL" json:"deleted_at"`
}

func (s *Stock) Prepare() {
	s.ID = 0
	s.Product = Product{}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *Stock) Validate() error {
	return nil
}

func (s *Stock) SaveStock(db *gorm.DB) (*Stock, error) {
	var err error
	err = db.Debug().Model(&Stock{}).Create(&s).Error
	if err != nil {
		return &Stock{}, err
	}

	return s, nil
}

func (s *Stock) FindAllStocks(db *gorm.DB) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	err = db.Debug().Model(&Stock{}).Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	if len(stocks) > 0 {
		for i := range stocks {
			err := db.Debug().Model(&Product{}).Where("id = ?", stocks[i].ProductId).Take(&stocks[i].Product).Error
			if err != nil {
				return &[]Stock{}, err
			}
		}
	}
	return &stocks, nil
}

func (s *Stock) FindStockByID(db *gorm.DB, pid uint64) (*Stock, error) {
	var err error
	err = db.Debug().Model(&Stock{}).Where("id = ?", pid).Take(&s).Error
	if err != nil {
		return &Stock{}, err
	}

	if s.ID != 0 {
		err = db.Debug().Model(&Product{}).Where("id = ?", s.ProductId).Take(&s.Product).Error
		if err != nil {
			return &Stock{}, err
		}
	}
	return s, nil
}

func (s *Stock) UpdateAStock(db *gorm.DB) (*Stock, error) {

	var err error
	err = db.Debug().Model(&Stock{}).Where("id = ?", s.ID).Updates(Stock{UsedQuantity: s.UsedQuantity, MaximumQuantity: s.MaximumQuantity, AvailableQuantity: s.AvailableQuantity, MinimumQuantity: s.MinimumQuantity, ProductId: s.ProductId, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Stock{}, err
	}

	return s, nil
}

func (s *Stock) DeleteAStock(db *gorm.DB, pid uint64) (*Stock, error) {

	var err error
	err = db.Debug().Model(&Stock{}).Where("id = ?", pid).Updates(Stock{Deleted: 1, DeletedAt: time.Now()}).Error
	if err != nil {
		return &Stock{}, err
	}

	return s, nil
}
