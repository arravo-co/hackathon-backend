package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/arravoco/hackathon_backend/exports"
)

func CreateParticipantRecord(dataToSave *exports.CreateParticipantRecordData) (*exports.ParticipantDocument, error) {
	participantCol, err := Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	participantEmails := append([]string{dataToSave.ParticipantEmail, dataToSave.TeamLeadEmail}, dataToSave.CoParticipantEmails...)
	participantEmails = slices.Compact[[]string](participantEmails)
	participantId, err := GenerateParticipantID(participantEmails)
	if err != nil {
		exports.MySugarLogger.Fatalln(err)
		return nil, errors.New("failed to generate participant id")
	}
	dat := exports.ParticipantDocument{
		ParticipantId:       participantId,
		HackathonId:         dataToSave.HackathonId,
		Type:                dataToSave.Type,
		TeamLeadEmail:       dataToSave.TeamLeadEmail,
		TeamName:            dataToSave.TeamName,
		CoParticipantEmails: dataToSave.CoParticipantEmails,
		GithubAddress:       dataToSave.GithubAddress,
		ParticipantEmail:    dataToSave.ParticipantEmail,
	}
	result, err := participantCol.InsertOne(ctx, dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	dat.Id = result.InsertedID
	return &dat, nil
}

func GenerateParticipantID(emails []string) (string, error) {
	slices.Sort[[]string](emails)
	joined := strings.Join(emails, ":")
	h := sha256.New()
	_, err := h.Write([]byte(joined))
	if err != nil {
		return "", err
	}
	hashByte := h.Sum(nil)
	hashedString := fmt.Sprintf("%x", hashByte)
	fmt.Println("hashedString")
	fmt.Println(joined)
	fmt.Println(hashedString)
	fmt.Println("hashedString")
	slicesOfHash := strings.Split(hashedString, "")
	prefixSlices := slicesOfHash[0:5]
	postFix := slicesOfHash[len(slicesOfHash)-5:]
	sub := strings.Join(append(prefixSlices, postFix...), "")
	return sub, nil
}
