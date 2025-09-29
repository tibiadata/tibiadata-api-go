package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Child of Highscores
type Highscore struct {
	Rank     int    `json:"rank"`            // The character's rank/postition.
	Name     string `json:"name"`            // The name of the character.
	Vocation string `json:"vocation"`        // The character's vocation.
	World    string `json:"world"`           // The character's world.
	Level    int    `json:"level"`           // The character's level.
	Value    int    `json:"value"`           // The character's value for the highscores or loyalty points.
	Title    string `json:"title,omitempty"` // The character's loyalty title. (when category: loyalty)
}

// Child of Highscore
type HighscorePage struct {
	CurrentPage     int `json:"current_page"`  // The current page being displayed.
	TotalPages      int `json:"total_pages"`   // The total number of pages.
	TotalHighscores int `json:"total_records"` // The total amount of highscore records.
}

// Child of JSONData
type Highscores struct {
	World         string        `json:"world"`          // The world the highscores belong to.
	Category      string        `json:"category"`       // The selected category being displayed.
	Vocation      string        `json:"vocation"`       // The selected vocation filtered on.
	HighscoreAge  int           `json:"highscore_age"`  // The age of the highscore page in minutes.
	HighscoreList []Highscore   `json:"highscore_list"` // List of highscore records.
	HighscorePage HighscorePage `json:"highscore_page"` // Information of highscore pages.
}

// The base includes two levels: Highscores and Information
type HighscoresResponse struct {
	Highscores  Highscores  `json:"highscores"`
	Information Information `json:"information"`
}

var (
	HighscoresAgeRegex  = regexp.MustCompile(`.*<div class="Text">Highscores.*Last Update: ([0-9]+) minutes ago.*`)
	HighscoresPageRegex = regexp.MustCompile(`.*<b>.*Pages:\ ?(.*)<\/b>.*<b>.*Results:\ ?([0-9,]+)<\/b>.*`)
	SevenColumnRegex    = regexp.MustCompile(`<td>(.*)<\/td><td.*">(.*)<\/a><\/td><td.*>(.*)<\/td><td.*>(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
	SixColumnRegex      = regexp.MustCompile(`<td>(.*)<\/td><td.*">(.*)<\/a><\/td><td.*">(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
)

func TibiaHighscoresImpl(world string, category validation.HighscoreCategory, vocationName string, currentPage int, BoxContentHTML string, url string) (HighscoresResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return HighscoresResponse{}, fmt.Errorf("[error] TibiaHighscoresImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty HighscoreData var
	var (
		HighscoreData                                                                                                          []Highscore
		HighscoreDataVocation, HighscoreDataWorld, HighscoreDataTitle                                                          string
		HighscoreDataRank, HighscoreDataLevel, HighscoreDataValue, HighscoreAge, HighscoreTotalPages, HighscoreTotalHighscores int
	)

	// getting age of data
	subma1 := HighscoresAgeRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		HighscoreAge = TibiaDataStringToInteger(subma1[0][1])
	}

	// getting amount of pages
	subma1 = HighscoresPageRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		HighscoreTotalPages = strings.Count(subma1[0][1], "class=\"PageLink")
		HighscoreTotalHighscores = TibiaDataStringToInteger(subma1[0][2])
	}

	if currentPage > HighscoreTotalPages {
		return HighscoresResponse{}, validation.ErrorHighscorePageTooBig
	}

	var insideError error

	// Running query over each div
	ReaderHTML.Find(".TableContent tr").First().NextAll().EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		HighscoreDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaHighscoresImpl failed at HighscoreDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Regex the data table..
		var subma1 [][]string

		/*
			Tibia highscore table columns
			Achievment	=>	Rank		Name	Vocation	World		Level	Points
			Axe 		=>	Rank		Name	Vocation	World		Level	Skill Level
			Boss		=>	Rank		Name	Vocation	World		Level	Points
			Charm		=>	Rank		Name	Vocation	World		Level	Points
			Club		=>	Rank		Name	Vocation	World		Level	Skill Level
			Distance	=>	Rank		Name	Vocation	World		Level	Skill Level
			Drome		=>	Rank		Name	Vocation	World		Level	Score
			Exp 		=>	Rank		Name	Vocation	World		Level	Points
			Fishing		=>	Rank		Name	Vocation	World		Level	Skill Level
			Fist		=>	Rank		Name	Vocation	World		Level	Skill Level
			Goshnar		=>	Rank		Name	Vocation	World		Level	Points
			Loyalty		=>	Rank		Name	Title		Vocation	World	Level		Points
			Magic lvl	=>	Rank		Name	Vocation	World		Level	Skill Level
			Shield		=>	Rank		Name	Vocation	World		Level	Skill Level
			Sword		=>	Rank		Name	Vocation	World		Level	Skill Level
		*/

		if category == validation.HighScoreLoyaltypoints {
			subma1 = SevenColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		} else {
			subma1 = SixColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		}

		if len(subma1) > 0 {

			HighscoreDataRank = TibiaDataStringToInteger(subma1[0][1])
			if category == validation.HighScoreLoyaltypoints {
				HighscoreDataTitle = subma1[0][3]
				HighscoreDataVocation = subma1[0][4]
				HighscoreDataWorld = subma1[0][5]
				HighscoreDataLevel = TibiaDataStringToInteger(subma1[0][6])
				HighscoreDataValue = TibiaDataStringToInteger(subma1[0][7])
			} else {
				HighscoreDataVocation = subma1[0][3]
				HighscoreDataWorld = subma1[0][4]
				HighscoreDataLevel = TibiaDataStringToInteger(subma1[0][5])
				HighscoreDataValue = TibiaDataStringToInteger(subma1[0][6])
			}

			HighscoreData = append(HighscoreData, Highscore{
				Rank:     HighscoreDataRank,
				Name:     TibiaDataSanitizeEscapedString(subma1[0][2]),
				Vocation: HighscoreDataVocation,
				World:    HighscoreDataWorld,
				Level:    HighscoreDataLevel,
				Value:    HighscoreDataValue,
				Title:    HighscoreDataTitle,
			})
		}

		return true
	})

	categoryString, _ := category.String()

	if insideError != nil {
		return HighscoresResponse{}, insideError
	}

	//
	// Build the data-blob
	return HighscoresResponse{
		Highscores{
			World:         cases.Title(language.English).String(world),
			Category:      categoryString,
			Vocation:      vocationName,
			HighscoreAge:  HighscoreAge,
			HighscoreList: HighscoreData,
			HighscorePage: HighscorePage{
				CurrentPage:     currentPage,
				TotalPages:      HighscoreTotalPages,
				TotalHighscores: HighscoreTotalHighscores,
			},
		},
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{url},
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, err
}
