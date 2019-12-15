package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

func InsertSession(sid string, ttl int64, uname string) error {
	ttlStr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare(`INSERT INTO sessions (session_id, ttl, login_name) VALUES(?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(sid, ttlStr, uname)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(`SELECT ttl, login_name FROM sessions WHERE session_id = ?`)
	if err != nil {
		return ss, err
	}

	var (
		ttl   string
		uname string
	)
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.UserName = uname
	} else {
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare(`SELECT * FROM sessions`)
	if err != nil {
		log.Printf("%s", err)
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id, ttlStr, loginName string
		err := rows.Scan(&id, &ttlStr, &loginName)
		if err != nil {
			log.Printf("retrieve sessions error: %s", err)
			break
		}

		if ttl, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
			ss := &defs.SimpleSession{
				UserName: loginName,
				TTL:      ttl,
			}
			m.Store(id, ss)
			log.Printf(" session id:  %s, ttl: %d", id, ss.TTL)
		} else {
			log.Printf("parse TTL error: %s", err)
			break
		}
	}

	return m, nil
}

func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare(`DELETE FROM sessions WHERE session_id=?`)
	if err != nil {
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	defer stmtOut.Close()
	return nil
}
