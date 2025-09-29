package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

func TestHighscoresAll(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/highscores/all.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	highscoresJson, err := TibiaHighscoresImpl("", validation.HighScoreExperience, "all", 1, string(data), "https://www.tibia.com/community/?subtopic=highscores&world=&category=experience&profession=all&currentpage=1")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := highscoresJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=highscores&world=&category=experience&profession=all&currentpage=1", information.TibiaURLs[0])

	assert.Empty(highscoresJson.Highscores.World)

	assert.Equal("experience", highscoresJson.Highscores.Category)
	assert.Equal("all", highscoresJson.Highscores.Vocation)
	assert.Equal(12, highscoresJson.Highscores.HighscoreAge)

	assert.Equal(50, len(highscoresJson.Highscores.HighscoreList))

	assert.Equal(1, highscoresJson.Highscores.HighscorePage.CurrentPage)
	assert.Equal(20, highscoresJson.Highscores.HighscorePage.TotalPages)
	assert.Equal(1000, highscoresJson.Highscores.HighscorePage.TotalHighscores)

	firstHighscore := highscoresJson.Highscores.HighscoreList[0]
	assert.Equal(1, firstHighscore.Rank)
	assert.Equal("Goraca", firstHighscore.Name)
	assert.Equal("Master Sorcerer", firstHighscore.Vocation)
	assert.Equal("Bona", firstHighscore.World)
	assert.Equal(2197, firstHighscore.Level)
	assert.Equal(176271164607, firstHighscore.Value)
	assert.Empty(firstHighscore.Title)

	lastHighscore := highscoresJson.Highscores.HighscoreList[49]
	assert.Equal(50, lastHighscore.Rank)
	assert.Equal("Wujo Daro", lastHighscore.Name)
	assert.Equal("Elite Knight", lastHighscore.Vocation)
	assert.Equal("Refugia", lastHighscore.World)
	assert.Equal(1701, lastHighscore.Level)
	assert.Equal(81816135617, lastHighscore.Value)
	assert.Empty(lastHighscore.Title)
}

func TestHighscoresLoyalty(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/highscores/loyalty.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	highscoresJson, err := TibiaHighscoresImpl("Vunira", validation.HighScoreLoyaltypoints, "druids", 4, string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Vunira", highscoresJson.Highscores.World)
	assert.Equal("loyaltypoints", highscoresJson.Highscores.Category)
	assert.Equal("druids", highscoresJson.Highscores.Vocation)
	assert.Equal(12, highscoresJson.Highscores.HighscoreAge)

	assert.Equal(50, len(highscoresJson.Highscores.HighscoreList))
}
