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

	cursor, limit := validateFields(req.Cursor, req.Limit)

	fmt.Println("Executing SELECT query")
	var rows *sql.Rows
	var err error
	if *req.Reverse {
		rows, err = database.ReadAll("SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time DESC LIMIT ?", req.Chat, cursor, limit)
	} else {
		rows, err = database.ReadAll("SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time LIMIT ?", req.Chat, cursor, limit)
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

	// find next cursor
	hasMore, nextCursor := findNextCursorIfExist(req.Chat, cursor, limit, *req.Reverse)

	if hasMore {
		resp.Code, resp.Msg, resp.Messages, resp.HasMore, resp.NextCursor = 0, "Successful read", messages, &hasMore, &nextCursor
	} else {
		resp.Code, resp.Msg, resp.Messages, resp.HasMore = 0, "Successful read", messages, &hasMore
	}
	return resp, nil
}

func validateFields(reqCursor int64, reqLimit int32) (validatedCursor int64, validatedLimit int32) {
	var DEFAULT_CURSOR_VALUE int64 = 0
	var DEFAULT_LIMIT_VALUE int32 = 10

	// ensure cursor is at least 0
	var cursor = DEFAULT_CURSOR_VALUE
	if reqCursor > 0 {
		cursor = reqCursor
	}

	// ensure limit is at least 10
	var limit = DEFAULT_LIMIT_VALUE
	if reqLimit > 10 {
		limit = reqLimit
	}

	return cursor, limit
}

func findNextCursorIfExist(reqChat string, cursor int64, limit int32, reqReverse bool) (hasMoreValue bool, nextCursorValue int64) {
	var count int32
	row, err := database.Read("SELECT COUNT(*) FROM messages WHERE chat = ? AND send_time >= ?", reqChat, cursor)
	row.Scan(&count)

	var hasMore bool = false
	var nextCursor int64 = 0

	if count > limit {
		hasMore = true
	}

	var nextCursorRow *sql.Rows

	if reqReverse {
		nextCursorRow, err = database.ReadAll("SELECT send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time DESC LIMIT 1 OFFSET ?", reqChat, cursor, limit)
	} else {
		nextCursorRow, err = database.ReadAll("SELECT send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time LIMIT 1 OFFSET ?", reqChat, cursor, limit)
	}
	if err != nil {
		return false, 0
	}

	// Iterate over the single row
	for nextCursorRow.Next() {
		nextCursorRow.Scan(&nextCursor)
	}

	return hasMore, nextCursor

}
