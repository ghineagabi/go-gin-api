package utils

import (
	"database/sql"
)

var CodeToTTL = make(map[string]VerificationTTL)
var SessionToEmailID = make(map[string]CachedLoginSessions)

var Db *sql.DB
var Cred FileCredentials
