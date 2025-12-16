package main

import (
	"context"
	"strconv"

	"github.com/quanghia24/fishystore/gen-go/base"
	"github.com/quanghia24/fishystore/gen-go/user"
	"github.com/quanghia24/fishystore/gen-go/yob"
)

// In-memory user store
var (
	userStore = map[int64]*user.User{
		1: {
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		2: {
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
		3: {
			ID:    3,
			Name:  "Bob Johnson",
			Email: "bob@example.com",
		},
	}
)

type UserHandler struct {
	yobCli *yob.YOBServiceClient
}

func NewUserHandler(yobClient *yob.YOBServiceClient) *UserHandler {
	return &UserHandler{
		yobCli: yobClient,
	}
}

func (p *UserHandler) GetUser(ctx context.Context, req *user.GetUserReq) (_r *user.GetUserResp, _err error) {
	usr, ok := userStore[req.ID]

	resp := &user.GetUserResp{
		BaseResp: &base.BaseResp{},
	}

	if !ok {
		resp.BaseResp.Code = "404"
		resp.BaseResp.Msg = "User not found"
		return resp, nil
	}

	resp.User = usr

	// get year of birth
	yobReq := &yob.GetUserYOBReq{
		UserID: req.ID,
	}
	yobResp, err := p.yobCli.GetUserYOB(context.Background(), yobReq)
	if err != nil {
		return nil, err
	}
	// i64 -> string
	resp.User.Yob = strconv.Itoa(int(yobResp.Yob))

	resp.BaseResp.Code = "200"
	resp.BaseResp.Msg = "Success"

	return resp, nil
}
