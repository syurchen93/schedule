package manager

import (
	"schedule/model"
	"schedule/model/league"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syurchen93/api-football-client/common"
)

func Test_getCompetitionStandings(t *testing.T) {
	standingModels := []league.Standing{
		createStandingModel("Team A2", 2, "Group A"),
		createStandingModel("Team B2", 2, "Group B"),
		createStandingModel("Team C1", 1, "Group C"),
		createStandingModel("Team A1", 1, "Group A"),
		createStandingModel("Team B1", 1, "Group B"),
		createStandingModel("Team C2", 2, "Group C"),
	}

	standingDatas := buildStandingDatas(standingModels)

	assert.Equal(t, 3, len(standingDatas))
	assert.Equal(t, "Group A", standingDatas[0].GroupName)
	assert.Equal(t, 2, len(standingDatas[0].Standings))
	assert.Equal(t, "Team A1", standingDatas[0].Standings[0].TeamName)
	assert.Equal(t, 1, standingDatas[0].Standings[0].Position)
	assert.Equal(t, "Team A2", standingDatas[0].Standings[1].TeamName)
	assert.Equal(t, 2, standingDatas[0].Standings[1].Position)
}

func Test_generateScoreString(t *testing.T) {
	fixture := league.Fixture{
		Status: common.NotStarted,
	}
	assert.Equal(t, "", generateScoreString(fixture))

	fixture.Status = common.Finished
	fixture.GoalsHome = 2
	fixture.GoalsAway = 1
	assert.Equal(t, "2 : 1", generateScoreString(fixture))

	fixture.Status = common.FinishedAfterExtra
	assert.Equal(t, "2 : 1 (ET)", generateScoreString(fixture))

	fixture.Status = common.FinishedAfterPenalty
	fixture.PenaltyHome = 4
	fixture.PenaltyAway = 3
	assert.Equal(t, "2 : 1 P(4 : 3)", generateScoreString(fixture))
}

func createStandingModel(teamName string, position int, groupName string) league.Standing {
	dummyCode := "test"
	return league.Standing{
		Group: groupName,
		Rank:  position,
		Team: model.Team{
			Name: teamName,
			Code: &dummyCode,
		},
	}
}
