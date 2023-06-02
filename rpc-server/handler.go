package main

import (
	"context"
	"strconv"
	"github.com/sebbycake/tiktok_assignment_2023/rpc-server/db"
	"github.com/sebbycake/tiktok_assignment_2023/rpc-server/kitex_gen/rpc"
	"math/rand"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()

	// Create a new database connection
	database, err := db.NewDB("admin", "password", "localhost", "instant_messaging_app")
	if err != nil {
		resp.Code, resp.Msg = 500, "Error connecting to database"
		return resp, err
	}
	defer database.Close()

	// Execute SQL insert query
	id, err := database.Create("INSERT INTO messages (chat, text, sender, send_time) VALUES (?, ?, ?, ?)", req.Message.Chat, req.Message.Text, req.Message.Sender, req.Message.SendTime)
	if err != nil {
		resp.Code, resp.Msg = 500, "Error executing INSERT query"
		return resp, err
	}
	resp.Code, resp.Msg = 0, strconv.FormatInt(id, 10)
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = areYouLucky()
	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}