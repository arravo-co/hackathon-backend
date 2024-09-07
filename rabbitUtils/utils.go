package rabbitutils

import (
	"github.com/aidarkhanov/nanoid"
)

/*
	func ListenToAllQueues(consumer *consum) {
		chUploadPicQicJobDelivery, err := ConsumeQueue("upload.profile_picture.cloudinary", GetConsumerTag())
		if err != nil {
			fmt.Println(err.Error())
		}

		chWelcomeEmailToJudgeDelivery, err := ConsumeQueue("send.judge.registered.admin.welcome_email", GetConsumerTag())
		if err != nil {
			fmt.Println(err.Error())
		}

		//send.participant.created.welcome_email_verification_email
		chCoordinateParticipantWelcomeVerDelivery, err := ConsumeQueue("send.participant.created.welcome_email_verification_email", GetConsumerTag())
		if err != nil {
			fmt.Println(err.Error())
		}

		chUploadSolutionPicQueue, err := ConsumeQueue("upload.solution_picture.cloudinary", GetConsumerTag())
		if err != nil {
			fmt.Println(err.Error())
		}

		for {
			select {
			case response := <-chUploadPicQicJobDelivery:
				fmt.Printf(response.ConsumerTag)
				HandleJudgeProfilePicUploadConsumption(&response)
			case response := <-chWelcomeEmailToJudgeDelivery:
				HandleSendEmailToJudgeConsumption(&response)
			case response := <-chCoordinateParticipantWelcomeVerDelivery:
				fmt.Printf(response.ConsumerTag)
			case response := <-chUploadSolutionPicQueue:
				HandleSolutionPicUploadConsumption(&response)
				//HandleCoordinateParticipantWelcomeVerificationConsumption(&response)
			}
		}
	}
*/
func GetConsumerTag() string {

	id := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789", 10))
	return id
}
