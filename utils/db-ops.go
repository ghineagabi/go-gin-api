package utils

import (
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
	row, err := Db.Query(sqlStatement, u.Email, SHA512(u.Pass))
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
		return &InvalidFieldsError{Location: "uri", AffectedField: "Email", Reason: "No Email found"}
	}
	return nil
}
