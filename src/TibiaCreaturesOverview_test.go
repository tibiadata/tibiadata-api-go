package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestOverview(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creatures.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	creaturesJson, err := TibiaCreaturesOverviewImpl(string(data), "https://www.tibia.com/library/?subtopic=creatures")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := creaturesJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=creatures", information.TibiaURLs[0])
	assert.Equal("Feral Werecrocodile", creaturesJson.Creatures.Boosted.Name)
	assert.Equal("feralwerecrocodile", creaturesJson.Creatures.Boosted.Race)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/feralwerecrocodile.gif", creaturesJson.Creatures.Boosted.ImageURL)

	assert.Equal(638, len(creaturesJson.Creatures.Creatures))

	acidblob := creaturesJson.Creatures.Creatures[0]
	assert.Equal("Acid Blobs", acidblob.Name)
	assert.Equal("acidblob", acidblob.Race)
	assert.Equal("https://static.tibia.com/images/library/acidblob.gif", acidblob.ImageURL)
	assert.False(acidblob.Featured)

	feralwerecrocodile := creaturesJson.Creatures.Creatures[204]
	assert.Equal("Feral Werecrocodiles", feralwerecrocodile.Name)
	assert.Equal("feralwerecrocodile", feralwerecrocodile.Race)
	assert.Equal("https://static.tibia.com/images/library/feralwerecrocodile.gif", feralwerecrocodile.ImageURL)
	assert.True(feralwerecrocodile.Featured)

	minotauramazon := creaturesJson.Creatures.Creatures[360]
	assert.Equal("Minotaur Amazons", minotauramazon.Name)
	assert.Equal("minotauramazon", minotauramazon.Race)
	assert.Equal("https://static.tibia.com/images/library/minotauramazon.gif", minotauramazon.ImageURL)
	assert.False(minotauramazon.Featured)

	quarapredator := creaturesJson.Creatures.Creatures[465]
	assert.Equal("Quara Predators", quarapredator.Name)
	assert.Equal("quarapredator", quarapredator.Race)
	assert.Equal("https://static.tibia.com/images/library/quarapredator.gif", quarapredator.ImageURL)
	assert.False(quarapredator.Featured)

	slimes := creaturesJson.Creatures.Creatures[510]
	assert.Equal("Slimes", slimes.Name)
	assert.Equal("slime", slimes.Race)
	assert.Equal("https://static.tibia.com/images/library/slime.gif", slimes.ImageURL)
	assert.False(slimes.Featured)
}

func TestOverviewBoostedFictitious(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creatures_fictitious.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	creaturesJson, err := TibiaCreaturesOverviewImpl(string(data), "https://www.tibia.com/library/?subtopic=creatures")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := creaturesJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=creatures", information.TibiaURLs[0])
	assert.Equal("Ragged Rabid Wolf", creaturesJson.Creatures.Boosted.Name)
	assert.Equal("raggedrabidwolf", creaturesJson.Creatures.Boosted.Race)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/raggedrabidwolf.gif", creaturesJson.Creatures.Boosted.ImageURL)
}
