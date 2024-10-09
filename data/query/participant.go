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
	"go.mongodb.org/mongo-driver/mongo"
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
	fmt.Printf("\n\n%+v\n\n", opts)
	participantCol, err := q.Datasource.GetParticipantCollection()
	if err != nil {
		return nil, err
	}
	ctx := context.Context(context.Background())
	var result []exports.ParticipantTeamMembersWithAccountsAggregateDocument
	var participant_id *string
	var participant_status *string
	var solution_like string
	var limit int = 1_000_000
	var review_ranking_top *int
	var review_ranking_eq *int
	var sort_review_ranking int = 1
	var participant_type *string
	if opts.ParticipantId != nil && *opts.ParticipantId != "" {
		participant_id = opts.ParticipantId
	}
	if opts.ParticipantStatus != nil && *opts.ParticipantStatus != "" {
		participant_status = opts.ParticipantStatus
	}
	if opts.Limit != nil && *opts.Limit > 0 {
		limit = *opts.Limit
	}
	if opts.Solution_Like != nil {
		solution_like = *opts.Solution_Like
	}
	if opts.ReviewRanking_Top != nil && *opts.ReviewRanking_Top > 0 {
		review_ranking_top = opts.ReviewRanking_Top
	}
	if opts.SortByReviewRanking_Desc != nil && *opts.SortByReviewRanking_Desc {
		sort_review_ranking = -1
	} else if opts.SortByReviewRanking_Asc != nil && *opts.SortByReviewRanking_Asc {
		sort_review_ranking = 1
	}
	if opts.ParticipantType != nil && *opts.ParticipantType != "" {
		participant_type = opts.ParticipantType
	}
	if opts.ReviewRanking_Eq != nil && *opts.ReviewRanking_Eq > 0 {
		review_ranking_eq = opts.ReviewRanking_Eq
	}
	fmt.Println(sort_review_ranking)
	var lookup_accounts_stage = bson.D{
		{"$lookup",
			bson.D{
				{"from", "accounts"},
				{"localField", "participant_id"},
				{"foreignField", "participant_id"},
				{"as", "accounts"},
			},
		},
	}

	var add_solution_as_object_id_field_stage = bson.D{
		{"$addFields", bson.D{
			{"solution_id_as_object_id", bson.D{
				{"$toObjectId", "$solution_id"},
			},
			},
		},
		},
	}

	var lookup_solutions_stage = bson.D{
		{"$lookup",
			bson.D{
				{"from", "solutions"},
				{"localField", "solution_id_as_object_id"},
				{"foreignField", "_id"},
				{"as", "solutions"},
			},
		},
	}

	var add_solution_document_stage = bson.D{
		{"$addFields", bson.D{
			{"solution_document", bson.D{
				{"$first", "$solutions"}},
			}},
		},
	}

	var add_team_lead_field = bson.D{
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
	}

	var add_co_participants_field = bson.D{
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
	}

	var add_team_lead_info_and_co_participants_fields_stage = bson.D{
		{"$addFields",
			bson.D{{"team_lead_info", add_team_lead_field},
				{"co_participants", add_co_participants_field},
			},
		},
	}

	var add_team_lead_info_id_sub_field = bson.D{{"$toString", "$team_lead_info._id"}}

	var add_each_co_participant_id_sub_field = bson.D{
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
	}

	var add_each_co_participant_i_sub_field = bson.D{
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
	}

	var filter_by_participant_id = bson.D{
		{"$eq",
			bson.A{
				"$participant_id",
				bson.D{
					{"$ifNull",
						bson.A{
							participant_id,
							"$participant_id",
						},
					},
				},
			},
		},
	}

	var filter_by_participant_type = bson.D{
		{"$eq",
			bson.A{
				"$type",
				bson.D{
					{"$ifNull",
						bson.A{
							participant_type,
							"$type",
						},
					},
				},
			},
		},
	}

	var filter_by_status = bson.D{
		{"$eq",
			bson.A{
				"$status",
				bson.D{
					{"$ifNull",
						bson.A{
							participant_status,
							"$status",
						},
					},
				},
			},
		},
	}

	var filter_by_solution = bson.D{
		{"$or",
			bson.A{
				bson.D{
					{"$regexMatch",
						bson.D{
							{"input", "$solution_document.title"},
							{"regex", solution_like},
							{"options", "i"},
						},
					},
				},
				bson.D{
					{"$regexMatch",
						bson.D{
							{"input", "$solution_document.obective"},
							{"regex", solution_like},
							{"options", "i"},
						},
					},
				},
				bson.D{
					{"$regexMatch",
						bson.D{
							{"input", "$solution_document.description"},
							{"regex", solution_like},
							{"options", "i"},
						},
					},
				},
			},
		},
	}

	var filter_by_ranking = bson.D{
		{"$eq",
			bson.A{
				bson.D{
					{"$ifNull",
						bson.A{
							review_ranking_eq,
							0,
						},
					},
				},
				bson.D{
					{"$convert",
						bson.D{
							{"input", "$review_ranking"},
							{"onError", 0},
							{"onNull", 0},
							{"to", "int"},
						},
					},
				},
			},
		},
	}
	var filters_match_stage primitive.D
	if solution_like != "" || participant_id != nil || participant_status != nil {
		var andFilters bson.A = bson.A{}
		if participant_id != nil {
			andFilters = append(andFilters, filter_by_participant_id)
		}
		if solution_like != "" {
			andFilters = append(andFilters, filter_by_solution)
		}
		if participant_status != nil {
			andFilters = append(andFilters, filter_by_status)
		}
		if participant_type != nil {
			fmt.Println(participant_type)
			andFilters = append(andFilters, filter_by_participant_type)
		}

		if review_ranking_eq != nil {
			andFilters = append(andFilters, filter_by_ranking)
		}

		filters_match_stage = bson.D{

			{"$match",
				bson.D{
					{"$expr",
						bson.D{
							{"$and", andFilters},
						},
					},
				},
			},
		}
	}

	var cleanup_stage = bson.D{
		{"$unset",
			bson.A{
				"solutions",
				"accounts",
				"solution_id_as_object_id",
			},
		},
	}

	/*bson.D{
			{
				"$sort", bson.A{
					"review_ranking", sort_review_ranking,
				},
			},
		}
	var sort_stage primitive.D
		if sort_stage != nil {
			sort_stage = bson.D{
				{
					"$sort", bson.A{},
				},
			}
		}*/
	var top_review_sort_stage primitive.D
	var top_review_limit_stage primitive.D
	if review_ranking_top != nil {
		top_review_sort_stage = bson.D{{"$sort", bson.D{{"review_ranking", -1}}}}
		top_review_limit_stage = bson.D{{"$limit", review_ranking_top}}
	}

	var limit_stage primitive.D
	if limit > 0 {
		limit_stage = bson.D{{"$limit", limit}}
	}
	/**/
	//panic("here")
	var pipeline mongo.Pipeline = make(mongo.Pipeline, 0)
	pipeline = append(pipeline, lookup_accounts_stage)
	pipeline = append(pipeline, add_solution_as_object_id_field_stage)
	pipeline = append(pipeline, lookup_solutions_stage)
	pipeline = append(pipeline, add_solution_document_stage)
	pipeline = append(pipeline, add_team_lead_info_and_co_participants_fields_stage)
	pipeline = append(pipeline,
		bson.D{
			{"$addFields",
				bson.D{
					{"team_lead_info.id", add_team_lead_info_id_sub_field},
					{"co_participants.id", add_each_co_participant_id_sub_field},
					{"co_participants.i", add_each_co_participant_i_sub_field},
				},
			},
		})
	if filters_match_stage != nil {

		pipeline = append(pipeline, filters_match_stage)
	}
	pipeline = append(pipeline,
		bson.D{
			{"$unset",
				bson.A{
					"solutions",
					"accounts",
					"solution_id_as_object_id",
				},
			},
		})
	pipeline = append(pipeline, cleanup_stage)
	//pipeline = append(pipeline, sort_stage)
	/**/
	if top_review_limit_stage != nil {
		pipeline = append(pipeline, top_review_sort_stage)
		pipeline = append(pipeline, top_review_limit_stage)
	}

	if limit_stage != nil {
		pipeline = append(pipeline, limit_stage)
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
	if err != nil {
		return nil, err
	}
	ctx := context.Context(context.Background())

	session, err := participantCol.Database().Client().StartSession()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, errors.New("failed to add to invite list: failed to match")
	}

	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {

		if err != nil {
			return nil, err
		}
		filter := bson.M{
			"team_lead_email": dataToSave.Email,
			"hackathon_id":    dataToSave.HackathonId,
			"$or": bson.A{
				bson.M{"co_participants.email": bson.M{"$in": bson.A{dataToSave.Email}}},
				bson.M{"invite_list.email": bson.M{"$in": bson.A{dataToSave.Email}}},
			},
		}

		partAnyRecord := &exports.ParticipantDocument{}
		resultAnyRecord := participantCol.FindOne(ctx, filter)

		err = resultAnyRecord.Decode(partAnyRecord)
		if err == nil {
			return nil, errors.New("email belongs to a participant or has been invited to another team")
		}
		if err != mongo.ErrNoDocuments {
			fmt.Println(err)
			return nil, errors.New("unable to add user to list.")
		}
		filter = bson.M{
			"participant_id": dataToSave.ParticipantId,
			"hackathon_id":   dataToSave.HackathonId,
			//"invite_list.email": bson.M{"$nin": bson.A{dataToSave.Email}},
		}

		part := &exports.ParticipantDocument{}
		result := participantCol.FindOne(ctx, filter)

		err = result.Decode(part)
		if err != nil {
			return nil, fmt.Errorf("unable to add user to list.")

		}
		if len(part.CoParticipants) > 0 {
			for _, v := range part.CoParticipants {
				if v.Email == dataToSave.Email {
					return nil, errors.New(fmt.Sprintf("email %s is already a co-participant.", dataToSave.Email))
				}
			}
			if len(part.CoParticipants) == 3 {
				return nil, errors.New("Maximum number of co-participants reached.")
			}

			if len(part.CoParticipants)+len(part.InviteList) >= 3 {
				return nil, errors.New("Cannot invite more co-participants reached. ")
			}
		}
		if len(part.InviteList) > 0 {
			for _, v := range part.InviteList {
				if v.Email == dataToSave.Email {
					return nil, fmt.Errorf("email %s is already invited", dataToSave.Email)
				}
			}
		}
		upd := bson.M{
			"$addToSet": bson.M{"invite_list": exports.ParticipantDocumentTeamInviteInfo{
				Email:     dataToSave.Email,
				InviterId: dataToSave.InviterEmail, Time: time.Now()}},
			"$set": bson.M{"updated_at": time.Now()},
		}
		updateResult, err := participantCol.UpdateOne(ctx, filter, upd)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return nil, errors.New("failed to add to invite list: failed to match")
		}
		fmt.Printf("%#v", result)
		if updateResult.MatchedCount == 0 {
			fmt.Printf("failed to add to invite list")
			return nil, errors.New("failed to add to invite list: failed to match")
		}
		if updateResult.ModifiedCount == 0 && updateResult.UpsertedCount == 0 {
			fmt.Printf("No changes made")
			return nil, errors.New("failed to add to invite list: failed to save")
		}
		fmt.Println("Added to invite list")

		return updateResult, nil
	})

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
	fmt.Println(fmt.Sprintf("Added to co participants' list %s\n", dataToSave.Email))
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
	fmt.Println(dataToSave.SolutionId, dataToSave.ParticipantId, "uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu")
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

	if solDoc.Id.Hex() != dataToSave.SolutionId {
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

func (q *Query) UpdateSingleParticipantRecord(filterOpts *exports.UpdateSingleParticipantRecordFilter, dataToSave *exports.UpdateParticipantRecordData) (*exports.ParticipantDocument, error) {
	participantCol, err := q.Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"$expr": bson.M{
			"$eq": bson.A{
				"$participant_id", bson.M{"$ifNull": bson.A{filterOpts.ParticipantId, "$participant_id"}},
			},
		},
	}
	up := bson.M{}
	if dataToSave.Status != "" {
		up["status"] = dataToSave.Status
	}
	if dataToSave.ReviewRanking > 0 {
		up["review_ranking"] = dataToSave.ReviewRanking
	}
	len_of_map := 0
	for range up {
		len_of_map += 0
	}
	if len_of_map > 0 {
		up["updated_at"] = time.Now()
	}
	upd := bson.M{
		"$set": up,
	}
	fmt.Println(upd)
	docPart := &exports.ParticipantDocument{}
	after := options.After
	result := participantCol.FindOneAndUpdate(ctx, filter, upd, &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	})
	if err := result.Decode(docPart); err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Println("Updated successfully")
	return docPart, err
}
