package repository

import (
	"context"
	"fmt"
	"math"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	FindAll(ctx context.Context, query entity.QueryString) ([]entity.ProductEntity, int64, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.ProductEntity, error)
	Create(ctx context.Context, req entity.ProductEntity) error
	Update(ctx context.Context, req entity.ProductEntity) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (p *productRepository) Create(ctx context.Context, req entity.ProductEntity) error {

	modelProduct := model.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
	}

	err = p.db.Create(&modelProduct).Error

	if err != nil {
		code = "[REPOSITORY] CreateProduct - 1"
		log.Errorw(code, err)
		return err
	}

	return nil

}

func (p *productRepository) Delete(ctx context.Context, id int64) error {

	err = p.db.Where("id = ?", id).Model(&model.Product{}).Delete(&model.Product{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteProduct - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (p *productRepository) Update(ctx context.Context, req entity.ProductEntity) error {

	modelProduct := model.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
	}

	err = p.db.Where("id = ?", req.ID).Updates(&modelProduct).Error
	if err != nil {
		code = "[REPOSITORY] UpdateProduct - 1"
		log.Errorw(code, err)
		return err
	}

	return nil

}

func (p *productRepository) FindAll(ctx context.Context, query entity.QueryString) ([]entity.ProductEntity, int64, int64, error) {
	var modelProduct []model.Product
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := p.db.Table("product").Preload(clause.Associations).
		Where("name ILIKE ? OR category ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")

	err = sqlMain.Model(&modelProduct).Count(&countData).Error
	if err != nil {
		code = "[REPOSITORY] FindAllProduct - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))

	err = sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelProduct).Error
	if err != nil {
		code = "[REPOSITORY] FindAllProduct - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := []entity.ProductEntity{}
	for _, val := range modelProduct {
		resp := entity.ProductEntity{
			ID:            val.ID,
			Name:          val.Name,
			Description:   val.Description,
			Price:         val.Price,
			StockQuantity: val.StockQuantity,
			Category:      val.Category,
			CreatedAt:     val.CreatedAt,
			UpdatedAt:     val.UpdatedAt}

		resps = append(resps, resp)
	}
	return resps, countData, int64(totalPages), nil

}

func (p *productRepository) GetByID(ctx context.Context, id int64) (*entity.ProductEntity, error) {

	var modelProduct model.Product
	err = p.db.Where("id = ?", id).Find(&modelProduct).Error
	if err != nil {
		code = "[REPOSITORY] GetProduct - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.ProductEntity{
		ID:            modelProduct.ID,
		Name:          modelProduct.Name,
		Description:   modelProduct.Description,
		Price:         modelProduct.Price,
		StockQuantity: modelProduct.StockQuantity,
		Category:      modelProduct.Category,
		CreatedAt:     modelProduct.CreatedAt,
		UpdatedAt:     modelProduct.UpdatedAt,
	}

	return &resp, nil

}
