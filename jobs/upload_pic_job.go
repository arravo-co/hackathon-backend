package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
	"github.com/arravoco/hackathon_backend/utils"
)

type UploadPicConsumer struct {
	Ch chan interface{}
}

// queue name = upload_pic
func StartConsumingJob(q_name string) error {
	queue, err := rmqUtils.GetQueue(q_name)
	if err != nil {
		fmt.Println("Error getting queue")
		return err
	}

	err = queue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func CreateUploadPicToCloudinaryConsumer() (*UploadPicConsumer, error) {
	queue, err := rmqUtils.GetQueue("upload_pic_cloudinary")
	if err != nil {
		fmt.Println("Error getting queue")
		return nil, err
	}

	err = queue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	uploadPicConsumer := &UploadPicConsumer{}
	str, err := queue.AddConsumer("upload_pic_cloudinary", uploadPicConsumer)
	if err != nil {
		return nil, err
	}
	fmt.Println(str)
	return uploadPicConsumer, nil
}

func (c *UploadPicConsumer) Consume(d rmq.Delivery) {
	payload := d.Payload()
	fmt.Println("Starting to consume")
	payloadStruct := exports.UploadPicQueuePayload{}
	err := json.Unmarshal([]byte(payload), &payloadStruct)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	uploadResult, err := utils.UploadPicToCloudinary(payloadStruct.FilePath)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	fmt.Println(uploadResult)

	judge, err := entity.GetJudgeEntity(payloadStruct.Email)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	fmt.Println("Here at this stage.")
	err = judge.UpdateJudgeProfile(dtos.UpdateJudgeDTO{
		ProfilePictureUrl: uploadResult.SecureURL,
	})
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
