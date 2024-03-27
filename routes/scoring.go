package routes

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/labstack/echo/v4"
)

func ScoreParticipant(c echo.Context) error {
	data.UpdateParticipantScoreRecordByJudge(&exports.UpdateParticipantScoreFilterData{}, &exports.UpdateParticipantScoreDataByJudge{
		ScoreInfo: exports.ScoreInfo{
			JudgeId: "",
			Criteria: exports.ScoreCriterium{
				Score:     0,
				Criterium: "",
			},
			Weight: 1,
		},
	})
	return nil
}
