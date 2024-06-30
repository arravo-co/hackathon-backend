package jobs

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/utils"
)

func ConsumeUploadPicQueue(payloadStruct *exports.UploadPicQueuePayload) error {
	fmt.Println("Starting to consume")
	uploadResult, err := utils.UploadPicToCloudinary(payloadStruct.FilePath)
	if err != nil {
		return err
	}
	fmt.Println(uploadResult)

	judge, err := repository.GetJudgeEntity(payloadStruct.Email)
	if err != nil {
		return err
	}
	fmt.Println("Here at this stage.")
	err = judge.UpdateJudgeProfile(dtos.UpdateJudgeDTO{
		ProfilePictureUrl: uploadResult.SecureURL,
	})
	if err != nil {
		return err
	}
	return nil
}
