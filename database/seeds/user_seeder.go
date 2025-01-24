package seeds

import (
	"nex-commerce-service/internal/core/domain/model"
	"nex-commerce-service/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func UserSeller(db *gorm.DB) {
	bytes, err := conv.HashPassword("seller123")
	if err != nil {
		log.Fatal().Err(err).Msg("error hashing password")
	}

	seller := model.User{
		Username: "seller",
		Email:    "seller@gmail.com",
		Role:     "seller",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&seller, model.User{
		Email: "seller@gmail.com",
	}).Error; err != nil {
		log.Fatal().Err(err).Msg("error seed seller")
	} else {
		log.Info().Msg("seed seller success")
	}
}

func UserAdmin(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("error hashing password")
	}

	admin := model.User{
		Username: "admin",
		Email:    "admin@gmail.com",
		Role:     "admin",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{
		Email: "admin@gmail.com",
	}).Error; err != nil {
		log.Fatal().Err(err).Msg("error seed admin")
	} else {
		log.Info().Msg("seed admin success")
	}
}
