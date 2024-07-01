package cache

import (
	"context"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
)

var cfKeyId = "emails_per_hackathon:hackathon_id"

func SetupDefaultCache() {
	/**/ ctx := context.Background()
	status := db.DefaultRedisClient.CFReserve(ctx, cfKeyId, 100000)
	if status.Err() != nil {
		utils.MySugarLogger.Errorln(status.Err().Error())
	}

}

func RemoveEmailFromCache(email string) bool {
	ctx := context.Background()
	cmd := db.DefaultRedisClient.CFDel(ctx, cfKeyId, email)
	if cmd.Err() != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return false
	}
	return cmd.Val()
}

func AddEmailToCache(email string) bool {
	ctx := context.Background()
	cmd := db.DefaultRedisClient.CFAddNX(ctx, cfKeyId, email)
	if cmd.Err() != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return false
	}
	return cmd.Val()
}

func FindEmailInCache(email string) bool {
	ctx := context.Background()
	cmd := db.DefaultRedisClient.CFExists(ctx, cfKeyId, email)
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

func FindEmailsInCache(emails []interface{}) ([]bool, error) {
	ctx := context.Background()
	cmd := db.DefaultRedisClient.CFMExists(ctx, cfKeyId, emails...)
	if cmd.Err() != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return nil, cmd.Err()
	}
	found, err := cmd.Result()
	if err != nil {
		utils.MySugarLogger.Error("%v", cmd.Err().Error())
		return nil, cmd.Err()
	}
	return found, nil
}
