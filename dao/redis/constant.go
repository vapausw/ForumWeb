package redis

import (
	"errors"
	"github.com/go-redis/redis"
)

const (
	KeyPostInfoHashPrefix = "bluebell:post:"
	KeyPostTimeZSet       = "bluebell:post:time"
	KeyPostScoreZSet      = "bluebell:post:score"
	//KeyPostVotedUpSetPrefix   = "bluebell:post:voted:down:"
	//KeyPostVotedDownSetPrefix = "bluebell:post:voted:up:"
	KeyPostVotedZSetPrefix = "bluebell:post:voted:"

	KeyCommunityPostSetPrefix = "bluebell:community:"

	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 432
	PostPerAge               = 20
)

var (
	rdb                 *redis.Client
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoted          = errors.New("已经投过票了")
)
