package controllers

import (
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"github.com/lib/pq"
	"strconv"
	"time"
)

func InsertComment(comm *models.CommentToCreate, emailID *int) error {
	sqlStatement := `INSERT INTO public.comments (post_id, content, email_id, date, is_edited) 
VALUES ($1, $2, $3, $4, $5);`

	_, err = utils.Db.Exec(sqlStatement, comm.PostID, comm.Content, *emailID, pq.FormatTimestamp(time.Now()),
		false)
	if err != nil {
		return err
	}

	return nil
}

func InsertRespondToComment(comm *models.CommentToCreate, emailID *int, commentID *int) error {
	sqlStatement := `INSERT INTO public.comments (post_id, content, email_id, date, is_edited, responding_to_id) 
VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = utils.Db.Exec(sqlStatement, comm.PostID, comm.Content, *emailID, pq.FormatTimestamp(time.Now()),
		false, *commentID)
	if err != nil {
		return err
	}

	return nil
}

func LikeComment(emailID *int, commentID *int) error {
	sqlStatement := `INSERT INTO public.comment_likes (comment_id, liked_by) VALUES ($1, $2);`
	res, err := utils.Db.Exec(sqlStatement, *commentID, *emailID)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		return &utils.InvalidFieldsError{Location: "Body", AffectedField: "id",
			Reason: "Could not perform action on the specified post ID"}
	}

	return nil
}

func UpdateComment(comm *models.CommentToCreate, emailID *int, commentID *int) error {
	sqlStatement := `UPDATE public.comments SET content = $2, date = $4, is_edited = true 
                       WHERE id = $5 AND post_id = $1 AND email_id = $3`
	_, err = utils.Db.Exec(sqlStatement, comm.PostID, comm.Content, *emailID, pq.FormatTimestamp(time.Now()),
		*commentID)
	if err != nil {
		return err
	}

	return nil
}

func GetCommentsFromPost(postID *int, comms *[]models.CommentToGet, commID *int) error {
	sqlStatement :=
		`SELECT concat_ws (' ', abstract_users.last_name, abstract_users.first_name) as full_name, 
 		post_id, content, date, is_edited, public.comments.id, 
 		(SELECT COUNT(a.liked_by) FROM comment_likes as a WHERE a.comment_id = public.comments.id) as number_of_likes,
 		responding_to_id
 	 FROM public.comments
 	 INNER JOIN public.abstract_users
	 ON comments.email_id = abstract_users.id
	 WHERE post_id = $1 AND responding_to_id`

	if *commID == 0 {
		sqlStatement += " IS NULL"
	} else {
		sqlStatement += " = " + strconv.Itoa(*commID)
	}

	rows, err := utils.Db.Query(sqlStatement, *postID)
	defer rows.Close()

	if err != nil {
		return err
	}
	var c models.CommentToGet
	for rows.Next() {
		if err = rows.Scan(&c.FullName, &c.PostID, &c.Content, &c.Date, &c.IsEdited, &c.CommentID,
			&c.NumberOfLikes, &c.RespondingToID); err != nil {
			return err
		}
		*comms = append(*comms, c)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
