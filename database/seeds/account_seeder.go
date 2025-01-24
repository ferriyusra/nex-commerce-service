package seeds

import (
	"nex-commerce-service/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Account(db *gorm.DB) {
	accounts := []model.Account{
		{
			UserID:  1,
			Balance: 0,
		},
		{
			UserID:  2,
			Balance: 0,
		},
	}

	if err := db.Table("account").Create(&accounts).Error; err != nil {
		log.Fatal().Err(err).Msg("error seed accounts")
		return
	}

	log.Info().Msg("seed accounts success")
}
