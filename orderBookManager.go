package matching

type OrderBookManager struct {
	MarketId     int
	AskOrderBook OrderBook
	BidOrderBook OrderBook
}

func InitializeOrderBookManager(marketId int, options map[string]string) (orderBookManager OrderBookManager) {
	orderBookManager.MarketId = marketId
	orderBookManager.AskOrderBook = InitializeOrderBook(marketId, "ask", options)
	orderBookManager.BidOrderBook = InitializeOrderBook(marketId, "bid", options)
	return
}

func (obm *OrderBookManager) GetBooks(stype string) (OrderBook, OrderBook) {
	if stype == "ask" {
		return obm.AskOrderBook, obm.BidOrderBook
	}
	return obm.BidOrderBook, obm.AskOrderBook
}
