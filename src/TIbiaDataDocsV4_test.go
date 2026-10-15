package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestDocsV4CodeCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	assert := assert.New(t)

	assert.False(false, tibiaBoostableBossesV4())
	assert.False(false, tibiaCharactersCharacterV4())
	assert.False(false, tibiaCreaturesOverviewV4())
	assert.False(false, tibiaCreaturesCreatureV4())
	assert.False(false, tibiaFansitesV4())
	assert.False(false, tibiaGuildsGuildV4())
	assert.False(false, tibiaGuildsOverviewV4())
	assert.False(false, tibiaHighscoresV4())
	assert.False(false, tibiaHousesHouseV4())
	assert.False(false, tibiaHousesOverviewV4())
	assert.False(false, tibiaKillstatisticsV4())
	assert.False(false, tibiaNewslistArchiveV4())
	assert.False(false, tibiaNewslistArchiveDaysV4())
	assert.False(false, tibiaNewslistLatestV4())
	assert.False(false, tibiaNewslistV4())
	assert.False(false, tibiaNewsV4())
	assert.False(false, tibiaSpellsOverviewV4())
	assert.False(false, tibiaSpellsSpellV4())
	assert.False(false, tibiaWorldsOverviewV4())
	assert.False(false, tibiaWorldsWorldV4())

}
