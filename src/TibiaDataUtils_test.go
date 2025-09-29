package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTibiaCETDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T08:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 CET"))
}

func TestTibiaCESTDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T07:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 CEST"))
}

func TestTibiaUTCDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T09:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 UTC"))
}

func TestIsEnvExist(t *testing.T) {
	assert := assert.New(t)

	// Test when environment variable exists and is not empty
	os.Setenv("TIBIADATA_ENV", "production")
	assert.True(isEnvExist("TIBIADATA_ENV"))

	// Test when environment variable exists and is empty
	os.Setenv("TIBIADATA_ENV", "")
	assert.False(isEnvExist("TIBIADATA_ENV"))

	// Test when environment variable does not exist
	os.Unsetenv("TIBIADATA_ENV")
	assert.False(isEnvExist("TIBIADATA_ENV"))
}

func TestGetEnv(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("TIBIADATA_ENV", "production")
	defer os.Unsetenv("TIBIADATA_ENV")

	// Test when environment variable is set
	assert.Equal("production", getEnv("TIBIADATA_ENV", "default"))

	// Test when environment variable is not set
	assert.Equal("default", getEnv("NON_EXISTENT_ENV", "default"))
}

func TestGetEnvAsBool(t *testing.T) {
	assert := assert.New(t)

	// Test when environment variable is not set
	assert.True(true, getEnvAsBool("TIBIADATA_ENV", true))
	assert.False(false, getEnvAsBool("TIBIADATA_ENV", false))

	// Test when environment variable is set to true
	os.Setenv("TIBIADATA_ENV", "true")
	assert.True(true, getEnvAsBool("TIBIADATA_ENV", false))

	// Test when environment variable is set to false
	os.Setenv("TIBIADATA_ENV", "false")
	assert.False(false, getEnvAsBool("TIBIADATA_ENV", true))

	os.Unsetenv("TIBIADATA_ENV")
}

func TestTibiaDataVocationValidator(t *testing.T) {
	assert := assert.New(t)

	var (
		x, y string
	)

	x, y = TibiaDataVocationValidator("none")
	assert.Equal(x, "none")
	assert.Equal(y, "1")
	x, y = TibiaDataVocationValidator("knight")
	assert.Equal(x, "knights")
	assert.Equal(y, "2")
	x, y = TibiaDataVocationValidator("paladin")
	assert.Equal(x, "paladins")
	assert.Equal(y, "3")
	x, y = TibiaDataVocationValidator("sorcerer")
	assert.Equal(x, "sorcerers")
	assert.Equal(y, "4")
	x, y = TibiaDataVocationValidator("druid")
	assert.Equal(x, "druids")
	assert.Equal(y, "5")
	x, y = TibiaDataVocationValidator("monk")
	assert.Equal(x, "monks")
	assert.Equal(y, "6")
	x, y = TibiaDataVocationValidator("")
	assert.Equal(x, "all")
	assert.Equal(y, "0")
}

func TestTibiaDataGetNewsCategory(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("cipsoft", TibiaDataGetNewsCategory("newsicon_cipsoft"))
	assert.Equal("community", TibiaDataGetNewsCategory("newsicon_community"))
	assert.Equal("development", TibiaDataGetNewsCategory("newsicon_development"))
	assert.Equal("support", TibiaDataGetNewsCategory("newsicon_support"))
	assert.Equal("technical", TibiaDataGetNewsCategory("newsicon_technical"))
	assert.Equal("unknown", TibiaDataGetNewsCategory("newsicon_tibiadata"))
}

func TestTibiaDataGetNewsType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("ticker", TibiaDataGetNewsType("News Ticker"))
	assert.Equal("article", TibiaDataGetNewsType("Featured Article"))
	assert.Equal("news", TibiaDataGetNewsType("News"))
	assert.Equal("unknown", TibiaDataGetNewsType("TibiaData"))
}

func TestHTMLLineBreakRemover(t *testing.T) {
	const str = "a\nb\nc\nd\ne\nf\ng\nh\n"

	sanitizedStr := TibiaDataHTMLRemoveLinebreaks(str)
	assert.Equal(t, sanitizedStr, "abcdefgh")
}

func TestURLsRemover(t *testing.T) {
	const (
		strOne = `<a href="https://www.tibia.com/community/?subtopic=characters&amp;name=Bobeek">Bobeek</a>`
		strTwo = `<div>Bobeek</div>`
	)

	// Test when input contains a URL
	sanitizedStr := TibiaDataRemoveURLs(strOne)
	assert.Equal(t, sanitizedStr, "Bobeek")

	// Test when input does not contain a URL
	sanitizedStr = TibiaDataRemoveURLs(strTwo)
	assert.Equal(t, sanitizedStr, "")
}

func TestWorldFormater(t *testing.T) {
	const str = "hEsThDIáÛõ"

	sanitizedStr := TibiaDataStringWorldFormatToTitle(str)
	assert.Equal(t, sanitizedStr, "Hesthdiáûõ")
}

func TestEscaper(t *testing.T) {
	const (
		strOne   = "god durin"
		strTwo   = "god+durin"
		strThree = "gód"
		strFour  = "Näurin"
	)

	sanitizedStrOne := TibiaDataQueryEscapeString(strOne)
	sanitizedStrTwo := TibiaDataQueryEscapeString(strTwo)
	sanitizedStrThree := TibiaDataQueryEscapeString(strThree)
	sanitizedStrFour := TibiaDataQueryEscapeString(strFour)

	assert := assert.New(t)
	assert.Equal(sanitizedStrOne, "god+durin")
	assert.Equal(sanitizedStrTwo, "god+durin")
	assert.Equal(sanitizedStrThree, "g%C3%B3d")
	assert.Equal(sanitizedStrFour, "N%C3%A4urin")
}

func TestDateParser(t *testing.T) {
	const (
		strYearMonthDay = "2022-03-09"
		strYearMonth    = "2022-03"
	)

	assert := assert.New(t)
	assert.Equal("0001-01-01", TibiaDataDate(""))

	// YearMonthDay
	assert.Equal(strYearMonthDay, TibiaDataDate("March 9 2022"))
	assert.Equal(strYearMonthDay, TibiaDataDate("Mar 09 2022"))

	// YearMonth
	assert.Equal(strYearMonth, TibiaDataDate("March 2022"))
	assert.Equal(strYearMonth, TibiaDataDate("Mar 2022"))
	assert.Equal(strYearMonth, TibiaDataDate("2022-03"))
	assert.Equal(strYearMonth, TibiaDataDate("03/22"))
}

func TestStringToInt(t *testing.T) {
	const (
		isInt = 123
		noInt = 0
	)

	assert := assert.New(t)

	// Test when input is a valid integer
	assert.Equal(isInt, TibiaDataStringToInteger("123"))

	// Test when input contains commas
	assert.Equal(isInt, TibiaDataStringToInteger("1,2,3"))

	// Test when input is not a valid integer
	assert.Equal(noInt, TibiaDataStringToInteger("not an integer"))
}

func TestHTMLRemover(t *testing.T) {
	const str = `<div id="DeactivationContainer" onclick="ActivateWebsiteFrame();$('.LightBoxContentToHide').css('display', 'none');">abc</div>`

	sanitizedString := RemoveHtmlTag(str)
	assert.Equal(t, sanitizedString, "abc")
}

func TestFake(t *testing.T) {
	assert := assert.New(t)
	const str = "&lt;"

	htmlString := TibiaDataSanitizeEscapedString(str)
	assert.Equal(htmlString, "<")

	const strTwo = `"`

	doubleString := TibiaDataSanitizeDoubleQuoteString(strTwo)
	assert.Equal(doubleString, "'")

	const strThree = "\u00A0"

	nbspString := TibiaDataSanitizeStrings(strThree)
	assert.Equal(nbspString, " ")

	const strFour = "\u0026#39;"

	string0026 := TibiaDataSanitize0026String(strFour)
	assert.Equal(string0026, "'")

	const strFive = "1kk"

	kkInt := TibiaDataConvertValuesWithK(strFive)
	assert.Equal(kkInt, 1000000)
}
