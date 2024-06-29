package taskmgt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
)

func FormatTaskCoordinatorKey(str string) string {
	return strings.Join([]string{"task_coordinators", str}, ":")
}

func GetSendParticipantWelcomeAndVerificationCoordinatorById(id string) (*exports.SendParticipantWelcomeAndVerificationEmailCoordinatorState, error) {
	cmd := db.DefaultRedisClient.HGetAll(context.Background(), FormatTaskKey(id))
	str, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	ttl, err := time.Parse(time.RFC3339, str["ttl"])
	if err != nil {
		fmt.Println(err.Error())
	}

	tsk := &exports.SendParticipantWelcomeAndVerificationEmailCoordinatorState{
		CoordinatorId:  str["id"],
		CurrentStateId: str["current_state_id"],
		Email:          str["email"],
		Token:          str["token"],
		LastName:       str["last_name"],
		FirstName:      str["first_name"],
		TeamName:       str["team_name"],
		TeamRole:       str["team_role"],
		TeamLeadName:   str["team_lead_name"],
		TTL:            ttl,
	}
	return tsk, nil
}

func SaveSendParticipantWelcomeAndVerificationEmailCoordinatorTaskById(tskInput *exports.SendParticipantWelcomeAndVerificationEmailCoordinatorState) error {
	mp := make(map[string]string)
	mp["id"] = tskInput.CoordinatorId
	mp["last_name"] = tskInput.LastName
	mp["first_name"] = tskInput.FirstName
	mp["email"] = tskInput.Email
	mp["token"] = tskInput.Token
	mp["ttl"] = tskInput.TTL.String()
	mp["team_lead_name"] = tskInput.TeamLeadName
	mp["team_role"] = tskInput.TeamRole
	mp["team_name"] = tskInput.TeamName
	pipe := db.DefaultRedisClient.TxPipeline()
	cmd := pipe.HSet(context.Background(), FormatTaskCoordinatorKey(tskInput.CoordinatorId), mp)
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("Redis CMD result id: %d\n", in)
	res, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}
	exports.MySugarLogger.Info(res)
	return nil
}

func UpdateSendParticipantWelcomeAndVerificationEmailCoordinatorById(id, status string) error {
	cmd := db.DefaultRedisClient.HSet(context.Background(), FormatTaskKey(id), "status", status)
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%d\n", in)
	return nil
}

func UpdateSendParticipantWelcomeAndVerificationEmailCoordinatorTaskStatusById(id, status string) error {
	cmd := db.DefaultRedisClient.HSet(context.Background(), FormatTaskKey(id), "status", status)
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%d\n", in)
	return nil
}

func DeleteSendParticipantWelcomeAndVerificationCoordinatorTaskStatusById(id string) error {
	cmd := db.DefaultRedisClient.HDel(context.Background(), FormatTaskKey(id))
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%d\n", in)
	return nil
}
