package exports

type CreateParticipantScoreData struct {
	ScoreId       string    `validate:"required" json:"score_id"`
	HackathonId   string    `validate:"required" json:"hackathon_id"`
	ParticipantId string    `validate:"required" json:"participant_id"`
	Stage         string    `validate:"required" json:"stage"`
	JudgeId       string    `validate:"required" json:"judge_id"`
	ScoreInfo     ScoreInfo `validate:"required" json:"score_info"`
}

type UpdateParticipantScoreFilterData struct {
	ScoreId       string `validate:"required" json:"score_id"`
	HackathonId   string `validate:"required" json:"hackathon_id"`
	ParticipantId string `validate:"required" json:"participant_id"`
	Stage         string `validate:"required" json:"stage"`
}

type UpdateParticipantScoreDataByJudge struct {
	ScoreInfo ScoreInfo `validate:"required" json:"score_info"`
}

type ScoreInfo struct {
	JudgeId  string         `bson:"judge" validate:"required" json:"judge"`
	Criteria ScoreCriterium `bson:"criteria" validate:"required" json:"criteria"`
	Weight   float32        `bson:"weight" validate:"required" json:"weight"`
}

type ScoreCriterium struct {
	Score     float32 `bson:"score" validate:"required" json:"score"`
	Criterium string  `bson:"criterium" validate:"required" json:"criterium"`
}
