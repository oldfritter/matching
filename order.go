package matching

import (
	"github.com/shopspring/decimal"
)

type Order struct {
	Id            int
	MarketId      int
	Timestamp     int64
	Type          string
	OrderType     string
	Price         decimal.Decimal
	Volume        decimal.Decimal
	Locked        decimal.Decimal
	BasePrecision int
}

func InitializeOrder(attrs map[string]string) (order Order, err error) {
	order.Id, _ = strconv.Atoi(attrs["id"])
	order.MarketId, _ = strconv.Atoi(attrs["market_id"])
	order.Timestamp, _ = strconv.ParseInt(attrs["timestamp"], 10, 64)
	order.Type = attrs["type"]
	order.OrderType = attrs["order_type"]
	order.Volume, _ = decimal.NewFromString(attrs["volume"])
	order.Price, _ = decimal.NewFromString(attrs["price"])
	order.Locked, _ = decimal.NewFromString(attrs["locked"])
	order.BasePrecision, _ = strconv.Atoi(attrs["base_precision"])
	return
}

func (order *Order) TradeWith(counterOrder Order, counterBook OrderBook) (trade Trade) {
	if counterOrder.OrderType == "LimitOrder" {
		if order.IsCrossed(counterOrder.Price) {
			trade.Price = counterOrder.Price
			trade.Volume = order.Volume
			if order.Volume > counterOrder.Volume {
				trade.Volume = counterOrder.Volume
			}
			trade.Funds = trade.Price.Mul(trade.Volume)
		}
	} else {
		trade.Volume = order.Volume
		if trade.Volume > counterOrder.Volume {
			trade.Volume = counterOrder.Volume
		}
		volumeLimit := counterOrder.VolumeLimit(order.Price)
		if trade.Volume > volumeLimit {
			trade.Volume = volumeLimit
		}
	}
	return
}

func (order *Order) IsValid() (result bool) {
	if order.Type != "ask" && order.Type != "bid" {
		return
	}
	zero := decimal.NewFromFloat(0.0)
	if order.OrderType == "LimitOrder" && (order.Price.LessThanOrEqual(zero) || order.Volume.LessThanOrEqual(zero)) {
		return
	}
	if order.OrderType == "MarketOrder" && order.Price.GreaterThan(zero) && order.Locked.LessThanOrEqual(zero) {
		return
	}
	if order.Id > 0 && order.Timestamp > 0 && order.MarketId > 0 {
		result = true
	}
	return
}

func (order *Order) isFilled() (result bool) {
	if order.OrderType == "LimitOrder" {
		result = order.Volume.LessThanOrEqual(decimal.NewFromFloat(0.0))
	} else if order.OrderType == "MarketOrder" {
		result = order.Volume.LessThanOrEqual(decimal.NewFromFloat(0.0)) || order.Locked.LessThanOrEqual(decimal.NewFromFloat(0.0))
	}
	return
}

func (order *Order) Fill(trade Trade) {
	if trade.Volume > order.Volume {
		return
	}
	if order.OrderType == "LimitOrder" {
		return
	}
	funds := trade.Funds
	if order.Type == "ask" {
		funds = trade.Volume
	}
	if funds > order.Locked {
		return
	}
	order.Locked -= funds
	return
}

// func for MarketOrder
func (order *Order) VolumeLimit(tradePrice decimal.Decimal) (result decimal.Decimal) {
	if order.Type == "ask" {
		result = order.Locked
	} else {
		result = order.Locked.DivRound(tradePrice, int32(order.BasePrecision))
	}
	return
}

// func for LimitOrder
func (order *Order) IsCrossed(price decimal.Decimal) (result bool) {
	if order.Type == "ask" {
		result = price.GreaterThanOrEqual(order.Price)
	} else {
		result = price.LessThanOrEqual(order.Price)
	}
	return
}

func (order *Order) Label() (label string) {
	if order.OrderType == "LimitOrder" {
		label = fmt.Sprintf("%d/$%s/%s", order.Id, order.Price.String(), order.Volume.String())
	} else if order.OrderType == "MarketOrder" {
		label = fmt.Sprintf("%d/%s", order.Id, order.Volume.String())
	}
	return
}
