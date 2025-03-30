package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(context.Context, *Product) error
	GetProductByID(context.Context, string) (*Product, error)
	ListProducts(context.Context, uint64, uint64) ([]*Product, error)
	ListProductsWithIDs(context.Context, []string) ([]*Product, error)
	SearchProducts(context.Context, string, uint64, uint64) ([]*Product, error)
}

type elasticRepo struct {
	Client *elastic.Client
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	Client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &elasticRepo{Client}, nil
}

func (r *elasticRepo) Close() {
	r.Client.Stop()
}

func (r *elasticRepo) PutProduct(ctx context.Context, product *Product) error {
	_, err := r.Client.Index().
		Index("catalog").
		Type("product").
		Id(product.ID).
		BodyJson(ProductDocument{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}).Do(ctx)

	return err
}

func (r *elasticRepo) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.Client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	p := &ProductDocument{}
	if err := json.Unmarshal(*res.Source, p); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepo) ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error) {
	res, err := r.Client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*Product{}
	for _, hit := range res.Hits.Hits {
		p := &ProductDocument{}

		if err := json.Unmarshal(*hit.Source, p); err == nil {
			products = append(products, &Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
			})
		}
	}

	return products, nil
}

func (r *elasticRepo) ListProductsWithIDs(ctx context.Context, ids []string) ([]*Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().
				Index("catalog").
				Type("product").
				Id(id),
		)
	}

	res, err := r.Client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for _, doc := range res.Docs {
		if doc.Found {
			p := &ProductDocument{}
			if err := json.Unmarshal(*doc.Source, p); err == nil {
				products = append(products, &Product{
					ID:          doc.Id,
					Name:        p.Name,
					Description: p.Description,
				})
			}
		}
	}

	return products, nil
}

func (r *elasticRepo) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error) {
	res, err := r.Client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*Product{}
	for _, hit := range res.Hits.Hits {
		p := &ProductDocument{}

		if err := json.Unmarshal(*hit.Source, p); err == nil {
			products = append(products, &Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
			})
		}
	}

	return products, nil
}
