package validation

import (
	"errors"
	"strings"
)

var (
	// validHighscoreCatregories stores all valid highscore categories
	validHighscoreCategories = []string{"achievements", "achievement", "axe", "axefighting", "charm", "charms", "charmpoints", "charmspoints", "club", "clubfighting", "distance", "distancefighting", "fishing", "fist", "fistfighting", "goshnar", "goshnars", "goshnarstaint", "loyalty", "loyaltypoints", "magic", "mlvl", "magiclevel", "shielding", "shield", "sword", "swordfighting", "drome", "dromescore", "experience", "boss", "bosses", "bosspoints", "bountypoints", "bountypoint", "bountypointsearned", "weeklytasks", "weeklytask", "weeklytaskscompleted"}
)

// IsHighscoreCategoryValid reports wheter the provided string represents a valid highscore category
// Check if error == nil to see whether the highscore category is valid or not
func IsHighscoreCategoryValid(hs string) error {
	for _, highscore := range validHighscoreCategories {
		if strings.EqualFold(hs, highscore) {
			return nil
		}
	}

	return ErrorHighscoreCategoryDoesNotExist
}

type HighscoreCategory int

const (
	HighScoreAchievements HighscoreCategory = iota + 1
	HighScoreAxefighting
	HighScoreCharmpoints
	HighScoreClubfighting
	HighScoreDistancefighting
	HighScoreExperience
	HighScoreFishing
	HighScoreFistfighting
	HighScoreGoshnarstaint
	HighScoreLoyaltypoints
	HighScoreMagiclevel
	HighScoreShielding
	HighScoreSwordfighting
	HighScoreDromescore
	HighScoreBosspoints
	HighScoreBountypoints
	HighScoreWeeklytasks
)

func (hc HighscoreCategory) String() (string, error) {
	seasons := [...]string{"achievements", "axefighting", "charmpoints", "clubfighting", "distancefighting", "experience", "fishing", "fistfighting", "goshnarstaint", "loyaltypoints", "magiclevel", "shielding", "swordfighting", "dromescore", "bosspoints"}
	if hc < HighScoreAchievements || hc > HighScoreBosspoints {
		return "", errors.New("invalid HighscoreCategory value")
	}
	return seasons[hc-1], nil
}

func HighscoreCategoryFromString(input string) HighscoreCategory {
	// Sanatize of category value
	input = strings.ToLower(input)
	switch input {
	case "achievements", "achievement":
		return HighScoreAchievements
	case "axe", "axefighting":
		return HighScoreAxefighting
	case "charm", "charms", "charmpoints", "charmspoints":
		return HighScoreCharmpoints
	case "club", "clubfighting":
		return HighScoreClubfighting
	case "distance", "distancefighting":
		return HighScoreDistancefighting
	case "fishing":
		return HighScoreFishing
	case "fist", "fistfighting":
		return HighScoreFistfighting
	case "goshnar", "goshnars", "goshnarstaint":
		return HighScoreGoshnarstaint
	case "loyalty", "loyaltypoints":
		return HighScoreLoyaltypoints
	case "magic", "mlvl", "magiclevel":
		return HighScoreMagiclevel
	case "shielding", "shield":
		return HighScoreShielding
	case "sword", "swordfighting":
		return HighScoreSwordfighting
	case "drome", "dromescore":
		return HighScoreDromescore
	case "boss", "bosses", "bosspoints":
		return HighScoreBosspoints
	case "bountypoints", "bountypoint", "bountypointsearned":
		return HighScoreBountypoints
	case "weeklytasks", "weeklytask", "weeklytaskscompleted":
		return HighScoreWeeklytasks
	default:
		return HighScoreExperience
	}
}
