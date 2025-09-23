package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Child of Guild
type Guildhall struct {
	Name  string `json:"name"`  // The name of the house.
	World string `json:"world"` // The world the guildhall belongs to.
	/*
		Town      string `json:"town"`       // We can collect that from cached info?
		Status    string `json:"status"`     // rented (but maybe also auctioned)
		Owner     string `json:"owner"`      // We can collect that from cached info?
		HouseID   int    `json:"houseid"`    // We can collect that from cached info?
	*/
	PaidUntil string `json:"paid_until"` // The date the last paid rent is due.
}

// Child of Guild
type GuildMember struct {
	Name     string `json:"name"`     // The name of the guild's member.
	Title    string `json:"title"`    // The member's title.
	Rank     string `json:"rank"`     // The rank the member does belong to.
	Vocation string `json:"vocation"` // The member's vocation.
	Level    int    `json:"level"`    // The member's level.
	Joined   string `json:"joined"`   // The day when the member joined.
	Status   string `json:"status"`   // Whether the member is online or offline.
}

// Child of Guild
type InvitedGuildMember struct {
	Name string `json:"name"` // The name of the character.
	Date string `json:"date"` // The date the character was invited.
}

// Child of JSONData
type Guild struct {
	Name               string               `json:"name"`              // The name of the guild.
	World              string               `json:"world"`             // The world the guild belongs to.
	LogoURL            string               `json:"logo_url"`          // The URL to the guild's logo.
	Description        string               `json:"description"`       // The description of the guild.
	Guildhalls         []Guildhall          `json:"guildhalls"`        // The guildhall the guild has as their home.
	Active             bool                 `json:"active"`            // Whether the guild is active or in formation.
	Founded            string               `json:"founded"`           // The day it was founded.
	Applications       bool                 `json:"open_applications"` // Whether applications are open or not.
	Homepage           string               `json:"homepage"`          // The guild's homepage.
	InWar              bool                 `json:"in_war"`            // Whether it is currently in war or not.
	DisbandedDate      string               `json:"disband_date"`      // The date when the guild will be disbanded, if the condition aren't meet.
	DisbandedCondition string               `json:"disband_condition"` // The reason why the guild will get disbanded.
	PlayersOnline      int                  `json:"players_online"`    // The number of online members in the guild.
	PlayersOffline     int                  `json:"players_offline"`   // The number of offline members in the guild.
	MembersTotal       int                  `json:"members_total"`     // The number of total members in the guild.
	MembersInvited     int                  `json:"members_invited"`   // The number of invited members in the guild.
	Members            []GuildMember        `json:"members"`           // List of all members in the guild.
	Invited            []InvitedGuildMember `json:"invites"`           // List of invited members.
}

// The base includes two levels: Guild and Information
type GuildResponse struct {
	Guild       Guild       `json:"guild"`
	Information Information `json:"information"`
}

var (
	GuildLogoRegex                     = regexp.MustCompile(`.*img src="(.*)" width=.*`)
	GuildWorldAndFoundationRegex       = regexp.MustCompile(`^The guild was founded on (.*) on (.*).<br/>`)
	GuildHomepageRegex                 = regexp.MustCompile(`<a href="(.*)" target=.*>`)
	GuildhallRegex                     = regexp.MustCompile(` is (.*). The rent is paid until (.*).<br/>`)
	GuildDisbaneRegex                  = regexp.MustCompile(`<b>It will be disbanded on (.*.[0-9]+.[0-9]+) (.*)\.<\/b>.*`)
	GuildMemberInformationRegex        = regexp.MustCompile(`<td>(.*)<\/td><td><a.*">(.*)<\/a>(.*)<\/td><td>(.*)<\/td><td>([0-9]+)<\/td><td>(.*)<\/td><td class.*class.*">(.*)<\/span><\/td>`)
	GuildMemberInvitesInformationRegex = regexp.MustCompile(`<td><a.*">(.*)<\/a><\/td><td>(.*)<\/td>`)
)

func TibiaGuildsGuildImpl(guild string, BoxContentHTML string, url string) (GuildResponse, error) {
	// Creating empty vars
	var (
		MembersData                                                                                                                                                    []GuildMember
		InvitedData                                                                                                                                                    []InvitedGuildMember
		GuildGuildhallData                                                                                                                                             []Guildhall
		MembersRank, MembersTitle, MembersStatus, GuildDescription, GuildDisbandedDate, GuildDisbandedCondition, GuildHomepage, GuildWorld, GuildLogoURL, GuildFounded string
		GuildActive, GuildApplications, GuildInWar, GuildDescriptionFinished                                                                                           bool
		MembersCountOnline, MembersCountOffline, MembersCountInvited                                                                                                   int
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return GuildResponse{}, fmt.Errorf("[error] TibiaGuildsGuildImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	guildName, err := ReaderHTML.Find("h1").First().Html()
	if err != nil {
		return GuildResponse{}, fmt.Errorf("[error] TibiaGuildsGuildImpl failed at ReaderHTML.Find, err: %s", err)
	}

	if guildName == "" {
		return GuildResponse{}, validation.ErrorGuildNotFound
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPA, err := ReaderHTML.Find(".BoxContent table").Html()
	if err != nil {
		return GuildResponse{}, fmt.Errorf("[error] TibiaGuildsGuildImpl failed at InnerTableContainerTMPA ReaderHTML.Find, err: %s", err)
	}

	subma1b := GuildLogoRegex.FindAllStringSubmatch(InnerTableContainerTMPA, -1)

	if len(subma1b) > 0 {
		GuildLogoURL = subma1b[0][1]
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find("#GuildInformationContainer").Html()
	if err != nil {
		return GuildResponse{}, fmt.Errorf("[error] TibiaGuildsGuildImpl failed at InnerTableContainerTMPB ReaderHTML.Find, err: %s", err)
	}

	for line := range strings.SplitSeq(strings.TrimSuffix(InnerTableContainerTMPB, "\n"), "\n") {
		// Guild information
		if !GuildDescriptionFinished {
			// First line is the description..
			GuildDescription += strings.ReplaceAll(line+"\n", "<br/><br/>\n", "")

			// Abort loop and continue wiht next section
			if strings.Contains(line, "<br/><br/>") {
				GuildDescription = strings.TrimSpace(TibiaDataSanitizeEscapedString(GuildDescription))
				GuildDescriptionFinished = true
			}

		}

		if GuildDescriptionFinished || strings.HasPrefix(line, "The guild was founded on ") {
			// The rest of the Guild information

			if strings.HasPrefix(GuildDescription, "The guild was founded on ") {
				GuildDescription = ""
				GuildDescriptionFinished = true
			}

			if strings.HasPrefix(line, "The guild was founded on") {
				// Regex to get GuildWorld and GuildFounded
				subma1b := GuildWorldAndFoundationRegex.FindAllStringSubmatch(line, -1)
				if len(subma1b) != 0 {
					GuildWorld = subma1b[0][1]
					GuildFounded = TibiaDataDate(subma1b[0][2])
				}
			}

			// If to get GuildActive
			if strings.HasPrefix(line, "It is currently active") {
				GuildActive = true
			}

			// If open for applications
			if strings.HasPrefix(line, "Guild is opened for applications.") {
				GuildApplications = true
			} else if strings.HasPrefix(line, "Guild is closed for applications during war.") {
				GuildInWar = true
			}

			if strings.HasPrefix(line, "The official homepage is") {
				subma1c := GuildHomepageRegex.FindAllStringSubmatch(line, -1)
				GuildHomepage = subma1c[0][1]
			}

			// If guildhall
			if strings.HasPrefix(line, "Their home on "+GuildWorld) {
				subma1b := GuildhallRegex.FindAllStringSubmatch(line, -1)

				GuildGuildhallData = append(GuildGuildhallData, Guildhall{
					Name:      TibiaDataSanitizeEscapedString(subma1b[0][1]),
					World:     GuildWorld,
					PaidUntil: TibiaDataDate(subma1b[0][2]),
				})
			}

			// If disbanded
			if strings.HasPrefix(line, "<b>It will be disbanded on ") {
				subma1c := GuildDisbaneRegex.FindAllStringSubmatch(line, -1)
				if len(subma1c) > 0 {
					GuildDisbandedDate = TibiaDataDate(subma1c[0][1])
					GuildDisbandedCondition = subma1c[0][2]
				}
			}
		}
	}

	var insideError error

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into GuildsDivHTML
		GuildsDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error]  TibiaGuildsGuildImpl failed at GuildsDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Removing linebreaks from HTML
		GuildsDivHTML = TibiaDataHTMLRemoveLinebreaks(GuildsDivHTML)

		// Regex to get data for record values
		subma1 := GuildMemberInformationRegex.FindAllStringSubmatch(GuildsDivHTML, -1)

		if len(subma1) > 0 {
			// Rank name
			if TibiaDataSanitizeStrings(subma1[0][1]) != " " {
				MembersRank = subma1[0][1]
			}

			// Title
			MembersTitle = strings.ReplaceAll(strings.ReplaceAll(subma1[0][3], " (", ""), ")", "")

			// Status
			if strings.Contains(subma1[0][7], "online") {
				MembersStatus = "online"
				MembersCountOnline++
			} else {
				MembersStatus = "offline"
				MembersCountOffline++
			}

			MembersData = append(MembersData, GuildMember{
				Name:     TibiaDataSanitizeStrings(subma1[0][2]),
				Title:    MembersTitle,
				Rank:     MembersRank,
				Vocation: subma1[0][4],
				Level:    TibiaDataStringToInteger(subma1[0][5]),
				Joined:   TibiaDataDate(subma1[0][6]),
				Status:   MembersStatus,
			})
		} else {
			// Regex to get data for record values
			subma2 := GuildMemberInvitesInformationRegex.FindAllStringSubmatch(GuildsDivHTML, -1)

			if len(subma2) > 0 {
				MembersCountInvited++
				InvitedData = append(InvitedData, InvitedGuildMember{
					Name: TibiaDataSanitizeStrings(subma2[0][1]),
					Date: TibiaDataDate(subma2[0][2]),
				})
			}
		}

		return true
	})

	if insideError != nil {
		return GuildResponse{}, insideError
	}

	//
	// Build the data-blob
	return GuildResponse{
		Guild{
			Name:               guildName,
			World:              GuildWorld,
			LogoURL:            GuildLogoURL,
			Description:        GuildDescription,
			Guildhalls:         GuildGuildhallData,
			Active:             GuildActive,
			Founded:            GuildFounded,
			Applications:       GuildApplications,
			Homepage:           GuildHomepage,
			InWar:              GuildInWar,
			DisbandedDate:      GuildDisbandedDate,
			DisbandedCondition: GuildDisbandedCondition,
			PlayersOnline:      MembersCountOnline,
			PlayersOffline:     MembersCountOffline,
			MembersTotal:       (MembersCountOnline + MembersCountOffline),
			MembersInvited:     MembersCountInvited,
			Members:            MembersData,
			Invited:            InvitedData,
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
