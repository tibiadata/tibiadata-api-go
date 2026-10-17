package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of JSONData
type NewsItem struct {
	ID       int    `json:"id"`                // The internal ID of the news.
	Date     string `json:"date"`              // The date when the news was published.
	News     string `json:"news"`              // The news in plain text.
	Category string `json:"category"`          // The category of the news.
	Type     string `json:"type"`              // The type of news.
	TibiaURL string `json:"url"`               // The URL for the news with id.
	ApiURL   string `json:"url_api,omitempty"` // The URL for the news in this API.
}

// The base
type NewsListResponse struct {
	News        []NewsItem  `json:"news"`
	Information Information `json:"information"`
}

func TibiaNewslistImpl(days int, BoxContentHTML, handlerURL, apiVersion string) (NewsListResponse, error) {
	// Declaring vars for later use..
	var NewsListData []NewsItem

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return NewsListResponse{}, fmt.Errorf("[error] TibiaNewslistImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var insideError error

	ReaderHTML.Find(".Odd,.Even").EachWithBreak(func(index int, s *goquery.Selection) bool {
		var OneNews NewsItem

		// getting category by image src
		CategoryImg, _ := s.Find("img").Attr("src")
		OneNews.Category = TibiaDataGetNewsCategory(CategoryImg)

		// getting type from headline
		NewsType := s.Nodes[0].FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild.Data
		OneNews.Type = TibiaDataGetNewsType(TibiaDataSanitizeStrings(NewsType))

		// getting date from headline
		OneNews.Date = TibiaDataDate(s.Nodes[0].FirstChild.NextSibling.FirstChild.Data)
		OneNews.News = s.Find("a").Text()

		// getting remaining things as URLs
		NewsURL, _ := s.Find("a").Attr("href")
		p, err := url.Parse(NewsURL)
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaNewslistImpl failed at p, err := url.Parse(NewsURL), err: %s", err)
			return false
		}
		NewsID := p.Query().Get("id")
		NewsSplit := strings.Split(NewsURL, NewsID)
		OneNews.ID = TibiaDataStringToInteger(NewsID)
		OneNews.TibiaURL = NewsSplit[0] + NewsID

		if TibiaDataHost != "" {
			OneNews.ApiURL = TibiaDataProtocol + "://" + TibiaDataHost + "/" + apiVersion + "/news/id/" + NewsID
		}

		// add to NewsListData for response
		NewsListData = append(NewsListData, OneNews)

		return true
	})

	if insideError != nil {
		return NewsListResponse{}, insideError
	}

	//
	// Build the data-blob
	return NewsListResponse{
		NewsListData,
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{handlerURL},
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
