package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Child of BoostableBoss (used for list of boostable bosses and boosted boss section)
type OverviewBoostableBoss struct {
	Name     string `json:"name"`      // The name of the boss.
	ImageURL string `json:"image_url"` // The URL to this boss's image.
	Featured bool   `json:"featured"`  // Whether it is featured of not.
}

// Child of JSONData
type BoostableBossesContainer struct {
	Boosted         OverviewBoostableBoss   `json:"boosted"`             // The current boosted boss.
	BoostableBosses []OverviewBoostableBoss `json:"boostable_boss_list"` // The list of boostable bosses.
}

// The base includes two levels: BoostableBosses and Information
type BoostableBossesOverviewResponse struct {
	BoostableBosses BoostableBossesContainer `json:"boostable_bosses"`
	Information     Information              `json:"information"`
}

func TibiaBoostableBossesOverviewImpl(BoxContentHTML string, url string) (BoostableBossesOverviewResponse, error) {
	const (
		bodyIndexer    = `<body`
		endBodyIndexer = `</body>`

		todayChecker  = `Today's boosted boss: `
		bossesChecker = `<div class="CaptionContainer" >`

		todayBossIndexer    = `title="` + todayChecker
		endTodayBossIndexer = `" src="`

		todayBossImgIndexer    = `https://static.tibia.com/images/global/header/monsters/`
		endTodayBossImgIndexer = `" onClick="`

		bossesImgIndexer    = `https://static.tibia.com/images/library/`
		endBossesImgIndexer = `"`

		bossesNameIndexer    = `border=0 />      <div>`
		endBossesNameIndexer = `</div>`
	)

	bodyIdx := strings.Index(
		BoxContentHTML, bodyIndexer,
	)

	if bodyIdx == -1 {
		return BoostableBossesOverviewResponse{}, errors.New("[error] body passd to TibiaBoostableBossesOverviewImpl is not valid")
	}

	endBodyIdx := strings.Index(
		BoxContentHTML[bodyIdx:], endBodyIndexer,
	) + bodyIdx + len(endBodyIndexer)

	if endBodyIdx == -1 {
		return BoostableBossesOverviewResponse{}, errors.New("[error] body passd to TibiaBoostableBossesOverviewImpl is not valid")
	}

	data := BoxContentHTML[bodyIdx:endBodyIdx]

	var (
		started bool

		boostedBossName string
		boostedBossImg  string

		bosses = make([]OverviewBoostableBoss, 0, validation.AmountOfBoostableBosses)
	)

	split := strings.Split(data, "\n")
	for _, cur := range split {
		isTodaysLine := strings.Contains(cur, todayChecker) && !started
		isBossesLine := strings.Contains(cur, bossesChecker)

		if !isTodaysLine && !isBossesLine {
			continue
		}

		if isTodaysLine {
			started = true

			todayBossIdx := strings.Index(
				cur, todayBossIndexer,
			) + len(todayBossIndexer)
			endTodayBossIdx := strings.Index(
				cur[todayBossIdx:], endTodayBossIndexer,
			) + todayBossIdx

			boostedBossName = TibiaDataSanitizeEscapedString(
				cur[todayBossIdx:endTodayBossIdx],
			)

			todayBossImgIdx := strings.Index(
				cur[todayBossIdx:], todayBossImgIndexer,
			) + todayBossIdx
			endTodayBossImgIdx := strings.Index(
				cur[todayBossImgIdx:], endTodayBossImgIndexer,
			) + todayBossImgIdx

			boostedBossImg = cur[todayBossImgIdx:endTodayBossImgIdx]
		}

		if isBossesLine {
			for idx := strings.Index(cur, bossesImgIndexer); idx != -1; idx = strings.Index(cur, bossesImgIndexer) {
				imgIdx := strings.Index(
					cur, bossesImgIndexer,
				)
				endImgIdx := strings.Index(
					cur[imgIdx:], endBossesImgIndexer,
				) + imgIdx
				img := cur[imgIdx:endImgIdx]

				nameIdx := strings.Index(
					cur, bossesNameIndexer,
				) + len(bossesNameIndexer)
				endNameIdx := strings.Index(
					cur[nameIdx:], endBossesNameIndexer,
				) + nameIdx
				name := TibiaDataSanitizeEscapedString(cur[nameIdx:endNameIdx])

				bosses = append(bosses, OverviewBoostableBoss{
					Name:     name,
					ImageURL: img,
					Featured: name == boostedBossName,
				})

				cur = cur[endNameIdx-1:]
			}

			break
		}
	}

	return BoostableBossesOverviewResponse{
		BoostableBossesContainer{
			Boosted: OverviewBoostableBoss{
				Name:     boostedBossName,
				ImageURL: boostedBossImg,
				Featured: true,
			},
			BoostableBosses: bosses,
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
