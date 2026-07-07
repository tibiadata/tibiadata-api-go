package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestNewsById(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/news/archive/6529.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	newsArticleJson, err := TibiaNewsImpl(6529, "https://www.tibia.com/news/?subtopic=newsarchive&id=6529", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := newsArticleJson.Information

	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6529", information.TibiaURLs[0])

	assert.Equal(6529, newsArticleJson.News.ID)
	assert.Equal("2022-01-12", newsArticleJson.News.Date)
	assert.Empty(newsArticleJson.News.Title)
	assert.Equal("development", newsArticleJson.News.Category)
	assert.Equal("ticker", newsArticleJson.News.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6529", newsArticleJson.News.TibiaURL)
	assert.Equal("A number of issues related to the 25 years activities have been fixed, including the following: Dragon pinatas that were stuck in the inbox have been changed into dragon pinata kits. The wind-up loco can now be taken even after defeating Lord Retro. The baby seals can now be painted even when the quest The Ice Islands is active. The weight of wallpapers and fairy lights has been increased and their market category has been changed to decoration. A number of typos and map issues have been fixed as well.", newsArticleJson.News.Content)
	assert.Equal("A number of issues related to the 25 years activities have been fixed, including the following: Dragon pinatas that were stuck in the inbox have been changed into dragon pinata kits. The wind-up loco can now be taken even after defeating Lord Retro. The baby seals can now be painted even when the quest The Ice Islands is active. The weight of wallpapers and fairy lights has been increased and their market category has been changed to decoration. A number of typos and map issues have been fixed as well.", newsArticleJson.News.ContentHTML)
}

func TestNews6512(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/news/archive/6512.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	newsArticleJson, err := TibiaNewsImpl(6512, "https://www.tibia.com/news/?subtopic=newsarchive&id=6512", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(6512, newsArticleJson.News.ID)
	assert.Equal("2022-01-04", newsArticleJson.News.Date)
	assert.Empty(newsArticleJson.News.Title)
	assert.Equal("community", newsArticleJson.News.Category)
	assert.Equal("ticker", newsArticleJson.News.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6512", newsArticleJson.News.TibiaURL)
	assert.Equal("TibiaData.com has some news to share! First of all, they invite you to participate in their Discord Server. Further, they are now present on GitHub. They are working on their v3, which is still in beta. If you are interested in such things, head on over there to see what is cooking.", newsArticleJson.News.Content)
	assert.Equal("<a href=\"https://tibiadata.com\" target=\"_blank\" rel=\"noopener noreferrer\" rel=\"noopener\">TibiaData.com</a> has some news to share! First of all, they invite you to participate in their <a href=\"https://tibiadata.com/2021/07/join-tibiadata-on-discord/\" target=\"_blank\" rel=\"noopener noreferrer\" rel=\"noopener\">Discord Server</a>. Further, they are now present on <a href=\"https://tibiadata.com/2021/12/tibiadata-has-gone-open-source/\" target=\"_blank\" rel=\"noopener noreferrer\" rel=\"noopener\">GitHub</a>. They are working on their <a href=\"https://tibiadata.com/doc-api-v3/v3-beta/\" target=\"_blank\" rel=\"noopener noreferrer\" rel=\"noopener\">v3</a>, which is still in beta. If you are interested in such things, head on over there to see what is cooking.", newsArticleJson.News.ContentHTML)
}

func TestNews504(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/news/archive/504.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	newsArticleJson, err := TibiaNewsImpl(504, "https://www.tibia.com/news/?subtopic=newsarchive&id=504", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(504, newsArticleJson.News.ID)
	assert.Equal("2007-04-27", newsArticleJson.News.Date)
	assert.Empty(newsArticleJson.News.Title)
	assert.Equal("community", newsArticleJson.News.Category)
	assert.Equal("ticker", newsArticleJson.News.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=504", newsArticleJson.News.TibiaURL)
	assert.Equal("A new feedback form has been released today. Help us to find out which websites and magazines are popular in your country by filling out the new questionnaire. In our current poll we are curious about your occupation.", newsArticleJson.News.Content)
	assert.Equal("<br/>A new feedback form has been released today. Help us to find out which websites and magazines are popular in your country by filling out the new questionnaire. In our current poll we are curious about your occupation.", newsArticleJson.News.ContentHTML)
}

func TestNews6481(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/news/archive/6481.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	newsArticleJson, err := TibiaNewsImpl(6481, "https://www.tibia.com/news/?subtopic=newsarchive&id=6481", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(6481, newsArticleJson.News.ID)
	assert.Equal("2021-12-22", newsArticleJson.News.Date)
	assert.Equal("New Mounts", newsArticleJson.News.Title)
	assert.Equal("development", newsArticleJson.News.Category)
	assert.Empty(newsArticleJson.News.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6481", newsArticleJson.News.TibiaURL)
	assert.Equal("New mounts have been added to the Store today!\nThe origins of the Emerald Raven, Mystic Raven, and Radiant Raven are shrouded in darkness, as no written record nor tale told by even the most knowing storytellers mentions but a trace of them. Superstition surrounds them, as some see these gigantic birds as an echo of a long forgotten past, while others believe them to herald hitherto unknown events. What is clear is that they are highly intelligent beings which make great companions if they deem somebody worthy.\n\nClick on image to enlarge.\n\nOnce bought, your character can use the mount ingame anytime, no matter if you are a free account or a Premium account.\nGet yourself a corvid companion!Your Community Managers ", newsArticleJson.News.Content)
	assert.Equal("<p>New mounts have been added to the Store today!</p>\n<p>The origins of the <strong>Emerald Raven</strong>, <strong>Mystic Raven</strong>, and <strong>Radiant Raven</strong> are shrouded in darkness, as no written record nor tale told by even the most knowing storytellers mentions but a trace of them. Superstition surrounds them, as some see these gigantic birds as an echo of a long forgotten past, while others believe them to herald hitherto unknown events. What is clear is that they are highly intelligent beings which make great companions if they deem somebody worthy.</p>\n<figure><center><img style=\"cursor: pointer;\" src=\"https://static.tibia.com/images/news/emeraldraven_small.jpg\" onclick=\"ImageInNewWindow(&#39;https://static.tibia.com/images/news/emeraldraven.jpg&#39;)\"/></center>\n<figcaption><center><em>Click on image to enlarge.</em></center></figcaption>\n</figure>\n<p>Once bought, your character can use the mount ingame anytime, no matter if you are a free account or a Premium account.</p>\n<p>Get yourself a corvid companion!<br/>Your Community Managers</p> ", newsArticleJson.News.ContentHTML)
}
