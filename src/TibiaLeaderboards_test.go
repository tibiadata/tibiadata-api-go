package main

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

func readLeaderboardsTestFile(t *testing.T, name string) string {
	t.Helper()

	file, err := static.TestFiles.Open("testdata/leaderboards/" + name)
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	return string(data)
}

func TestLeaderboardsPastRotation(t *testing.T) {
	const url = "https://www.tibia.com/community/?subtopic=leaderboards&world=Antica&rotation=138&currentpage=1"

	leaderboardsJson, err := TibiaLeaderboardsImpl("Antica", "138", 1, readLeaderboardsTestFile(t, "Antica-138-1.html"), url)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	leaderboards := leaderboardsJson.Leaderboards

	assert.Equal(url, leaderboardsJson.Information.TibiaURLs[0])
	assert.Equal("Antica", leaderboards.World)
	assert.Equal(LeaderboardRotation{RotationID: 138, Current: false, EndDate: "2026-10-28T09:00:00Z"}, leaderboards.Rotation)

	assert.Equal(2, len(leaderboards.AvailableRotations))
	assert.Equal(LeaderboardRotation{RotationID: 139, Current: true, EndDate: "2026-11-11T09:00:00Z"}, leaderboards.AvailableRotations[0])
	assert.Equal(LeaderboardRotation{RotationID: 138, Current: false, EndDate: "2026-10-28T09:00:00Z"}, leaderboards.AvailableRotations[1])

	assert.Empty(leaderboards.LastUpdate)

	assert.Equal(50, len(leaderboards.LeaderboardList))
	assert.Equal(LeaderboardEntry{Rank: 1, Name: "Ziobi", DromeLevel: 114}, leaderboards.LeaderboardList[0])
	assert.Equal(LeaderboardEntry{Rank: 2, Name: "Sabat czaruje", DromeLevel: 114}, leaderboards.LeaderboardList[1])
	assert.Equal(50, leaderboards.LeaderboardList[49].Rank)

	assert.Equal(LeaderboardPage{CurrentPage: 1, TotalPages: 21, TotalRecords: 1003}, leaderboards.LeaderboardPage)
}

func TestLeaderboardsLastPage(t *testing.T) {
	leaderboardsJson, err := TibiaLeaderboardsImpl("Antica", "138", 21, readLeaderboardsTestFile(t, "Antica-138-21.html"), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	leaderboards := leaderboardsJson.Leaderboards

	assert.Equal(3, len(leaderboards.LeaderboardList))
	assert.Equal(1001, leaderboards.LeaderboardList[0].Rank)
	assert.Equal(1003, leaderboards.LeaderboardList[2].Rank)
	assert.Equal(LeaderboardPage{CurrentPage: 21, TotalPages: 21, TotalRecords: 1003}, leaderboards.LeaderboardPage)
}

func TestLeaderboardsCurrent(t *testing.T) {
	leaderboardsJson, err := TibiaLeaderboardsImpl("Antica", "current", 1, readLeaderboardsTestFile(t, "Antica-current.html"), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	leaderboards := leaderboardsJson.Leaderboards

	assert.Equal(LeaderboardRotation{RotationID: 139, Current: true, EndDate: "2026-11-11T09:00:00Z"}, leaderboards.Rotation)

	lastUpdate, err := time.Parse(time.RFC3339, leaderboards.LastUpdate)
	assert.NoError(err)
	assert.True(lastUpdate.Before(time.Now()))
	assert.WithinDuration(time.Now().Add(-49*time.Minute), lastUpdate, time.Minute)

	// the "Last update" table must not be read as a ranking row
	assert.Equal(LeaderboardEntry{Rank: 1, Name: "Ranny Zolnierz Elder", DromeLevel: 110}, leaderboards.LeaderboardList[0])
	assert.Equal(LeaderboardPage{CurrentPage: 1, TotalPages: 3, TotalRecords: 109}, leaderboards.LeaderboardPage)
}

func TestLeaderboardsCurrentByID(t *testing.T) {
	leaderboardsJson, err := TibiaLeaderboardsImpl("Antica", "139", 1, readLeaderboardsTestFile(t, "Antica-current.html"), "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 139, leaderboardsJson.Leaderboards.Rotation.RotationID)
	assert.True(t, leaderboardsJson.Leaderboards.Rotation.Current)
}

func TestLeaderboardsUnknownRotation(t *testing.T) {
	// tibia.com serves the current rotation when asked for one it does not list
	_, err := TibiaLeaderboardsImpl("Antica", "100", 1, readLeaderboardsTestFile(t, "Antica-current.html"), "")
	assert.ErrorIs(t, err, validation.ErrorLeaderboardRotationDoesNotExist)
}

func TestLeaderboardsPageTooBig(t *testing.T) {
	// tibia.com serves the last page when asked for one past the end
	_, err := TibiaLeaderboardsImpl("Antica", "138", 22, readLeaderboardsTestFile(t, "Antica-138-21.html"), "")
	assert.ErrorIs(t, err, validation.ErrorHighscorePageTooBig)
}

func TestLeaderboardsEmptyCurrent(t *testing.T) {
	leaderboardsJson, err := TibiaLeaderboardsImpl("Antica", "current", 1, readLeaderboardsTestFile(t, "Antica-current-never.html"), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	leaderboards := leaderboardsJson.Leaderboards

	assert.Empty(leaderboards.LastUpdate)
	assert.NotNil(leaderboards.LeaderboardList)
	assert.Equal(0, len(leaderboards.LeaderboardList))
	assert.Equal(LeaderboardPage{CurrentPage: 1, TotalPages: 0, TotalRecords: 0}, leaderboards.LeaderboardPage)

	data, err := json.Marshal(leaderboardsJson)
	assert.NoError(err)
	assert.Contains(string(data), `"leaderboard_list":[]`)
}
