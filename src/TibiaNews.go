package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of JSONData
type News struct {
	ID          int    `json:"id"`              // The internal ID of the news.
	Date        string `json:"date"`            // The date when the news was published.
	Title       string `json:"title,omitempty"` // The title of the news.
	Category    string `json:"category"`        // The category of the news.
	Type        string `json:"type,omitempty"`  // The type of news.
	TibiaURL    string `json:"url"`             // The URL for the news with id.
	Content     string `json:"content"`         // The news in plain text.
	ContentHTML string `json:"content_html"`    // The news in HTML format.
}

// The base
type NewsResponse struct {
	News        News        `json:"news"`
	Information Information `json:"information"`
}

var (
	martelRegex                      = regexp.MustCompile(`<img src=\"https:\/\/static\.tibia\.com\/images\/global\/letters\/letter_martel_(.)\.gif\" ([^\/>]+..)`)
	newsTickerAnchorRelFirstRegex    = regexp.MustCompile(`<a href="([^"]+)" rel="noopener noreferrer" target="_blank">`)
	newsTickerAnchorTargetFirstRegex = regexp.MustCompile(`<a href="([^"]+)" target="_blank" rel="noopener noreferrer">`)
)

func normalizeTickerContentHTML(contentHTML string) string {
	contentHTML = newsTickerAnchorRelFirstRegex.ReplaceAllString(contentHTML, `<a href="$1" target="_blank" rel="noopener noreferrer" rel="noopener">`)
	contentHTML = newsTickerAnchorTargetFirstRegex.ReplaceAllString(contentHTML, `<a href="$1" target="_blank" rel="noopener noreferrer" rel="noopener">`)
	return contentHTML
}

func TibiaNewsImpl(NewsID int, rawUrl string, BoxContentHTML string) (NewsResponse, error) {
	// Declaring vars for later use..
	var (
		NewsData    News
		tmp1        *goquery.Selection
		tmp2        string
		insideError error
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return NewsResponse{}, fmt.Errorf("[error] TibiaNewsImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	NewsData.ID = NewsID
	NewsData.TibiaURL = rawUrl

	ReaderHTML.Find(".NewsHeadline").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// getting category by image src
		CategoryImg, _ := s.Find("img").Attr("src")
		NewsData.Category = TibiaDataGetNewsCategory(CategoryImg)
		// getting date from headline
		tmp1 = s.Find(".NewsHeadlineDate")
		tmp2, err = tmp1.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaNewsImpl failed at tmp2, err = tmp1.Html(), NewsHeadlineDate, err: %s", err)
			return false
		}

		NewsData.Date = TibiaDataDate(strings.ReplaceAll(tmp2, " - ", ""))

		// getting headline text (which could be title or also type)
		tmp1 = s.Find(".NewsHeadlineText")
		tmp2, err = tmp1.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaNewsImpl failed at tmp2, err = tmp1.Html(), NewsHeadlineText, err: %s", err)
			return false
		}

		NewsData.Title = RemoveHtmlTag(tmp2)
		if NewsData.Title == "News Ticker" {
			NewsData.Type = "ticker"
			NewsData.Title = ""
		}

		return true
	})

	if insideError != nil {
		return NewsResponse{}, insideError
	}

	ReaderHTML.Find(".NewsTableContainer").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// checking if its a ticker..
		if NewsData.Type == "ticker" {
			tmp1 = s.Find("p")
			if tmp1.Text() == "" {
				NewsData.Content = s.Text()
				NewsData.ContentHTML, err = s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaNewsImpl failed at NewsData.ContentHTML, err = s.Html(), err: %s", err)
					return false
				}
				NewsData.ContentHTML = normalizeTickerContentHTML(NewsData.ContentHTML)
			} else {
				NewsData.Content = tmp1.Text()
				NewsData.ContentHTML, err = tmp1.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaNewsImpl failed at NewsData.ContentHTML, err = tmp1.Html(), err: %s", err)
					return false
				}
				NewsData.ContentHTML = normalizeTickerContentHTML(NewsData.ContentHTML)
			}
		} else {
			// getting html
			tmp2, err = s.First().Html()
			if err != nil {
				insideError = fmt.Errorf("[error] TibiaNewsImpl failed at NewsData.ContentHTML, tmp2, err = s.First().Html(), err: %s", err)
				return false
			}

			// replacing martel letter in articles with real letters
			tmp2 = martelRegex.ReplaceAllString(tmp2, "$1")
			s.ReplaceWithHtml(tmp2)

			// storing html content
			NewsData.ContentHTML = tmp2

			// reading string again after replacing letters
			tmp1, err := goquery.NewDocumentFromReader(strings.NewReader(tmp2))
			if err != nil {
				insideError = fmt.Errorf("[error] TibiaNewsImpl failed attmp1, err := goquery.NewDocumentFromReader, err: %s", err)
				return false
			}

			// storing text content
			NewsData.Content = tmp1.Text()
		}

		return true
	})

	if insideError != nil {
		return NewsResponse{}, insideError
	}

	//
	// Build the data-blob
	return NewsResponse{
		NewsData,
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{rawUrl},
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
