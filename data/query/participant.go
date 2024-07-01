package query

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (q *Query) CreateParticipantRecord(dataToSave *exports.CreateParticipantRecordData) (*exports.ParticipantDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{
		ParticipantId:    dataToSave.ParticipantId,
		HackathonId:      dataToSave.HackathonId,
		Type:             dataToSave.Type,
		TeamLeadEmail:    dataToSave.TeamLeadEmail,
		TeamName:         dataToSave.TeamName,
		CoParticipants:   dataToSave.CoParticipants,
		GithubAddress:    dataToSave.GithubAddress,
		ParticipantEmail: dataToSave.ParticipantEmail,
		InviteList:       []exports.InviteInfo{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           "UNREVIEWED",
	}
	result, err := participantCol.InsertOne(ctx, dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	dat.Id = result.InsertedID
	return &dat, nil
}

func (q *Query) GetParticipantsRecords() ([]exports.ParticipantDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := &[]exports.ParticipantDocument{}
	result, err := participantCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		fmt.Printf("\n%s\n", result.Err())
		return nil, result.Err()
	}
	err = result.All(context.Background(), dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return *dat, nil
}

func (q *Query) GetParticipantsRecordsAggregate() ([]exports.ParticipantAccountWithCoParticipantsDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	if err != nil {
		return nil, err
	}
	ctx := context.Context(context.Background())
	var result []exports.ParticipantAccountWithCoParticipantsDocument
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"status", "UNREVIEWED"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "accounts"},
					{"localField", "participant_id"},
					{"foreignField", "participant_id"},
					{"as", "accounts"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"team_lead_info",
						bson.D{
							{"$ifNull",
								bson.A{
									bson.D{
										{"$arrayElemAt",
											bson.A{
												bson.D{
													{"$filter",
														bson.D{
															{"input", "$accounts"},
															{"as", "arr"},
															{"cond",
																bson.D{
																	{"$eq",
																		bson.A{
																			bson.D{{"$getField", "team_lead_email"}},
																			"$$arr.email",
																		},
																	},
																},
															},
														},
													},
												},
												0,
											},
										},
									},
									"",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$replaceRoot",
				bson.D{
					{"newRoot",
						bson.D{
							{"$mergeObjects",
								bson.A{
									"$team_lead_info",
									"$$ROOT",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"co_participants",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$co_participants"},
									{"as", "arr"},
									{"in",
										bson.D{
											{"$mergeObjects",
												bson.A{
													"$$arr",
													bson.D{{"team_role", "$$arr.role"}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$unset", "co_participants.role"}},
		bson.D{
			{"$addFields",
				bson.D{
					{"co_participants",
						bson.D{
							{"$map",
								bson.D{
									{"input",
										bson.D{
											{"$filter",
												bson.D{
													{"input", "$accounts"},
													{"as", "arr"},
													{"cond",
														bson.D{
															{"$not",
																bson.A{
																	bson.D{
																		{"$eq",
																			bson.A{
																				"$$arr.email",
																				"$team_lead_email",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									{"as", "item"},
									{"in",
										bson.D{
											{"$mergeObjects",
												bson.A{
													bson.D{
														{"$setField",
															bson.D{
																{"input", "$$item"},
																{"field", bson.D{{"$literal", "account_role"}}},
																{"value", "$$item.role"},
															},
														},
													},
													bson.D{
														{"$arrayElemAt",
															bson.A{
																bson.D{
																	{"$filter",
																		bson.D{
																			{"input", "$co_participants"},
																			{"as", "arr"},
																			{"cond",
																				bson.D{
																					{"$eq",
																						bson.A{
																							"$$arr.email",
																							"$$item.email",
																						},
																					},
																				},
																			},
																		},
																	},
																},
																0,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$unset",
				bson.A{
					"password_hash",
					"co_participants.password_hash",
					"team_lead_info",
				},
			},
		},
	}

	cursor, err := participantCol.Aggregate(ctx, pipeline)
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *Query) GetParticipantRecord(participantId string) (*exports.ParticipantDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{}
	result := participantCol.FindOne(ctx, bson.M{"participant_id": participantId})
	if result.Err() != nil {
		fmt.Printf("\n%s\n", result.Err())
		return nil, result.Err()
	}
	err = result.Decode(&dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return &dat, nil
}

func (q *Query) AddToTeamInviteList(dataToSave *exports.AddToTeamInviteListData) (interface{}, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": bson.M{"$nin": bson.A{dataToSave.Email}},
	}
	upd := bson.M{
		"$addToSet": bson.M{"invite_list": exports.InviteInfo{Email: dataToSave.Email,
			InviterId: dataToSave.InviterEmail, Time: time.Now()}},
		"$set": bson.M{"updated_at": time.Now()},
	}
	fmt.Println(upd)

	result, err := participantCol.UpdateOne(ctx, filter, upd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result)
	if result.MatchedCount == 0 {
		fmt.Printf("failed to add to invite list")
		return nil, errors.New("failed to add to invite list: failed to match")
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		fmt.Printf("No changes made")
		return nil, errors.New("failed to add to invite list: failed to save")
	}
	return result, err
}

func (q *Query) AddSolutionToTeam(dataToSave *exports.AddSolutionToTeamData) (interface{}, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id": dataToSave.ParticipantId,
		"hackathon_id":   dataToSave.HackathonId,
	}
	upd := bson.M{
		"$set": bson.M{"solution_id": dataToSave.SolutionId},
	}
	fmt.Println(upd)

	result, err := participantCol.UpdateOne(ctx, filter, upd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result)
	if result.MatchedCount == 0 {
		fmt.Printf("failed to add to invite list")
		return nil, errors.New("failed to add to invite list: failed to match")
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		fmt.Printf("No changes made")
		return nil, errors.New("failed to add to invite list: failed to save")
	}
	return result, err
}

func (q *Query) AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error) {
	partDoc := &exports.ParticipantDocument{}
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": dataToSave.Email,
	}

	fmt.Println("\n\n\n", filter, "\n\n\n")

	upd := bson.M{
		"$addToSet": bson.M{"co_participants": bson.M{"email": dataToSave.Email, "role": dataToSave.TeamRole}},
		"$pull":     bson.M{"invite_list": bson.M{"email": dataToSave.Email}},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	retDoc := options.After
	result := participantCol.FindOneAndUpdate(ctx, filter, upd, &options.FindOneAndUpdateOptions{ReturnDocument: &retDoc})
	err = result.Decode(partDoc)
	if err != nil {
		fmt.Printf("%s\n\n\n\n", err.Error())
		return nil, err
	}
	return partDoc, err
}

func (q *Query) RemoveMemberFromParticipatingTeam(dataToSave *exports.RemoveMemberFromParticipatingTeamData) (interface{}, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id": dataToSave.ParticipantId,
		"hackathon_id":   dataToSave.HackathonId,
	}

	fmt.Println("\n\n", dataToSave.MemberEmail, "\n\n")

	upd := bson.M{
		"$pull": bson.M{"co_participants.email": dataToSave.MemberEmail},
	}
	result, err := participantCol.UpdateOne(ctx, filter, upd, &options.UpdateOptions{})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		if strings.Contains(err.Error(), "mongo") {
			return nil, fmt.Errorf("Unexpected error")
		}
		return nil, err
	}
	fmt.Printf("%#v", result.ModifiedCount)
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no records found")
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no changes made.")
	}
	return participantCol, err
}
