package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Child of Leaderboards
type LeaderboardRotation struct {
	RotationID int    `json:"rotation_id"` // The ID of the Tibiadrome rotation.
	Current    bool   `json:"current"`     // Whether this is the currently running rotation.
	EndDate    string `json:"end_date"`    // The time the rotation ends or ended.
}

// Child of Leaderboards
type LeaderboardEntry struct {
	Rank       int    `json:"rank"`        // The character's rank/position.
	Name       string `json:"name"`        // The name of the character.
	DromeLevel int    `json:"drome_level"` // The character's drome level.
}

// Child of Leaderboards
type LeaderboardPage struct {
	CurrentPage  int `json:"current_page"`  // The current page being displayed.
	TotalPages   int `json:"total_pages"`   // The total number of pages.
	TotalRecords int `json:"total_records"` // The total amount of leaderboard records.
}

// Child of JSONData
type Leaderboards struct {
	World              string                `json:"world"`               // The world the leaderboard belongs to.
	Rotation           LeaderboardRotation   `json:"rotation"`            // The rotation being displayed.
	AvailableRotations []LeaderboardRotation `json:"available_rotations"` // The rotations tibia.com offers.
	LastUpdate         string                `json:"last_update"`         // The time the leaderboard was last updated (empty if never or not shown).
	LeaderboardList    []LeaderboardEntry    `json:"leaderboard_list"`    // List of leaderboard records.
	LeaderboardPage    LeaderboardPage       `json:"leaderboard_page"`    // Information of leaderboard pages.
}

// The base includes two levels: Leaderboards and Information
type LeaderboardsResponse struct {
	Leaderboards Leaderboards `json:"leaderboards"`
	Information  Information  `json:"information"`
}

var (
	LeaderboardsLastUpdateRegex = regexp.MustCompile(`Last update: ([0-9]+) minutes? ago`)
	LeaderboardsResultsRegex    = regexp.MustCompile(`Results:\s*([0-9,]+)`)
)

const leaderboardsCurrentRotationPrefix = "Current rotation (ends on "

func TibiaLeaderboardsImpl(world, rotation string, currentPage int, BoxContentHTML, url string) (LeaderboardsResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return LeaderboardsResponse{}, fmt.Errorf("[error] TibiaLeaderboardsImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var (
		availableRotations []LeaderboardRotation
		selectedRotation   LeaderboardRotation
	)

	ReaderHTML.Find(`select[name="rotation"] option`).Each(func(index int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		current := strings.HasPrefix(text, leaderboardsCurrentRotationPrefix)
		if current {
			text = strings.TrimSuffix(strings.TrimPrefix(text, leaderboardsCurrentRotationPrefix), ")")
		}

		id, _ := s.Attr("value")
		leaderboardRotation := LeaderboardRotation{
			RotationID: TibiaDataStringToInteger(id),
			Current:    current,
			EndDate:    TibiaDataDatetime(text),
		}
		availableRotations = append(availableRotations, leaderboardRotation)

		// tibia.com marks the shown rotation; without a mark the first (current) one is shown
		if _, selected := s.Attr("selected"); selected || index == 0 {
			selectedRotation = leaderboardRotation
		}
	})

	// Without the rotation list we cannot tell which rotation tibia.com served
	if len(availableRotations) == 0 {
		return LeaderboardsResponse{}, fmt.Errorf("[error] TibiaLeaderboardsImpl found no rotation list")
	}

	if rotation != "current" {
		requestedRotation, err := strconv.Atoi(rotation)
		if err != nil || requestedRotation != selectedRotation.RotationID {
			return LeaderboardsResponse{}, validation.ErrorLeaderboardRotationDoesNotExist
		}
	}

	var lastUpdate string
	if subma1 := LeaderboardsLastUpdateRegex.FindStringSubmatch(BoxContentHTML); len(subma1) > 0 {
		minutes := TibiaDataStringToInteger(subma1[1])
		lastUpdate = time.Now().Add(-time.Duration(minutes) * time.Minute).UTC().Format(time.RFC3339)
	}

	var totalPages, totalRecords int
	if pageNavigation := ReaderHTML.Find(".PageNavigation").First(); pageNavigation.Length() > 0 {
		totalPages = pageNavigation.Find(".PageLink").Length()
		if subma1 := LeaderboardsResultsRegex.FindStringSubmatch(pageNavigation.Text()); len(subma1) > 0 {
			totalRecords = TibiaDataStringToInteger(subma1[1])
		}
	}

	if totalPages > 0 && currentPage > totalPages {
		return LeaderboardsResponse{}, validation.ErrorHighscorePageTooBig
	}

	leaderboardList := make([]LeaderboardEntry, 0, 50)

	// Only the table with the Rank/Name/Drome Level header holds ranking rows
	ReaderHTML.Find(".TableContent").Has("tr.LabelH").First().Find("tr.LabelH").NextAll().Each(func(index int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() < 3 {
			return
		}

		leaderboardList = append(leaderboardList, LeaderboardEntry{
			Rank:       TibiaDataStringToInteger(strings.TrimSuffix(strings.TrimSpace(cells.Eq(0).Text()), ".")),
			Name:       strings.TrimSpace(cells.Eq(1).Text()),
			DromeLevel: TibiaDataStringToInteger(strings.TrimSpace(cells.Eq(2).Text())),
		})
	})

	//
	// Build the data-blob
	return LeaderboardsResponse{
		Leaderboards{
			World:              world,
			Rotation:           selectedRotation,
			AvailableRotations: availableRotations,
			LastUpdate:         lastUpdate,
			LeaderboardList:    leaderboardList,
			LeaderboardPage: LeaderboardPage{
				CurrentPage:  currentPage,
				TotalPages:   totalPages,
				TotalRecords: totalRecords,
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
	}, nil
}
