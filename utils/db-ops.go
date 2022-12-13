package utils

import (
	"example/web-service-gin/errors"
	"example/web-service-gin/models"
	"github.com/lib/pq"
	"time"
)

func CreateSession(emailID *int, sessID *string) error {
	sqlStatement := `WITH A AS (INSERT INTO public.sessions (id, "end", email_id) VALUES ($1, $2, $3))
					 	  UPDATE public.abstract_users SET last_login = $4 WHERE id = $3`
	_, err := Db.Exec(sqlStatement, *sessID, pq.FormatTimestamp(time.Now().Add(time.Hour*24)), *emailID,
		pq.FormatTimestamp(time.Now()))
	if err != nil {
		return err
	}
	return nil
}

func CheckCredentials(u *models.UserCredentials, e *int) error {
	sqlStatement := `SELECT id FROM public.abstract_users WHERE abstract_users.email = $1 AND 
                                                         					   abstract_users.password = $2 
                     LIMIT 1`
	row, err := Db.Query(sqlStatement, u.Email, SHA512(u.Password))
	defer row.Close()
	if err != nil {
		return err
	}
	if !row.Next() {
		return &InvalidFieldsError{Location: "Basic auth", AffectedField: "email and password", Reason: "EmailID and/or password mismatch"}
	}
	err = row.Scan(e)
	if err != nil {
		return err
	}
	return nil
}

func GetEmailIDByEmail(email *string, emailID *int) error {
	sqlStatement := `SELECT id FROM public.abstract_users WHERE abstract_users.email = $1
					 LIMIT 1`
	rows, err := Db.Query(sqlStatement, *email)
	defer rows.Close()
	if err != nil {
		return err
	}
	if rows.Next() {
		if err = rows.Scan(emailID); err != nil {
			return err
		}
	} else {
		return errors.New(errors.EmailUnfound)
	}
	return nil
}

func GetSessionsAfterRestart(mapper map[string]CachedLoginSessions) error {
	sqlStatement := `SELECT id, "end", email_id FROM public.sessions`
	rows, err := Db.Query(sqlStatement)
	defer rows.Close()
	if err != nil {
		return err
	}

	var info CachedLoginSessions
	var sessionID string
	for rows.Next() {
		if err = rows.Scan(&sessionID, &info.SessTTL, &info.EmailID); err != nil {
			return err
		}
		mapper[sessionID] = info
	}

	return nil
}

func DeleteSession(emailID *string) error {
	sqlStatement := `DELETE FROM public.sessions WHERE sessions.id = $1`
	res, err := Db.Exec(sqlStatement, *emailID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
