package dto

import (
	"time"

	domain_brand "github.com/shodruzhoshimzoda/tojtech/internal/domain/brand"
)

// BrandRequest - body of POST /brands and PATCH /brands/{uuid}.
// The admin form always sends all fields, so an update replaces them all (like categories).
type BrandRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	LogoURL     string `json:"logo_url"`
	Description string `json:"description"`
}

func (r BrandRequest) ToBrand() *domain_brand.Brand {
	return &domain_brand.Brand{
		Name:        r.Name,
		Slug:        r.Slug,
		LogoURL:     r.LogoURL,
		Description: r.Description,
	}
}

// BrandResponse - what the API returns for a brand.
type BrandResponse struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	LogoURL      string `json:"logo_url"`
	Description  string `json:"description"`
	ProductCount int    `json:"product_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func NewBrandResponse(b *domain_brand.Brand) BrandResponse {
	return BrandResponse{
		UUID:         b.UUID.String(),
		Name:         b.Name,
		Slug:         b.Slug,
		LogoURL:      b.LogoURL,
		Description:  b.Description,
		ProductCount: b.ProductCount,
		CreatedAt:    b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    b.UpdatedAt.Format(time.RFC3339),
	}
}

// NewBrandResponses never returns nil, so the JSON is [] and not null.
func NewBrandResponses(brands []*domain_brand.Brand) []BrandResponse {
	out := make([]BrandResponse, 0, len(brands))
	for _, b := range brands {
		out = append(out, NewBrandResponse(b))
	}
	return out
}