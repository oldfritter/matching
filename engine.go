package matching

import (
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

func InitializeEngine(marketId int, options Options) Engine {
	engine := Engine{MarketId: marketId, OrderBookManager: OrderBookManager{MarketId: marketId}, Options: options}
}

func (engine *Engine) AskOrderBook() OrderBook {
	return engine.OrderBookManager.AskOrderBook
}

func (engine *Engine) BidOrderBook() OrderBook {
	return engine.OrderBookManager.BidOrderBook
}

func (engine *Engine) Submit(order Order) {
	book, counterBook = engine.OrderBookManager.GetBooks(order.Type)
	engine.matchLimitOrder(order, counterBook)
	engine.addOrCancelLimitOrder(order, book)
}

func (engine *Engine) Cancel(order Order) {
	book, counterBook = engine.OrderBookManager.GetBooks(order.Type)
	removedOrder = book.Remove(order)
	if removedOrder != nil {
		engin.publishCancel(removedOrder, "cancelled by user")
	} else {
		// Matching.logger.warn "Cannot find order##{order.id} to cancel, skip."
	}

}

func (engine *Engine) match(order Order, counterBook OrderBook) {
	if order.Filled() || engine.isTiny(order) {
		return
	}
	counterOrder = counterBook.Top()
	if counterOrder == nil {
		return
	}
	trade := order.TradeWith(counterOrder, counterBook)
	if trade == nil {
		return
	}
	counterBook.FillTop(trade)
	order.Fill(trade)
	engine.publish(order, counterOrder, trade)
	engine.match(order, counterBook)
}

func (engine *Engine) addOrCancel(order Order, book OrderBook) {
	if order.Filled() {
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
	cas := 1
	for fixed > 0 {
		cas *= 10
		fixed--
	}
	minVolume := decimal.NewFromFloat(1.0).Div(cas)
	return order.Volume < minVolume
}

// TODO
func (engine *Engine) LimitOrders()

// TODO
func (engine *Engine) MarketOrders()
