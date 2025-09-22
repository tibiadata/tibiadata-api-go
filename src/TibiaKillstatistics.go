package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Child of KillStatistics
type Entry struct {
	Race                    string `json:"race"`                     // The name of the creature/race.
	LastDayKilledPlayers    int    `json:"last_day_players_killed"`  // Number of players killed by this race in the last day.
	LastDayKilledByPlayers  int    `json:"last_day_killed"`          // Number of creatures of this race killed in the last day.
	LastWeekKilledPlayers   int    `json:"last_week_players_killed"` // Number of players killed by this race in the last week.
	LastWeekKilledByPlayers int    `json:"last_week_killed"`         // Number of creatures of this race killed in the last week.
}

// Child of KillStatistics
type Total struct {
	LastDayKilledPlayers    int `json:"last_day_players_killed"`  // Total number of players killed in total in the last day.
	LastDayKilledByPlayers  int `json:"last_day_killed"`          // Total number of creatures in total killed in the last day.
	LastWeekKilledPlayers   int `json:"last_week_players_killed"` // Total number of players killed in total in the last week.
	LastWeekKilledByPlayers int `json:"last_week_killed"`         // Total number of creatures in total killed in the last week.
}

// Child of JSONData
type KillStatistics struct {
	World   string  `json:"world"`   // The world the statistics belong to.
	Entries []Entry `json:"entries"` // List of killstatistic.
	Total   Total   `json:"total"`   // List of total kills.
}

// The base includes two levels: KillStatistics and Information
type KillStatisticsResponse struct {
	KillStatistics KillStatistics `json:"killstatistics"`
	Information    Information    `json:"information"`
}

func TibiaKillstatisticsImpl(world string, BoxContentHTML string, url string) (KillStatisticsResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return KillStatisticsResponse{}, fmt.Errorf("[error] TibiaKillstatisticsImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty KillStatisticsData var
	var (
		KillStatisticsData                                                                                               []Entry
		TotalLastDayKilledPlayers, TotalLastDayKilledByPlayers, TotalLastWeekKilledPlayers, TotalLastWeekKilledByPlayers int
	)

	// Running query over each div
	ReaderHTML.Find("#KillStatisticsTable .TableContent tr.Odd,tr.Even").Each(func(index int, s *goquery.Selection) {
		// Extract kill statistics from table row
		killStats := extractKillStatistics(s.Find("td").Nodes)

		// Accumulate totals
		TotalLastDayKilledPlayers += killStats.LastDayKilledPlayers
		TotalLastDayKilledByPlayers += killStats.LastDayKilledByPlayers
		TotalLastWeekKilledPlayers += killStats.LastWeekKilledPlayers
		TotalLastWeekKilledByPlayers += killStats.LastWeekKilledByPlayers

		// Append new Entry item to KillStatisticsData
		KillStatisticsData = append(KillStatisticsData, killStats)
	})

	//
	// Build the data-blob
	return KillStatisticsResponse{
		KillStatistics{
			World:   world,
			Entries: KillStatisticsData,
			Total: Total{
				LastDayKilledPlayers:    TotalLastDayKilledPlayers,
				LastDayKilledByPlayers:  TotalLastDayKilledByPlayers,
				LastWeekKilledPlayers:   TotalLastWeekKilledPlayers,
				LastWeekKilledByPlayers: TotalLastWeekKilledByPlayers,
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

// Helper function to extract and convert kill statistics
func extractKillStatistics(dataColumns []*html.Node) Entry {
	return Entry{
		Race:                    TibiaDataSanitizeEscapedString(dataColumns[0].FirstChild.Data),
		LastDayKilledPlayers:    TibiaDataStringToInteger(dataColumns[1].FirstChild.Data),
		LastDayKilledByPlayers:  TibiaDataStringToInteger(dataColumns[2].FirstChild.Data),
		LastWeekKilledPlayers:   TibiaDataStringToInteger(dataColumns[3].FirstChild.Data),
		LastWeekKilledByPlayers: TibiaDataStringToInteger(dataColumns[4].FirstChild.Data),
	}
}
