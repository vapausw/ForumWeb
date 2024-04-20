package mysql

import (
	"ForumWeb/model"
	"database/sql"
	"go.uber.org/zap"
)

func CreateComment(comment *model.Comment) (err error) {
	sqlStr := `insert into comment(comment_id, author_id, post_id, content, parent_id) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.AuthorID, comment.PostID, comment.Content, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		return ErrInsertFailed
	}
	return
}

func GetCommentListByIDs(ids []string) (comments []*model.Comment, err error) {
	sqlStr := `select comment_id, author_id, post_id, content from comment where comment_id in (?)`
	rows, err := db.Query(sqlStr, ids)
	if err != nil {
		zap.L().Error("query comment list failed", zap.Error(err))
		return nil, ErrServiceBusy
	}
	for rows.Next() {
		var comment model.Comment
		if err = rows.Scan(&comment.CommentID, &comment.AuthorID, &comment.PostID, &comment.Content); err != nil {
			zap.L().Error("scan comment failed", zap.Error(err))
			continue
		}
		comments = append(comments, &comment)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrServiceBusy
	}
	if len(comments) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}
