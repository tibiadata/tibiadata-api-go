package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of Guilds
type OverviewGuild struct {
	Name        string `json:"name"`        // The name of the guild.
	LogoURL     string `json:"logo_url"`    // The URL to the guild's logo.
	Description string `json:"description"` // The description of the guild.
}

// Child of JSONData
type OverviewGuilds struct {
	World     string          `json:"world"`     // The world the guilds belongs to.
	Active    []OverviewGuild `json:"active"`    // List of active guilds.
	Formation []OverviewGuild `json:"formation"` // List of guilds under formation.
}

// The base includes two levels: Guilds and Information
type GuildsOverviewResponse struct {
	Guilds      OverviewGuilds `json:"guilds"`
	Information Information    `json:"information"`
}

func TibiaGuildsOverviewImpl(world string, BoxContentHTML string, url string) (GuildsOverviewResponse, error) {
	// Creating empty vars
	var (
		ActiveGuilds, FormationGuilds []OverviewGuild
		GuildCategory                 string
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return GuildsOverviewResponse{}, fmt.Errorf("[error] TibiaGuildsOverviewImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Running query over each div
	ReaderHTML.Find(".TableContainer").Each(func(index int, s *goquery.Selection) {
		// Figure out the guild category
		s.Find(".Text").Each(func(index int, s *goquery.Selection) {
			tableName := s.Nodes[0].FirstChild.Data
			switch {
			case strings.Contains(tableName, "Active Guilds"):
				GuildCategory = "active"
			case strings.Contains(tableName, "Guilds in Course of Formation"):
				GuildCategory = "formation"
			}
		})

		if GuildCategory != "" {
			// Extract guilds
			s.Find(".TableContent tbody").Children().NextAll().Each(func(index int, s *goquery.Selection) {
				tableRow := s.Nodes[0]
				nameAndDescriptionNode := tableRow.FirstChild.NextSibling.NextSibling

				name := nameAndDescriptionNode.FirstChild.FirstChild.Data
				logoURL := tableRow.FirstChild.FirstChild.Attr[0].Val
				description := ""

				// Check if there's a description to fetch.
				if nameAndDescriptionNode.FirstChild.NextSibling != nil && nameAndDescriptionNode.FirstChild.NextSibling.NextSibling != nil {
					description = strings.TrimSpace(nameAndDescriptionNode.FirstChild.NextSibling.NextSibling.Data)
				}

				OneGuild := OverviewGuild{
					Name:        name,
					LogoURL:     logoURL,
					Description: description,
				}

				// Adding OneGuild to correct category
				switch GuildCategory {
				case "active":
					ActiveGuilds = append(ActiveGuilds, OneGuild)
				case "formation":
					FormationGuilds = append(FormationGuilds, OneGuild)
				}
			})
		}
	})

	//
	// Build the data-blob
	return GuildsOverviewResponse{
		OverviewGuilds{
			World:     world,
			Active:    ActiveGuilds,
			Formation: FormationGuilds,
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
