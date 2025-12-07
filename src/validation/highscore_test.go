package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighscoreCategoryAchievementsString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreAchievements
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("achievements", stringValue)
	assert.Equal(HighscoreCategory(1), highscoreCategory)
}

func TestHighscoreCategoryExperienceString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreExperience
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("experience", stringValue)
	assert.Equal(HighscoreCategory(6), highscoreCategory)
}

func TestHighscoreCategoryDromescoreString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreDromescore
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("dromescore", stringValue)
	assert.Equal(HighscoreCategory(14), highscoreCategory)
}

func TestHighscoreCategoryBosspointsString(t *testing.T) {
	assert := assert.New(t)
	highscoreCategory := HighScoreBosspoints
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("bosspoints", stringValue)
	assert.Equal(HighscoreCategory(15), highscoreCategory)
}

func TestHighscoreCategoryBountypointsString(t *testing.T) {
	assert := assert.New(t)
	highscoreCategory := HighScoreBountypoints
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("bountypoints", stringValue)
	assert.Equal(HighscoreCategory(16), highscoreCategory)
}

func TestHighscoreCategoryWeeklytasksString(t *testing.T) {
	assert := assert.New(t)
	highscoreCategory := HighScoreWeeklytasks
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("weeklytasks", stringValue)
	assert.Equal(HighscoreCategory(17), highscoreCategory)
}

func TestHighscoreCategoryInvalidValueString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighscoreCategory(99)
	_, err := highscoreCategory.String()

	assert.NotNil(err)
}

func TestHighscoreCategoryFromString(t *testing.T) {
	assert := assert.New(t)

	categoryTests := map[string]struct {
		inputs   []string
		expected HighscoreCategory
	}{
		"Experience": {
			inputs:   []string{"experience", ""},
			expected: HighScoreExperience,
		},
		"Achievements": {
			inputs:   []string{"achievements", "achievement"},
			expected: HighScoreAchievements,
		},
		"Axefighting": {
			inputs:   []string{"axe", "axefighting"},
			expected: HighScoreAxefighting,
		},
		"Charmpoints": {
			inputs:   []string{"charm", "charms", "charmpoints", "charmspoints"},
			expected: HighScoreCharmpoints,
		},
		"Clubfighting": {
			inputs:   []string{"club", "clubfighting"},
			expected: HighScoreClubfighting,
		},
		"Distancefighting": {
			inputs:   []string{"distance", "distancefighting"},
			expected: HighScoreDistancefighting,
		},
		"Fishing": {
			inputs:   []string{"fishing"},
			expected: HighScoreFishing,
		},
		"Fistfighting": {
			inputs:   []string{"fist", "fistfighting"},
			expected: HighScoreFistfighting,
		},
		"Goshnarstaint": {
			inputs:   []string{"goshnar", "goshnars", "goshnarstaint"},
			expected: HighScoreGoshnarstaint,
		},
		"Loyaltypoints": {
			inputs:   []string{"loyalty", "loyaltypoints"},
			expected: HighScoreLoyaltypoints,
		},
		"Magiclevel": {
			inputs:   []string{"magic", "mlvl", "magiclevel"},
			expected: HighScoreMagiclevel,
		},
		"Shielding": {
			inputs:   []string{"shielding", "shield"},
			expected: HighScoreShielding,
		},
		"Swordfighting": {
			inputs:   []string{"sword", "swordfighting"},
			expected: HighScoreSwordfighting,
		},
		"Dromescore": {
			inputs:   []string{"drome", "dromescore"},
			expected: HighScoreDromescore,
		},
		"Bosspoints": {
			inputs:   []string{"boss", "bosses", "bosspoints"},
			expected: HighScoreBosspoints,
		},
		"Bountypoints": {
			inputs:   []string{"bountypoints", "bountypoint", "bountypointsearned"},
			expected: HighscoreBountypoints,
		},
		"Weeklytasks": {
			inputs:   []string{"weeklytasks", "weeklytask", "weeklytaskscompleted"},
			expected: HighscoreWeeklytasks,
		},
	}

	for category, data := range categoryTests {
		t.Run(category, func(t *testing.T) {
			for _, input := range data.inputs {
				t.Run(input, func(t *testing.T) {
					result := HighscoreCategoryFromString(input)
					assert.Equal(data.expected, result, "unexpected result for input: %s", input)
				})
			}
		})
	}
}
