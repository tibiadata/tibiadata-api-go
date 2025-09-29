package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestPremia(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/Premia.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	premiaGuildsJson, err := TibiaGuildsOverviewImpl("Premia", string(data), "https://www.tibia.com/community/?subtopic=guilds&world=Premia")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	information := premiaGuildsJson.Information

	assert.Equal("Premia", premiaGuildsJson.Guilds.World)
	assert.Equal(38, len(premiaGuildsJson.Guilds.Active))
	assert.Equal(3, len(premiaGuildsJson.Guilds.Formation))

	orderOfGloryActiveGuild := premiaGuildsJson.Guilds.Active[28]
	assert.Equal("Order of Glory", orderOfGloryActiveGuild.Name)
	assert.Equal("https://static.tibia.com/images/guildlogos/Order_of_Glory.gif", orderOfGloryActiveGuild.LogoURL)
	assert.Equal("We are an English speaking guild of friends and allies from around the world who seek only peaceful questing, exploring, team hunts and a chill place to hang out. Message any of our leaders for an invitation. Contact Zyb with any problems.", orderOfGloryActiveGuild.Description)

	secondGuildInFormation := premiaGuildsJson.Guilds.Formation[1]
	assert.Equal("Konungen", secondGuildInFormation.Name)
	assert.Equal("https://static.tibia.com/images/community/default_logo.gif", secondGuildInFormation.LogoURL)
	assert.Empty(secondGuildInFormation.Description)

	assert.Equal("https://www.tibia.com/community/?subtopic=guilds&world=Premia", information.TibiaURLs[0])
}
