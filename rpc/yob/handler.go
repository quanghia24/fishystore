package main

import (
	"context"

	"github.com/quanghia24/fishystore/gen-go/base"
	"github.com/quanghia24/fishystore/gen-go/yob"
)

// In-memory yob store
var (
	yobStore = map[int64]int64{
		1: 1990,
		2: 1992,
		3: 1988,
	}
)

type YOBHandler struct {
}

func NewYOBHandler() *YOBHandler {
	return &YOBHandler{}
}

func (p *YOBHandler) GetUserYOB(ctx context.Context, req *yob.GetUserYOBReq) (_r *yob.GetUserYOBResp, _err error) {
	yobValue, ok := yobStore[req.UserID]

	resp := &yob.GetUserYOBResp{
		BaseResp: &base.BaseResp{},
	}

	if !ok {
		resp.BaseResp.Code = "404"
		resp.BaseResp.Msg = "User not found"
		return resp, nil
	}

	resp.Yob = yobValue
	resp.BaseResp.Code = "200"
	resp.BaseResp.Msg = "Success"

	return resp, nil
}
