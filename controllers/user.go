package controllers

import (
	"example/web-service-gin/errors"
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
)

func DeleteUser(emailID *int) error {
	sqlStatement := `
DELETE FROM public.abstract_users
WHERE id = $1;`
	res, err := utils.Db.Exec(sqlStatement, *emailID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func UpdateAbstractUser(u *models.AbstractUserToUpdate, emailID *int) error {
	var query strings.Builder
	params := make([]interface{}, 0)
	params = append(params, emailID)
	query.WriteString("UPDATE public.abstract_users SET")

	if u.FirstName != "" {
		query.WriteString(fmt.Sprintf(" first_name=$%d,", len(params)+1))
		params = append(params, u.FirstName)
	}

	if u.LastName != "" {
		query.WriteString(fmt.Sprintf(" last_name=$%d,", len(params)+1))
		params = append(params, u.LastName)
	}
	if len(params) < 2 {
		return errors.New(errors.NotEnoughParameters)
	}
	queryString := fmt.Sprintf("%s WHERE id=$1", strings.TrimSuffix(query.String(), ","))

	_, err := utils.Db.Exec(queryString, params...)
	if err != nil {
		return err
	}
	return nil
}

func InsertAbstractUser(absUsr *models.AbstractUser) error {
	sqlStatement := `INSERT INTO public.abstract_users (first_name, last_name, password, email, date_joined, last_login) 
VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := utils.Db.Exec(sqlStatement, absUsr.FirstName, absUsr.LastName,
		utils.SHA512(absUsr.Password), absUsr.Email, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()))
	if err != nil {
		return err
	}

	return nil
}

func EmailExists(email string) error {
	sqlStatement := `SELECT email FROM public.abstract_users WHERE abstract_users.email = $1
					 LIMIT 1`
	rows, err := utils.Db.Query(sqlStatement, email)
	defer rows.Close()
	if err != nil {
		return err
	}
	if rows.Next() {
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "email", Reason: "duplicated email"}
	}
	return nil
}

func UpdatePassword(emailID *int, newPass *string) error {
	sqlStatement := `UPDATE public.abstract_users
	SET password = $1
	WHERE id = $2;`
	_, err := utils.Db.Exec(sqlStatement, utils.SHA512(*newPass), *emailID)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePasswordByEmail(email *string, newPass *string) error {
	sqlStatement := `UPDATE public.abstract_users
	SET password = $1
	WHERE id = $2;`
	_, err := utils.Db.Exec(sqlStatement, utils.SHA512(*newPass), *email)
	if err != nil {
		return err
	}

	return nil
}

func GetEmailByEmailID(emailID *int, email *string) error {
	sqlStatement := `SELECT email FROM public.abstract_users WHERE abstract_users.id = $1
					 LIMIT 1`
	rows, err := utils.Db.Query(sqlStatement, *emailID)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(email); err != nil {
			return err
		}
		return nil
	}

	return &utils.InvalidFieldsError{Location: "", AffectedField: "email", Reason: "emailID not existent"}
}
