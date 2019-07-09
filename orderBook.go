package matching

import (
	"fmt"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	MarketId     int
	Side         string
	LimitOrders  *rbt.Tree
	MarketOrders *rbt.Tree
	Broadcast    bool
}
type Data struct {
	Action string `json:"action"`
	Order  Order  `json:"order"`
}

func InitializeOrderBook(marketId int, side string, options map[string]string) (ob OrderBook) {
	ob.MarketId = marketId
	ob.Side = side
	ob.LimitOrders = rbt.NewWithStringComparator()
	ob.MarketOrders = rbt.NewWithIntComparator()
	ob.Broadcast = true
	if options["broadcast"] != "" {
		if options["broadcast"] == "true" || options["broadcast"] == "True" || options["broadcast"] == "TRUE" {
			ob.Broadcast = true
		} else {
			ob.Broadcast = false
		}
	}
	ob.broadcast(struct {
		Action   string
		MarketId int
		Side     string
	}{Action: "new", MarketId: marketId, Side: side})
	return
}

func (ob *OrderBook) BestLimitPrice() decimal.Decimal {
	order := ob.LimitTop()
	if order.Id == 0 {
		return decimal.NewFromFloat(0.0)
	}
	return order.Price
}

func (ob *OrderBook) Top() (order Order) {
	if ob.MarketOrders.Empty() {
		return ob.LimitTop()
	}
	order = ob.MarketOrders.Left().Value.(Order)
	return order
}

func (ob *OrderBook) FillTop(trade Trade) {
	order := ob.Top()
	if order.Id == 0 {
		return
	}
	order.Fill(trade)
	if order.IsFilled() {
		ob.Remove(order)
	} else {
		ob.broadcast(Data{Action: "update", Order: order})
	}
}

func (ob *OrderBook) Find(order Order) (o Order) {
	if order.OrderType == "LimitOrder" {
		values, _ := ob.LimitOrders.Get(order.Price.String())
		priceLevel := values.(PriceLevel)
		o = priceLevel.Find(order.Id)
	} else if order.OrderType == "MarketOrder" {
		value, _ := ob.MarketOrders.Get(order.Id)
		o = value.(Order)
	} else {
		return
	}
	return
}

func (ob *OrderBook) LimitTop() (order Order) {
	if ob.MarketOrders.Empty() {
		return
	}
	var level PriceLevel
	if ob.Side == "ask" {
		level = ob.LimitOrders.Left().Value.(PriceLevel)
	} else if ob.Side == "bid" {
		level = ob.LimitOrders.Right().Value.(PriceLevel)
	}
	order = level.Top()
	return
}

func (ob *OrderBook) LimitOrdersMap() (orders map[string][]Order) {
	orders = make(map[string][]Order)
	for _, key := range ob.LimitOrders.Keys() {
		values, _ := ob.LimitOrders.Get(key)
		orders[key.(string)] = values.([]Order)
	}
	return
}

func (ob *OrderBook) MarketOrdersMap() (orders []Order) {
	for _, value := range ob.MarketOrders.Values() {
		orders = append(orders, value.(Order))
	}
	return
}

func (ob *OrderBook) Add(order Order) error {
	if order.Volume.LessThanOrEqual(decimal.NewFromFloat(0.0)) {
		return fmt.Errorf("volume is zero")
	}
	if order.OrderType == "LimitOrder" {
		priceLevel := InitializePriceLevel(order.Price)
		values, found := ob.LimitOrders.Get(order.Price.String())
		if found {
			priceLevel = values.(PriceLevel)
		}
		priceLevel.Add(order)
		ob.LimitOrders.Put(order.Price.String(), priceLevel)
	} else if order.OrderType == "MarketOrder" {
		ob.MarketOrders.Put(order.Id, order)
	} else {
		return fmt.Errorf("Unknown order type")
	}
	ob.broadcast(Data{Action: "add", Order: order})
	return nil
}

func (ob *OrderBook) Remove(order Order) (o Order, err error) {
	if order.OrderType == "LimitOrder" {
		o = ob.removeLimitOrder(order)
	} else if order.OrderType == "MarketOrder" {
		o = ob.removeMarketOrder(order)
	} else {
		err = fmt.Errorf("Unknown order type")
	}
	return
}

func (ob *OrderBook) removeLimitOrder(order Order) (o Order) {
	values, found := ob.LimitOrders.Get(order.Price.String())
	if !found {
		return
	}
	priceLevel := values.(PriceLevel)
	o = priceLevel.Find(order.Id)
	if o.Id == 0 {
		return
	}
	priceLevel.Remove(order)
	if priceLevel.IsEmpty() {
		ob.LimitOrders.Remove(order.Price.String())
	}
	ob.broadcast(Data{Action: "remove", Order: order})
	return
}

func (ob *OrderBook) removeMarketOrder(order Order) (o Order) {
	value, found := ob.MarketOrders.Get(order.Id)
	if !found {
		return
	}
	o = value.(Order)
	ob.MarketOrders.Remove(order.Id)
	ob.broadcast(Data{Action: "remove", Order: order})
	return
}

func (ob *OrderBook) broadcast(data interface{}) {

}
