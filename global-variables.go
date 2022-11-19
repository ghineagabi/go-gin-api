package main

import "database/sql"

var codeToTTL = make(map[string]VerificationTTL)
var sessionToEmailID = make(map[string]CachedLoginSessions)

var db *sql.DB
var err error
var cred FileCredentials
