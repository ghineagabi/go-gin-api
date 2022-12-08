package controllers

import (
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
)

func DeleteUser(email string) error {
	sqlStatement := `
DELETE FROM public.abstract_users
WHERE email = $1;`
	res, err := utils.Db.Exec(sqlStatement, email)
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
	if u.Age != 0 {
		query.WriteString(fmt.Sprintf(" age=$%d,", len(params)+1))
		params = append(params, u.Age)
	}
	if len(params) < 2 {
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "firstName/lastName/age",
			Reason: "Could not map any of the provided fields"}
	}
	queryString := fmt.Sprintf("%s WHERE id=$1", strings.TrimSuffix(query.String(), ","))

	_, err = utils.Db.Exec(queryString, params...)
	if err != nil {
		return err
	}
	return nil
}

func InsertAbstractUser(absUsr *models.AbstractUser) error {
	sqlStatement := `INSERT INTO public.abstract_users (age, first_name, last_name, password, email, date_joined, last_login) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = utils.Db.Exec(sqlStatement, absUsr.Age, absUsr.FirstName, absUsr.LastName,
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
