package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)



// 
var (
	ErrProductNotFound = errors.New("product not  found")
)



type ProductRespository struct {
	dbPool 	*pgxpool.Pool
}




func NewProductRepository(dbPool *pgxpool.Pool) *ProductRespository {
	return &ProductRespository{
		dbPool: dbPool,
	}
}



func (p *ProductRespository) GetById(ctx context.Context, id int64)  (*domain.Product, error) {


		query := "SELECT id, name, description, price from products where id=$1;"


		var product domain.Product

		if err := p.dbPool.QueryRow(ctx, query, id).Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrProductNotFound
			}

			return nil, err
		}


		return &product, nil


}