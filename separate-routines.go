package main

import (
	"sync"
	"time"
)

var mutex = &sync.RWMutex{}

func emptyVerificationCodes(ttl *map[string]VerificationTTL) {
	mutex.Lock()
	for code, verCode := range *ttl {
		if verCode.expired() {
			delete(*ttl, code)
		}
	}
	mutex.Unlock()
}

func emptyVerificationCodesRoutine(seconds uint32) {
	for {
		emptyVerificationCodes(&codeToTTL)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func clearExpiredSessions(seconds uint32) {
	for {
		emptyDBSessions()
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}
