package main

import (
	"net/http"

	"github.com/labstack/echo"
//	"database/sql"
	"time"
)

func fetchUnread(c echo.Context) error {
	userID := sessUserID(c)
	if userID == 0 {
		return c.NoContent(http.StatusForbidden)
	}

	time.Sleep(time.Second)

/*
	channels, err := queryChannels()
	if err != nil {
		return err
	}
*/
/*
	rows, err := db.Query("select m.channel_id, count(m.id) from message m group by m.channel_id")
	if err != nil{
		return err
	}
	all := make(map[int64]int64)
	for rows.Next() {
		var cid, c int64
		rows.Scan(&cid, &c)
		all[cid] = c
	}
	rows, err = db.Query("select c.id, count(m.id), max(h.message_id) from channel c left join haveread h on c.id = h.channel_id and h.user_id = ? left join message m on m.channel_id = c.id and m.id > h.message_id group by c.id", userID)
	if err != nil{
		return err
	}
	for rows.Next() {
		var id, cnt int64
		var mid sql.NullInt64
		rows.Scan(&id, &cnt, &mid)
		var midv int64
		if err = mid.Scan(&midv); err != nil {
		if val, ok := all[id]; ok{
		r := map[string]interface{}{
			"channel_id": id,
			"unread":     val}
		resp = append(resp, r)
		}else{
		r := map[string]interface{}{
			"channel_id": id,
			"unread":     0}
		resp = append(resp, r)
		}
		}else{
		r := map[string]interface{}{
			"channel_id": id,
			"unread":     cnt}
		resp = append(resp, r)
		}
	}
*/
	resp := []map[string]interface{}{}
	rows, err := db.Query("select c.id, count(m.id) from channel c left join haveread h on c.id = h.channel_id and h.user_id = ? left join message m on m.channel_id = c.id and m.id > IFNULL(h.message_id,0) group by c.id", userID)
	if err != nil{
		return err
	}
	for rows.Next() {
		var id, cnt int64
		rows.Scan(&id, &cnt)
		r := map[string]interface{}{
			"channel_id": id,
			"unread":     cnt}
		resp = append(resp, r)
	}

/*
	for _, chID := range channels {
		lastID, err := queryHaveRead(userID, chID)
		if err != nil {
			return err
		}

		var cnt int64
		if lastID > 0 {
			err = db.Get(&cnt,
				"SELECT COUNT(*) as cnt FROM message WHERE channel_id = ? AND ? < id",
				chID, lastID)
		} else {
			err = db.Get(&cnt,
				"SELECT COUNT(*) as cnt FROM message WHERE channel_id = ?",
				chID)
		}
		if err != nil {
			return err
		}
		r := map[string]interface{}{
			"channel_id": chID,
			"unread":     cnt}
		resp = append(resp, r)
	}
*/


	return c.JSON(http.StatusOK, resp)
}
