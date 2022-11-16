package main

import (
	_ "encoding/json"
	"github.com/lib/pq"
	"strings"
	"time"
)

// TODO: Change these to return an error (currently, most of them return empty strings if no error occured).

func deleteUser(userID int) {
	sqlStatement := `
DELETE FROM users
WHERE id = $1;`
	res, err := db.Exec(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}
}

func updateUser(newFirstName string, newLastName string, id int) {
	sqlStatement := `
UPDATE users
SET first_name = $2, last_name = $3
WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, newFirstName, newLastName)
	if err != nil {
		panic(err)
	}
}

func insertUser(age int, email string, firstName string, lastName string) error {
	sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, age, email, firstName, lastName).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func findUsersByID(firstNames []string) ([]string, error) {
	sqlStatement := `SELECT first_name FROM users WHERE id = ANY($1);`
	var users []string
	rows, err := db.Query(sqlStatement, pq.Array(firstNames))
	if err != nil {
		return users, err
	}
	var u string
	for rows.Next() {
		if err = rows.Scan(&u); err != nil {
			return users, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil

}

func insertAbstractUser(absUsr *AbstractUser) error {
	sqlStatement := `INSERT INTO public."abstract-users" (age, "firstName", "lastName", password, email, "dateJoined", "lastLogin" ) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = db.Exec(sqlStatement, absUsr.Age, absUsr.FirstName, absUsr.LastName,
		absUsr.Password, absUsr.Email, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()))
	if err != nil {
		return err
	}

	return nil
}

func insertPost(post *Post) error {
	sqlStatement := `INSERT INTO public."posts" (email, "dateUpdated", "dateCreated", title, "groupName", content) 
VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = db.Exec(sqlStatement, post.Email, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()),
		post.Title, post.GroupName, post.Content)
	if err != nil {
		return err
	}

	return nil
}

func findNameByPostTitle(users *[]PostInfo, s string) error {
	sqlStatement := `SELECT "abstract-users"."firstName", "abstract-users"."lastName", posts.content
FROM public.posts
INNER JOIN public."abstract-users"
ON posts.email = "abstract-users".email
WHERE public.posts.title LIKE '%` + s + `%'`

	if strings.TrimSpace(s) == "" {
		return &InvalidFieldsError{affectedField: "title", reason: "empty field"}
	}
	rows, err := db.Query(sqlStatement)
	defer rows.Close()

	if err != nil {
		return err
	}
	var user PostInfo
	for rows.Next() {
		if err = rows.Scan(&user.FirstName, &user.LastName, &user.Content); err != nil {
			return err
		}
		*users = append(*users, user)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func updatePost(post *ToUpdatePost) error {
	sqlStatement := `UPDATE public.posts
					 SET "dateUpdated"=$1, title=$2, content=$3
					WHERE posts.id = $4`
	_, err = db.Exec(sqlStatement, pq.FormatTimestamp(time.Now()), post.Title, post.Content, post.Id)
	if err != nil {
		return err
	}
	return nil
}

func createSession(s *AbstractUserSession, uc *UserCredentials) error {
	sqlStatement := `INSERT INTO public."sessions" ("id", "start", "end", "email") VALUES ($1, $2, $3, $4);`
	_id := SHA512(uc.Email + time.Now().String())
	_, err = db.Exec(sqlStatement, _id, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now().Add(time.Hour*24)), uc.Email)
	if err != nil {
		return err
	}
	s.Id = _id
	return nil
}

func checkCredentials(u *UserCredentials) error {
	sqlStatement := `SELECT email, password FROM public."abstract-users" WHERE "abstract-users".email = $1 AND 
                                                         "abstract-users".password = $2`
	var uc UserCredentials
	err = db.QueryRow(sqlStatement, u.Email, u.Pass).Scan(&uc.Email, &uc.Pass)
	if err != nil {
		return err
	}
	return nil
}
