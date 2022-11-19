package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func deleteUser(email string) error {
	sqlStatement := `
DELETE FROM public.abstract_users
WHERE email = $1;`
	res, err := db.Exec(sqlStatement, email)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func updateAbstractUser(u *AbstractUser, emailID *int) error {
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
		return &InvalidFieldsError{location: "Body", affectedField: "firstName/lastName/age",
			reason: "Could not map any of the provided fields"}
	}
	queryString := fmt.Sprintf("%s WHERE email_id=$1", strings.TrimSuffix(query.String(), ","))

	_, err = db.Exec(queryString, params...)
	if err != nil {
		return err
	}
	return nil
}

func insertAbstractUser(absUsr *AbstractUser) error {
	sqlStatement := `INSERT INTO public.abstract_users (age, first_name, last_name, password, email, date_joined, last_login) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = db.Exec(sqlStatement, absUsr.Age, absUsr.FirstName, absUsr.LastName,
		SHA512(absUsr.Password), absUsr.Email, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()))
	if err != nil {
		return err
	}

	return nil
}

func insertPost(post *PostToCreate, emailID *int) error {
	sqlStatement := `INSERT INTO public.posts (email_id, date_updated, date_created, title, content) 
VALUES ($1, $2, $3, $4, $5);`

	_, err = db.Exec(sqlStatement, emailID, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()),
		post.Title, post.Content)
	if err != nil {
		return err
	}

	return nil
}

func findPostByPostTitle(users *[]PostToGet, s *string) error {
	sqlStatement := `SELECT	concat_ws (' ', abstract_users.last_name, abstract_users.first_name) as full_name, 
       posts.title, posts.content
FROM public.posts
INNER JOIN public.abstract_users
ON posts.email_id = abstract_users.id
WHERE LOWER(posts.title) LIKE '%' || $1 || '%'`

	if strings.TrimSpace(*s) == "" {
		return &InvalidFieldsError{location: "query param", affectedField: "title", reason: "empty field"}
	}
	rows, err := db.Query(sqlStatement, s)
	defer rows.Close()

	if err != nil {
		return err
	}
	var user PostToGet
	for rows.Next() {
		if err = rows.Scan(&user.FullName, &user.Title, &user.Content); err != nil {
			return err
		}
		*users = append(*users, user)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func findPostTitles(titles *[]string, s *string, limit int) error {
	sqlStatement := `SELECT	DISTINCT posts.title
FROM public.posts
INNER JOIN public.abstract_users
ON posts.email_id = abstract_users.id
WHERE LOWER(posts.title) LIKE '%' || $1 || '%'
LIMIT ` + strconv.Itoa(limit)

	if strings.TrimSpace(*s) == "" {
		return &InvalidFieldsError{location: "query param", affectedField: "title", reason: "empty field"}
	}
	rows, err := db.Query(sqlStatement, s)
	defer rows.Close()

	if err != nil {
		return err
	}
	var title string
	for rows.Next() {
		if err = rows.Scan(&title); err != nil {
			return err
		}
		*titles = append(*titles, title)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func updatePost(post *PostToUpdate, emailID *int) error {
	var query strings.Builder
	params := []any{post.Id, emailID, pq.FormatTimestamp(time.Now())}
	query.WriteString("UPDATE public.posts SET date_updated=$3")

	if post.Title != "" {
		query.WriteString(fmt.Sprintf(", title=$%d", len(params)+1))
		params = append(params, post.Title)
	}

	if post.Content != "" {
		query.WriteString(fmt.Sprintf(", content=$%d", len(params)+1))
		params = append(params, post.Content)
	}
	if len(params) < 4 {
		return &InvalidFieldsError{location: "Body", affectedField: "title/content",
			reason: "Could not map any of the provided fields"}
	}
	query.WriteString(fmt.Sprintf(" WHERE posts.id = $1 AND posts.email_id = $2"))

	res, err := db.Exec(query.String(), params...)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return &InvalidFieldsError{location: "Body", affectedField: "id/email",
			reason: "Could not change the specified post ID. Wrong ID/email combination"}
	}
	return nil
}

func createSession(emailID *int, sessID *string) error {
	sqlStatement := `INSERT INTO public.sessions (id, "end", email_id) VALUES ($1, $2, $3);`
	_, err = db.Exec(sqlStatement, *sessID, pq.FormatTimestamp(time.Now().Add(time.Hour*24)), *emailID)
	if err != nil {
		return err
	}
	return nil
}

func checkCredentials(u *UserCredentials) error {
	sqlStatement := `SELECT email, password FROM public.abstract_users WHERE abstract_users.email = $1 AND 
                                                         					   abstract_users.password = $2 
                     LIMIT 1`
	row, err := db.Query(sqlStatement, u.Email, SHA512(u.Pass))
	defer row.Close()
	if err != nil {
		return err
	}
	if !row.Next() {
		return &InvalidFieldsError{location: "Basic auth", affectedField: "email and password", reason: "EmailID and/or password mismatch"}
	}
	return nil
}

func emailExists(email string) error {
	sqlStatement := `SELECT email FROM public.abstract_users WHERE abstract_users.email = $1
					 LIMIT 1`
	rows, err := db.Query(sqlStatement, email)
	defer rows.Close()
	if err != nil {
		return err
	}
	if rows.Next() {
		return &InvalidFieldsError{location: "Body", affectedField: "email", reason: "duplicated email"}
	}
	return nil
}

func emptyDBSessions() {
	sqlStatement := `DELETE FROM public.sessions
					 WHERE "end" < $1;`
	_, err = db.Exec(sqlStatement, pq.FormatTimestamp(time.Now()))
}

func deletePost(emailID *int, postID *int) (int, error) {
	sqlStatement := `DELETE FROM public.posts WHERE posts.id = $1 AND 
                                        					 posts.email_id = $2;`
	res, err := db.Exec(sqlStatement, postID, emailID)
	if err != nil {
		return http.StatusBadRequest, err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return http.StatusUnauthorized, &InvalidFieldsError{location: "Body", affectedField: "id",
			reason: "Could not perform action on the specified post ID"}
	}
	return http.StatusAccepted, nil
}

func likePost(emailID *int, postID *int) error {
	sqlStatement := `INSERT INTO public.likes (post_id, liked_by) VALUES ($1, $2);`
	res, err := db.Exec(sqlStatement, postID, emailID)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return &InvalidFieldsError{location: "Body", affectedField: "id",
			reason: "Could not perform action on the specified post ID"}
	}

	return nil
}

func getEmailIDByEmail(email *string, emailID *int) error {
	sqlStatement := `SELECT id FROM public.abstract_users WHERE abstract_users.email = $1
					 LIMIT 1`
	rows, err := db.Query(sqlStatement, *email)
	defer rows.Close()
	if err != nil {
		return err
	}
	if rows.Next() {
		if err = rows.Scan(emailID); err != nil {
			return err
		}
	} else {
		return &InvalidFieldsError{location: "uri", affectedField: "Email", reason: "No Email found"}
	}
	return nil
}
