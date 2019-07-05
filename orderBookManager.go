package matching

type OrderBookManager struct {
	MarketId     int
	AskOrderBook OrderBook
	BidOrderBook OrderBook
}

func InitializeOrderBookManager(marketId int, options map[string]interface{}) (orderBookManager OrderBookManager) {
	orderBookManager.MarketId = marketId
	orderBookManager.AskOrderBook = InitializeOrderBook(marketId, "ask", options)
	orderBookManager.BidOrderBook = InitializeOrderBook(marketId, "bid", options)
}

func (obm *OrderBookManager) GetBooks(stype string) (OrderBook, OrderBook) {
	if stype == "ask" {
		return AskOrderBook, BidOrderBook
	}
	return BidOrderBook, AskOrderBook
}
