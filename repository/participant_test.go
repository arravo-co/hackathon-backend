package repository

import (
	"testing"

	"github.com/arravoco/hackathon_backend/dtos"
)

func TestRegisterTeamLead(t *testing.T) {
	p := ParticipantRepository{}
	args := dtos.RegisterNewParticipantDTO{
		FirstName: "Temitope",
		LastName:  "Alabi",
		Password:  "david",
		Skillset:  []string{"nodejs", "sql"},
		State:     "OSUN",
	}

	t.Run("RegisterTeamLead", func(t *testing.T) {

		res, err := p.RegisterTeamLead(args)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		if res.HackathonId != "" {

		}
	})
}
