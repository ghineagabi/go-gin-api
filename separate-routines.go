package main

import (
	"sync"
	"time"
)

var mutexSession = &sync.RWMutex{}
var mutexVerification = &sync.RWMutex{}

func emptyVerificationCodes(ttl *map[string]VerificationTTL) {
	mutexVerification.Lock()
	for code, verCode := range *ttl {
		if verCode.expired() {
			delete(*ttl, code)
		}
	}
	mutexVerification.Unlock()
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
