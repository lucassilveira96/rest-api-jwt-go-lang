package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "administrador do sistema",
		Email:    "administrator@gmail.com",
		Password: "12345000",
	},
	models.User{
		Nickname: "Lucas Silveira",
		Email:    "lucas.silva.silveira@rede.ulbra.br",
		Password: "password",
	},
}

var product_categories = []models.ProductCategory{
	models.ProductCategory{
		ID:          1,
		Description: "acougue",
	},
	models.ProductCategory{
		ID:          2,
		Description: "padaria",
	},
}

var products = []models.Product{
	models.Product{
		Description:       "Carne de Primeira",
		Price:             45.50,
		ProductCategoryId: 1,
	},
	models.Product{
		Description:       "Pao Frances",
		Price:             9.90,
		ProductCategoryId: 2,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Product{}, &models.User{}, &models.ProductCategory{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.ProductCategory{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Product{}).AddForeignKey("product_category_id", "product_categories(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i := range product_categories {
		err = db.Debug().Model(&models.ProductCategory{}).Create(&product_categories[i]).Error
		if err != nil {
			log.Fatalf("cannot seed product_categories table: %v", err)
		}
	}

	for i := range products {
		err = db.Debug().Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed products table: %v", err)
		}
	}
}
