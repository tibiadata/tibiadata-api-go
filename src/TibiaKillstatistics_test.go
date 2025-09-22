package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
	"golang.org/x/net/html"
)

func TestAntica(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/killstatistics/Antica.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	anticaJson, err := TibiaKillstatisticsImpl("Antica", string(data), "https://www.tibia.com/community/?subtopic=killstatistics&world=Antica")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := anticaJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=killstatistics&world=Antica", information.TibiaURLs[0])

	assert.Equal("Antica", anticaJson.KillStatistics.World)
	assert.Equal(1159, len(anticaJson.KillStatistics.Entries))

	elementalForces := anticaJson.KillStatistics.Entries[0]
	assert.Equal("(elemental forces)", elementalForces.Race)
	assert.Equal(6, elementalForces.LastDayKilledPlayers)
	assert.Equal(0, elementalForces.LastDayKilledByPlayers)
	assert.Equal(103, elementalForces.LastWeekKilledPlayers)
	assert.Equal(0, elementalForces.LastWeekKilledByPlayers)

	caveRats := anticaJson.KillStatistics.Entries[386]
	assert.Equal("cave rats", caveRats.Race)
	assert.Equal(2, caveRats.LastDayKilledPlayers)
	assert.Equal(1618, caveRats.LastDayKilledByPlayers)
	assert.Equal(78, caveRats.LastWeekKilledPlayers)
	assert.Equal(12876, caveRats.LastWeekKilledByPlayers)
}

func BenchmarkAntica(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/killstatistics/Antica.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	assert := assert.New(b)

	for i := 0; i < b.N; i++ {
		anticaJson, err := TibiaKillstatisticsImpl("Antica", string(data), "")
		if err != nil {
			b.Fatal(err)
		}

		assert.Equal("Antica", anticaJson.KillStatistics.World)
	}
}

func TestExtractKillStatistics(t *testing.T) {
	// Create mock HTML nodes for testing
	raceNode := &html.Node{
		FirstChild: &html.Node{
			Data: "Dragon",
		},
	}

	lastDayKilledPlayersNode := &html.Node{
		FirstChild: &html.Node{
			Data: "150",
		},
	}

	lastDayKilledByPlayersNode := &html.Node{
		FirstChild: &html.Node{
			Data: "25",
		},
	}

	lastWeekKilledPlayersNode := &html.Node{
		FirstChild: &html.Node{
			Data: "1200",
		},
	}

	lastWeekKilledByPlayersNode := &html.Node{
		FirstChild: &html.Node{
			Data: "85",
		},
	}

	// Create dataColumns slice
	dataColumns := []*html.Node{
		raceNode,
		lastDayKilledPlayersNode,
		lastDayKilledByPlayersNode,
		lastWeekKilledPlayersNode,
		lastWeekKilledByPlayersNode,
	}

	// Call the function
	result := extractKillStatistics(dataColumns)

	// Define expected values
	expected := Entry{
		Race:                    "Dragon",
		LastDayKilledPlayers:    150,
		LastDayKilledByPlayers:  25,
		LastWeekKilledPlayers:   1200,
		LastWeekKilledByPlayers: 85,
	}

	// Assert results
	if result.Race != expected.Race {
		t.Errorf("Race: got %v, want %v", result.Race, expected.Race)
	}

	if result.LastDayKilledPlayers != expected.LastDayKilledPlayers {
		t.Errorf("LastDayKilledPlayers: got %v, want %v", result.LastDayKilledPlayers, expected.LastDayKilledPlayers)
	}

	if result.LastDayKilledByPlayers != expected.LastDayKilledByPlayers {
		t.Errorf("LastDayKilledByPlayers: got %v, want %v", result.LastDayKilledByPlayers, expected.LastDayKilledByPlayers)
	}

	if result.LastWeekKilledPlayers != expected.LastWeekKilledPlayers {
		t.Errorf("LastWeekKilledPlayers: got %v, want %v", result.LastWeekKilledPlayers, expected.LastWeekKilledPlayers)
	}

	if result.LastWeekKilledByPlayers != expected.LastWeekKilledByPlayers {
		t.Errorf("LastWeekKilledByPlayers: got %v, want %v", result.LastWeekKilledByPlayers, expected.LastWeekKilledByPlayers)
	}
}

// Test with edge cases
func TestExtractKillStatisticsEdgeCases(t *testing.T) {
	// Test with empty strings
	emptyNode := &html.Node{
		FirstChild: &html.Node{
			Data: "",
		},
	}

	// Test with zero values
	zeroNode := &html.Node{
		FirstChild: &html.Node{
			Data: "0",
		},
	}

	dataColumns := []*html.Node{
		emptyNode, // Race
		zeroNode,  // LastDayKilledPlayers
		zeroNode,  // LastDayKilledByPlayers
		zeroNode,  // LastWeekKilledPlayers
		zeroNode,  // LastWeekKilledByPlayers
	}

	result := extractKillStatistics(dataColumns)

	if result.Race != "" {
		t.Errorf("Empty race: got %v, want empty string", result.Race)
	}

	if result.LastDayKilledPlayers != 0 {
		t.Errorf("Zero last day killed players: got %v, want 0", result.LastDayKilledPlayers)
	}

	if result.LastWeekKilledByPlayers != 0 {
		t.Errorf("Zero last week killed by players: got %v, want 0", result.LastWeekKilledByPlayers)
	}
}
