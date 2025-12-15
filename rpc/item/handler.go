package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex/client"
	item "github.com/quanghia24/fishystore/kitex_gen/fishystore/item"
	"github.com/quanghia24/fishystore/kitex_gen/fishystore/stock"
	"github.com/quanghia24/fishystore/kitex_gen/fishystore/stock/stockservice"
)

// ItemServiceImpl implements the last service interface defined in the IDL.
type ItemServiceImpl struct {
	stockCli stockservice.Client
}

func NewStockClient(addr string) (stockservice.Client, error) {
	return stockservice.NewClient("example.shop.stock", client.WithHostPorts(addr))
}

// GetItem implements the ItemServiceImpl interface.
func (s *ItemServiceImpl) GetItem(ctx context.Context, req *item.GetItemReq) (resp *item.GetItemResp, err error) {
	resp = item.NewGetItemResp()
	Item := item.NewItem()
	Item.Id = req.GetId()
	Item.Title = "Kitexxx"
	Item.Description = "description kitex"

	stockReq := stock.NewGetItemStockReq()
	stockReq.ItemId = req.GetId()
	stockResp, err := s.stockCli.GetItemStock(ctx, stockReq)
	if err != nil {
		log.Println(err)
		stockResp.Stock = 0
	}
	Item.Stock = stockResp.GetStock()
	resp.Item = Item
	return
}
