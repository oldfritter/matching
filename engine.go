package matching

import (
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	DEFAULT_PRECISION = 8
)

type Engine struct {
	MarketId         int
	OrderBookManager OrderBookManager
	Options          Options
}
type Options struct {
	Id              int
	Code            string
	Name            string
	BaseUnit        string
	QuoteUnit       string
	PriceGroupFixed int
	Bid             Fee
	Ask             Fee
	SortOrder       int
}
type Fee struct {
	Fee      decimal.Decimal
	Currency string
	Fixed    int
}
type Trade struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
	Funds  decimal.Decimal
}

func (trade *Trade) isNotValidated() bool {
	return trade.Price.IsZero() || trade.Volume.IsZero()
}

func InitializeEngine(marketId int, options Options) Engine {
	return Engine{
		MarketId:         marketId,
		OrderBookManager: InitializeOrderBookManager(marketId, map[string]string{}),
		Options:          options,
	}
}

func (engine *Engine) AskOrderBook() OrderBook {
	return engine.OrderBookManager.AskOrderBook
}

func (engine *Engine) BidOrderBook() OrderBook {
	return engine.OrderBookManager.BidOrderBook
}

func (engine *Engine) Submit(order Order) {
	book, counterBook := engine.OrderBookManager.GetBooks(order.Type)
	engine.match(order, counterBook)
	engine.addOrCancel(order, book)
}

func (engine *Engine) Cancel(order Order) {
	book, _ := engine.OrderBookManager.GetBooks(order.Type)
	removedOrder, err := book.Remove(order)
	if err == nil {
		engine.publishCancel(removedOrder, "cancelled by user")
	} else {
		// Matching.logger.warn "Cannot find order##{order.id} to cancel, skip."
	}

}

func (engine *Engine) match(order Order, counterBook OrderBook) {
	if order.IsFilled() || engine.isTiny(order) {
		return
	}
	counterOrder := counterBook.Top()
	if counterOrder.Id == 0 {
		return
	}
	trade := order.TradeWith(counterOrder, counterBook)
	if trade.isNotValidated() {
		return
	}
	counterBook.FillTop(trade)
	order.Fill(trade)
	engine.publish(order, counterOrder, trade)
	engine.match(order, counterBook)
}

func (engine *Engine) addOrCancel(order Order, book OrderBook) {
	if order.IsFilled() {
		return
	}
	if order.OrderType == "LimitOrder" {
		book.Add(order)
	} else if order.OrderType == "MarketOrder" {
		engine.publishCancel(order, "fill or kill market order")
	}
	return
}

func (engine *Engine) publish(order, counterOrder Order, trade Trade) {
	var ask, bid Order
	if order.Type == "ask" {
		ask = order
		bid = counterOrder
	} else {
		ask = counterOrder
		bid = order
	}
	fmt.Println(ask)
	fmt.Println(bid)
	// logger 记录订单成交
	return
}

func (engine *Engine) publishCancel(order Order, reson string) {
	// logger 记录订单取消
	return
}

func (engine *Engine) isTiny(order Order) (result bool) {
	var fixed = DEFAULT_PRECISION
	if engine.Options.Ask.Fixed != 0 {
		fixed = engine.Options.Ask.Fixed
	}
	cas := decimal.NewFromFloat(1)
	for fixed > 0 {
		cas = cas.Mul(decimal.NewFromFloat(10))
		fixed--
	}
	minVolume := decimal.NewFromFloat(1.0).Div(cas)
	return order.Volume.LessThan(minVolume)
}

func (engine *Engine) LimitOrders() (result map[string]map[string][]Order) {
	askOrderBook := engine.AskOrderBook()
	bidOrderBook := engine.BidOrderBook()
	result["ask"] = askOrderBook.LimitOrdersMap()
	result["bid"] = bidOrderBook.LimitOrdersMap()
	return
}

func (engine *Engine) MarketOrders() (result map[string][]Order) {
	askOrderBook := engine.AskOrderBook()
	bidOrderBook := engine.BidOrderBook()
	result["ask"] = askOrderBook.MarketOrdersMap()
	result["bid"] = bidOrderBook.MarketOrdersMap()
	// result["ask"] = engine.AskOrderBook.MarketOrders()
	// result["bid"] = engine.BidOrderBook.MarketOrders()
	return
}
