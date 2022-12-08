package utils

import (
	"github.com/lib/pq"
	"sync"
	"time"
)

var MutexSession = &sync.RWMutex{}
var MutexVerification = &sync.RWMutex{}

func emptyVerificationCodes(ttl *map[string]VerificationTTL) {
	MutexVerification.Lock()
	for code, verCode := range *ttl {
		if verCode.Expired() {
			delete(*ttl, code)
		}
	}
	MutexVerification.Unlock()
}

func EmptyVerificationCodesRoutine(seconds uint32) {
	for {
		emptyVerificationCodes(&CodeToTTL)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func ClearExpiredSessions(seconds uint32) {
	for {
		EmptyDBSessions()
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func EmptyDBSessions() {
	sqlStatement := `DELETE FROM public.sessions
					 WHERE "end" < $1;`
	_, _ = Db.Exec(sqlStatement, pq.FormatTimestamp(time.Now()))
}
