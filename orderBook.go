package matching

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	MarketId     int
	Side         string
	LimitOrders  rbt.Tree
	MarketOrders rbt.Tree
	Broadcast    bool
}
type Data struct {
	Action string `json:"action"`
	Order  Order  `json:"order"`
}

func InitializeOrderBook(marketId int, side string, options map[string]interface{}) (ob OrderBook) {
	ob.MarketId = marketId
	ob.Side = side
	ob.LimitOrders = rbt.NewWithStringComparator()
	ob.MarketOrders = rbt.NewWithIntComparator()
	ob.Broadcast = true
	if options["broadcast"] != nil {
		ob.Broadcast = options["broadcast"].(bool)
	}
	ob.broadcast(struct {
		Action   string
		MarketId int
		Side     string
	}{Action: "new", MarketId: marketId, Side: side})
}

func (ob *OrderBook) BestLimitPrice() decimal.Decimal {
	order := ob.LimitTop()
	if order == nil {
		return decimal.NewFromFloat(0.0)
	}
	return order.Price
}

func (ob *OrderBook) Top() (order Order) {
	if ob.MarketOrders.Empty() {
		return ob.LimitTop()
	}
	_, order = ob.MarketOrders.Left()
	return order
}

func (ob *OrderBook) FillTop(trade Trade) {
	order = ob.Top()
	if order == nil {
		return
	}
	order.Fill(trade)
	if order.isFilled() {
		ob.Remove(order)
	} else {
		broadcast(Data{Action: "update", Order: order})
	}
}

func (ob *OrderBook) Find(order Order) (order Order) {
	if order.OrderType == "LimitOrder" {
		order = ob.LimitOrders.Get(order.Price.String()).Find(order.Id)
	} else if order.OrderType == "MarketOrder" {
		order = ob.MarketOrders.Get(order.Id)
	} else {
		return
	}
	return
}

func (ob *OrderBook) LimitTop() (order Order) {
	if ob.MarketOrders.Empty() {
		return
	}
	if ob.Side == "ask" {
		price, level := ob.Left()
	} else if ob.Side == "bid" {
		price, level := ob.tree.Right()
	}
	order = level.Top()
	return
}

func (ob *OrderBook) LimitOrders() (orders map[string][]Order) {
	orders = make(map[string][]Order)
	for _, key := range ob.LimitOrders.Keys() {
		orders[key] = ob.LimitOrders.Get(key).Orders
	}
	return
}

func (ob *OrderBook) Add(order Order) error {
	if order.Volume.LessThanOrEqual(decimal.NewFromFloat(0.0)) {
		return fmt.Errorf("volume is zero")
	}
	if order.OrderType == "LimitOrder" {
		priceLevel := InitializePriceLevel()
		if ob.LimitOrders.Get(order.Price.String()) != nil {
			priceLevel = ob.LimitOrders.Get(order.Price.String()).(PriceLevel)
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

func (ob *OrderBook) Remove(order Order) (order Order, err error) {
	if order.OrderType == "LimitOrder" {
		order = ob.removeLimitOrder(order)
	} else if order.OrderType == "MarketOrder" {
		order = ob.removeMarketOrder(order)
	} else {
		err = fmt.Errorf("Unknown order type")
	}
	return
}

func (ob *OrderBook) removeLimitOrder(order Order) (order Order) {
	priceLevel = ob.LimitOrders.Get(order.Price.String()).(PriceLevel)
	if priceLevel == nil {
		return
	}
	order = priceLevel.Find(order.Id)
	if order == nil {
		return
	}
	priceLevel.Remove(order)
	if priceLevel.Empty() {
		ob.LimitOrders.Remove(order.Price.String())
	}
	ob.broadcast(Data{Action: "remove", Order: order})
	return
}

func (ob *OrderBook) removeMarketOrder(order Order) (order Order) {
	order = ob.MarketOrders.Get(order.id)
	if order == nil {
		return
	}
	ob.MarketOrders.Remove(order.Id)
	ob.broadcast(Data{Action: "remove", Order: order})
	return
}

func (ob *OrderBook) broadcast(data interface{}) {

}
