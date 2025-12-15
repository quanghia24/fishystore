namespace go fishystore.stock

include "base.thrift"

struct GetItemStockReq {
  1: required i64 item_id
}

struct GetItemStockResp {
  1: i64 stock

  255: base.BaseResp baseResp
}

service StockService {
  GetItemStockResp GetItemStock(1: GetItemStockReq req)
}
