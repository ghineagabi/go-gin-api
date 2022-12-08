package controllers

import (
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InsertPost(post *models.PostToCreate, emailID *int) error {
	sqlStatement := `INSERT INTO public.posts (email_id, date_updated, date_created, title, content) 
VALUES ($1, $2, $3, $4, $5);`

	_, err := utils.Db.Exec(sqlStatement, emailID, pq.FormatTimestamp(time.Now()), pq.FormatTimestamp(time.Now()),
		post.Title, post.Content)
	if err != nil {
		return err
	}

	return nil
}

func FindPostByPostTitle(users *[]models.PostToGet, s *string) error {
	sqlStatement := `SELECT	concat_ws (' ', abstract_users.last_name, abstract_users.first_name) as full_name, 
       posts.title, posts.content
FROM public.posts
INNER JOIN public.abstract_users
ON posts.email_id = abstract_users.id
WHERE LOWER(posts.title) LIKE '%' || $1 || '%'`

	rows, err := utils.Db.Query(sqlStatement, s)
	defer rows.Close()

	if err != nil {
		return err
	}
	var user models.PostToGet
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

func FindPostTitles(titles *[]string, s *string, limit int) error {
	sqlStatement := `SELECT	DISTINCT posts.title
FROM public.posts
INNER JOIN public.abstract_users
ON posts.email_id = abstract_users.id
WHERE LOWER(posts.title) LIKE '%' || $1 || '%'
LIMIT $2`

	rows, err := utils.Db.Query(sqlStatement, s, strconv.Itoa(limit))
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

func UpdatePost(post *models.PostToUpdate, emailID *int) error {
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
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "title/content",
			Reason: "Could not map any of the provided fields"}
	}
	query.WriteString(fmt.Sprintf(" WHERE posts.id = $1 AND posts.email_id = $2"))

	res, err := utils.Db.Exec(query.String(), params...)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "id/email",
			Reason: "Could not change the specified post ID. Wrong ID/email combination"}
	}
	return nil
}

func DeletePost(emailID *int, postID *int) (int, error) {
	sqlStatement := `DELETE FROM public.posts WHERE posts.id = $1 AND 
                                        					 posts.email_id = $2;`
	res, err := utils.Db.Exec(sqlStatement, postID, emailID)
	if err != nil {
		return http.StatusBadRequest, err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return http.StatusUnauthorized, &utils.InvalidFieldsError{Location: "Body", AffectedField: "id",
			Reason: "Could not perform action on the specified post ID"}
	}
	return http.StatusAccepted, nil
}

func LikePost(emailID *int, postID *int) error {
	sqlStatement := `INSERT INTO public.likes (post_id, liked_by) VALUES ($1, $2);`
	res, err := utils.Db.Exec(sqlStatement, postID, emailID)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "id",
			Reason: "Could not perform action on the specified post ID"}
	}

	return nil
}
