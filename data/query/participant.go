package query

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		GithubAddress:    dataToSave.GithubAddress,
		ParticipantEmail: dataToSave.ParticipantEmail,
		InviteList:       []exports.ParticipantDocumentTeamInviteInfo{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           "UNREVIEWED",
	}
	result, err := participantCol.InsertOne(ctx, dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	dat.Id = result.InsertedID.(primitive.ObjectID)
	return &dat, nil
}

func (q *Query) GetParticipantsRecords() ([]exports.ParticipantDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := []exports.ParticipantDocument{}
	result, err := participantCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		fmt.Printf("\n%s\n", result.Err())
		return nil, result.Err()
	}
	err = result.All(context.Background(), &dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return dat, nil
}

func (q *Query) GetParticipantsWithAccountsAggregate(opts exports.GetParticipantsWithAccountsAggregateFilterOpts) ([]exports.ParticipantTeamMembersWithAccountsAggregateDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	if err != nil {
		return nil, err
	}
	ctx := context.Context(context.Background())
	var result []exports.ParticipantTeamMembersWithAccountsAggregateDocument
	pipeline := bson.A{
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
		bson.D{{"$addFields", bson.D{{"solution_id_as_object_id", bson.D{{"$toObjectId", "$solution_id"}}}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "solutions"},
					{"localField", "solution_id_as_object_id"},
					{"foreignField", "_id"},
					{"as", "solutions"},
				},
			},
		},
		bson.D{{"$addFields", bson.D{{"solution_document", bson.D{{"$first", "$solutions"}}}}}},
		bson.D{
			{"$addFields",
				bson.D{
					{"team_lead_info",
						bson.D{
							{"$first",
								bson.A{
									bson.D{
										{"$filter",
											bson.D{
												{"input", "$accounts"},
												{"as", "acc"},
												{"cond",
													bson.D{
														{"$eq",
															bson.A{
																"$$acc.email",
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
					{"co_participants",
						bson.D{
							{"$let",
								bson.D{
									{"vars",
										bson.D{
											{"co_parts",
												bson.D{
													{"$filter",
														bson.D{
															{"input", "$accounts"},
															{"as", "acc"},
															{"cond",
																bson.D{
																	{"$ne",
																		bson.A{
																			"$$acc.email",
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
									{"in",
										bson.D{
											{"$map",
												bson.D{
													{"input", "$$co_parts"},
													{"as", "item"},
													{"in",
														bson.D{
															{"$mergeObjects",
																bson.A{
																	bson.D{
																		{"$mergeObjects",
																			bson.A{
																				"$$item",
																				bson.D{{"account_role", "$$item.role"}},
																			},
																		},
																	},
																	bson.D{
																		{"$let",
																			bson.D{
																				{"vars",
																					bson.D{
																						{"obj",
																							bson.D{
																								{"$first",
																									bson.D{
																										{"$filter",
																											bson.D{
																												{"input", "$co_participants"},
																												{"as", "acc"},
																												{"cond",
																													bson.D{
																														{"$eq",
																															bson.A{
																																"$$item.email",
																																"$$acc.email",
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
																				{"in",
																					bson.D{
																						{"$mergeObjects",
																							bson.A{
																								"$$obj",
																								bson.D{{"team_role", "$$obj.role"}},
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
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$addFields", bson.D{{"account_role.account_role", "$team_lead_info.role"}}}},
		bson.D{
			{"$addFields",
				bson.D{
					{"team_lead_info.id", bson.D{{"$toString", "$team_lead_info._id"}}},
					{"co_participants.id",
						bson.D{
							{"$first",
								bson.D{
									{"$map",
										bson.D{
											{"input", "$co_participants"},
											{"as", "co_part"},
											{"in", bson.D{{"$toString", "$$co_part._id"}}},
										},
									},
								},
							},
						},
					},
					{"co_participants.i",
						bson.D{
							{"$first",
								bson.D{
									{"$map",
										bson.D{
											{"input", "$co_participants"},
											{"as", "co_part"},
											{"in", bson.D{{"$toString", "$$co_part._id"}}},
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
			{"$match",
				bson.D{
					{"$expr",
						bson.D{
							{"$and",
								bson.A{
									bson.D{
										{"$eq",
											bson.A{
												"$type",
												bson.D{
													{"$ifNull",
														bson.A{
															"TEAM",
															"$type",
														},
													},
												},
											},
										},
									},
									bson.D{
										{"$eq",
											bson.A{
												"$status",
												bson.D{
													{"$ifNull",
														bson.A{
															opts.ParticipantStatus,
															"$status",
														},
													},
												},
											},
										},
									},
									bson.D{
										{"$regexMatch",
											bson.D{
												{"input", "status"},
												{"regex", ""},
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
					"solutions",
					"accounts",
					"solution_id_as_object_id",
				},
			},
		},
	}

	cursor, err := participantCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
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
	if err != nil {
		return nil, err
	}
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
		"$addToSet": bson.M{"invite_list": exports.ParticipantDocumentTeamInviteInfo{Email: dataToSave.Email,
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

	//fmt.Println("\n\n\n", filter, "\n\n\n")

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

	//fmt.Println("\n\n", dataToSave.MemberEmail, "\n\n")

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

func (q *Query) SelectSolutionForTeam(dataToSave *exports.SelectTeamSolutionData) (*exports.ParticipantDocument, error) {
	fmt.Println(dataToSave.SolutionId, "uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu")
	solId, err := primitive.ObjectIDFromHex(dataToSave.SolutionId)
	if err != nil {
		return nil, err
	}
	solCol, err := q.Datasource.GetSolutionCollection()
	if err != nil {
		return nil, err
	}
	fmt.Println(solId, "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy")
	solDoc := exports.SolutionDocument{}
	solResult := solCol.FindOne(context.Background(), bson.M{"_id": solId})
	if solResult.Err() != nil {
		return nil, solResult.Err()
	}
	if err := solResult.Decode(&solDoc); err != nil {
		return nil, err
	}

	if solDoc.Title == "" {
		return nil, errors.New("No document with id " + dataToSave.SolutionId)
	}
	participantCol, err := q.Datasource.GetParticipantCollection()
	if err != nil {
		return nil, err
	}
	ctx := context.Context(context.Background())
	filter := bson.M{
		"participant_id": dataToSave.ParticipantId,
		"hackathon_id":   dataToSave.HackathonId,
	}

	//fmt.Println("\n\n", dataToSave.MemberEmail, "\n\n")

	upd := bson.M{
		"$set": bson.M{"solution_id": dataToSave.SolutionId},
	}
	partDoc := exports.ParticipantDocument{}
	ret := options.After
	resultPartDoc := participantCol.FindOneAndUpdate(ctx, filter, upd, &options.FindOneAndUpdateOptions{
		ReturnDocument: &ret,
	})
	if err := resultPartDoc.Decode(&partDoc); err != nil {
		return nil, err
	}
	partDoc.Solution = exports.ParticipantDocumentParticipantSelectedSolution{}
	return &partDoc, err
}
