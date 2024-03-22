package cache

import (
	"context"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
)

var cfKeyId = "emails_per_hackathon:hackathon_id"

func init() {
	ctx := context.Background()
	status := db.RedisClient.CFReserve(ctx, cfKeyId, 100000)
	if status.Err() != nil {
		utils.MySugarLogger.Error("%v", status.Err().Error())
	}
}

func AddEmailToCache(email string) {
	ctx := context.Background()
	cmd := db.RedisClient.CFAdd(ctx, cfKeyId, email)
	if cmd.Err() != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
	}
}

func FindEmailInCache(email string) bool {
	ctx := context.Background()
	cmd := db.RedisClient.CFExists(ctx, cfKeyId, email)
	if cmd.Err() != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return false
	}
	found, err := cmd.Result()
	if err != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return false
	}
	return found
}
