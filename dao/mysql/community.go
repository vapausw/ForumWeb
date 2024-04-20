package mysql

import (
	"ForumWeb/model"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*model.Community, err error) {
	sqlStr := "select community_id, community_name from community"

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, ErrServiceBusy
	}
	for rows.Next() {
		var community model.Community
		if err = rows.Scan(&community.CommunityID, &community.CommunityName); err != nil {
			return nil, ErrServiceBusy
		}
		communityList = append(communityList, &community)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrServiceBusy
	}
	if len(communityList) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}

func GetCommunityByID(communityID string) (community *model.CommunityDetail, err error) {
	sqlStr := "select community_id, community_name, introduction, create_time from community where community_id = ?"

	community = new(model.CommunityDetail)
	err = db.QueryRow(sqlStr, communityID).Scan(&community.CommunityID, &community.CommunityName, &community.Introduction, &community.CreateTime)
	if err != nil {
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		return nil, ErrServiceBusy
	}
	return
}
