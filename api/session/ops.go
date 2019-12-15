package session

import (
	"fmt"
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

var sessionMap *sync.Map

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})

	return
}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // 30 min
	ss := &defs.SimpleSession{
		UserName: un,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss)

	err := dbops.InsertSession(id, ttl, un)
	if err != nil {
		return fmt.Sprintf("Error of GenerateNewSessionId: %s", err)
	}
	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		return ss.(*defs.SimpleSession).UserName, false
	} else {
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}
		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		sessionMap.Store(sid, ss)
		return ss.UserName, false
	}

	//return "", true
}
