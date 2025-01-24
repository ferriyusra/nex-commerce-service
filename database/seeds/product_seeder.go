package seeds

import (
	"nex-commerce-service/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Product(db *gorm.DB) {
	products := []model.Product{
		{
			UserID:        1,
			Name:          "Laptop Pro",
			Description:   "High-performance laptop for professionals",
			Price:         1500000,
			StockQuantity: 10,
			Category:      "Electronics",
		},
		{
			UserID:        1,
			Name:          "Wireless Mouse",
			Description:   "Ergonomic wireless mouse",
			Price:         250000,
			StockQuantity: 50,
			Category:      "Accessories",
		},
		{
			UserID:        1,
			Name:          "Gaming Keyboard",
			Description:   "Mechanical gaming keyboard with RGB",
			Price:         500000,
			StockQuantity: 20,
			Category:      "Electronics",
		},
	}

	if err := db.Table("product").Create(&products).Error; err != nil {
		log.Fatal().Err(err).Msg("error seed products")
		return
	}

	log.Info().Msg("seed products success")
}
