package matching

import (
	"github.com/shopspring/decimal"
)

type PriceLevel struct {
	Price  decimal.Decimal
	Orders []Order
}

func InitializePriceLevel(price decimal.Decimal) (pl PriceLevel) {
	pl.Price = price
	pl.Orders = []Order{}
	return
}

func (pl *PriceLevel) Top() Order {
	return pl.Orders[0]
}

func (pl *PriceLevel) IsEmpty() (result bool) {
	if len(pl.Orders) == 0 {
		result = true
	}
	return
}

func (pl *PriceLevel) Add(order Order) {
	pl.Orders = append(pl.Orders, order)
	return
}

func (pl *PriceLevel) Remove(order Order) {
	var orders []Order
	for _, o := range pl.Orders {
		if o.Id != order.Id {
			orders = append(orders, o)
		}
	}
	pl.Orders = orders
	return
}

func (pl *PriceLevel) Find(id int) (order Order) {
	for _, o := range pl.Orders {
		if o.Id == id {
			order = o
			break
		}
	}
	return
}
