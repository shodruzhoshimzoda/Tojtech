package dto

import (
	"github.com/google/uuid"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

type ProductDTO struct {
	UUID        uuid.UUID         `json:"uuid"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	Category    *CategoryDTO      `json:"category,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Images      []ProductImageDTO `json:"images,omitempty"`
}

type ProductImageDTO struct {
	UUID     uuid.UUID `json:"uuid"`
	ImageURL string    `json:"image_url"`
	IsMain   bool      `json:"is_main"`
}

func NewProductDTO(p *domain_product.Product) ProductDTO {
	imagesDTO := make([]ProductImageDTO, 0, len(p.Images))
	for _, img := range p.Images {
		imagesDTO = append(imagesDTO, ProductImageDTO{
			UUID:     img.UUID,
			ImageURL: img.ImageURL,
			IsMain:   img.IsMain,
		})
	}

	var categoryDTO *CategoryDTO
	if p.Category != nil {
		categoryDTO = &CategoryDTO{
			UUID: p.Category.UUID,
		}
	}

	return ProductDTO{
		UUID:        p.UUID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    categoryDTO,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
		Images:      imagesDTO,
	}
}
