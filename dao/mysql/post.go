package mysql

import (
	"ForumWeb/model"
	"go.uber.org/zap"
)

func CreatePost(p *model.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content, p.AuthorId, p.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrInsertFailed
		return
	}
	return
}

func GetPostByID(pid string) (post *model.ApiPostDetail, err error) {
	post = new(model.ApiPostDetail)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	err = db.QueryRow(sqlStr, pid).Scan(&post.PostID, &post.Title, &post.Content, &post.AuthorId, &post.CommunityID, &post.CreateTime)
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrServiceBusy
		return
	}
	return
}
