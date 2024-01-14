package orders

import (
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/entities"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/products"
)

type OrderFilter struct {
	Search    string `query:"search"`
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginationReq
	*entities.SortReq
}

type Order struct {
	Id           string           `db:"id" json:"id"`
	UserId       string           `db:"user_id" json:"user_id"`
	TransferSlip *TransferSlip    `db:"transfer_slip" json:"transfer_slip"`
	Products     []*ProductsOrder `db:"products" json:"products"`
	Address      string           `db:"address" json:"address"`
	Contact      string           `db:"contact" json:"contact"`
	Status       string           `db:"status" json:"status"`
	TotalPaid    float64          `db:"total_pain" json:"total_paid"`
	CreatedAt    string           `db:"created_at" json:"created_at"`
	UpdatedAt    string           `db:"updated_at" json:"updated_at"`
}

type TransferSlip struct {
	Id        string `json:"id"`
	FileName  string `json:"filename"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	Id      string            `db:"id" json:"id"`
	Qty     int               `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}
