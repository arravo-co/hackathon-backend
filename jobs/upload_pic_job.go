package jobs

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/services"
	"github.com/arravoco/hackathon_backend/utils"
)

func ConsumeUploadJudgeProfilePicQueue(payloadStruct *exports.UploadJudgeProfilePicQueuePayload) error {
	fmt.Println("Starting to consume")
	uploadResult, err := utils.UploadPicToCloudinary(payloadStruct.FilePath)
	if err != nil {
		return err
	}
	fmt.Println(uploadResult)

	serv := services.GetServiceWithDefaultRepositories()
	judge, err := serv.GetJudgeByEmail(payloadStruct.Email)
	if err != nil {
		return err
	}
	fmt.Println("Here at this stage.")
	_, err = serv.UpdateJudgeInfo(judge.Email, &services.UpdateJudgeDTO{
		ProfilePictureUrl: uploadResult.SecureURL,
	})
	if err != nil {
		return err
	}
	return nil
}

func ConsumeUploadSolutionPicQueue(payloadStruct *exports.UploadSolutionPicQueuePayload) error {
	uploadResult, err := utils.UploadPicToCloudinary(payloadStruct.FilePath)
	if err != nil {
		return err
	}
	fmt.Println(uploadResult)

	solServ := services.NewSolutionService()
	_, err = solServ.UpdateSolutionDataById(payloadStruct.SolutionId, &exports.UpdateSolutionData{
		SolutionImageUrl: uploadResult.SecureURL,
	})
	if err != nil {
		return err
	}
	return nil
}
