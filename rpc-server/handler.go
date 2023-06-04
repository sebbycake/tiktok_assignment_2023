package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sebbycake/tiktok_assignment_2023/rpc-server/db"
	"github.com/sebbycake/tiktok_assignment_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

var database, connectionError = db.ConnectToDB()

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()

	if connectionError != nil {
		resp.Code, resp.Msg = 500, "Error connecting to database"
		return resp, connectionError
	}

	fmt.Println("Executing INSERT query")
	// Execute SQL insert query
	id, err := database.Create("INSERT INTO messages (chat, text, sender, send_time) VALUES (?, ?, ?, ?)", req.Message.Chat, req.Message.Text, req.Message.Sender, req.Message.SendTime)
	if err != nil {
		resp.Code, resp.Msg = 500, "Error executing INSERT query"
		return resp, err
	}
	fmt.Println("Executed INSERT query")
	
	resp.Code, resp.Msg = 0, strconv.FormatInt(id, 10)
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	if connectionError != nil {
		resp.Code, resp.Msg = 500, "Error connecting to database"
		return resp, connectionError
	}

	var rows *sql.Rows
	var err error

	fmt.Println("Executing SELECT query")
	if *req.Reverse {
		rows, err = database.ReadAll("SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time DESC LIMIT ?", req.Chat, req.Cursor, req.Limit)
	} else {
		rows, err = database.ReadAll("SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time LIMIT ?", req.Chat, req.Cursor, req.Limit)
	}
	if err != nil {
		resp.Code, resp.Msg = 500, "Error making a select query"
		return resp, err
	}
	fmt.Println("Executed SELECT query")

	var messages []*rpc.Message

	fmt.Println("Processing rows")
	for rows.Next() {
		message := rpc.NewMessage()
		err := rows.Scan(&message.Chat, &message.Text, &message.Sender, &message.SendTime)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}
	fmt.Println("Finished processing rows")
	
	resp.Code, resp.Msg, resp.Messages = 0, "Successful read", messages
	return resp, nil
}
