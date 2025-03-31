// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package main

import (
	"time"
)

type AccountInput struct {
	Name string `json:"name"`
}

type Mutation struct {
}

type Order struct {
	ID         string            `json:"id"`
	Products   []*OrderedProduct `json:"products"`
	TotalPrice float64           `json:"totalPrice"`
	CreatedAt  time.Time         `json:"createdAt"`
}

type OrderInput struct {
	AccountID string               `json:"accountId"`
	Products  []*OrderProductInput `json:"products"`
}

type OrderProductInput struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type PaginationInput struct {
	Skip int `json:"skip"`
	Take int `json:"take"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

type ProductInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

type Query struct {
}
