package service

import (
	"context"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ProductService interface {
	FindAll(ctx context.Context, query entity.QueryString) ([]entity.ProductEntity, int64, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.ProductEntity, error)
	Create(ctx context.Context, req entity.ProductEntity) error
	Update(ctx context.Context, req entity.ProductEntity) error
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		productRepository: repo,
	}
}

func (p *productService) Create(ctx context.Context, req entity.ProductEntity) error {

	err := p.productRepository.Create(ctx, req)
	if err != nil {
		code = "[SERVICE] Create - 1"
		log.Errorw(code, err)
		return err
	}

	return nil

}

func (p *productService) Delete(ctx context.Context, id int64) error {
	err := p.productRepository.Delete(ctx, id)
	if err != nil {
		code = "[SERVICE] Delete - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (p *productService) GetByID(ctx context.Context, id int64) (*entity.ProductEntity, error) {
	result, err := p.productRepository.GetByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (p *productService) FindAll(ctx context.Context, query entity.QueryString) ([]entity.ProductEntity, int64, int64, error) {
	results, totalData, totalPages, err := p.productRepository.FindAll(ctx, query)
	if err != nil {
		code = "[SERVICE] Get - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPages, nil
}

func (p *productService) Update(ctx context.Context, req entity.ProductEntity) error {
	err := p.productRepository.Update(ctx, req)
	if err != nil {
		code = "[SERVICE] Update - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}
