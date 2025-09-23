package main

import (
	"html"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// TibiaDataDatetime func
func TibiaDataDatetime(date string) string {
	//TODO: Normalization needs to happen above this layer
	date = norm.NFKC.String(date)

	var (
		returnDate time.Time
		err        error
	)

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now()
	} else {
		// timezone use in html: CET/CEST
		loc, _ := time.LoadLocation("Europe/Berlin")

		// format used in datetime on html: Jan 02 2007, 19:20:30 CET
		formatting := "Jan 02 2006, 15:04:05 MST"

		// parsing html in time with location set in loc
		returnDate, err = time.ParseInLocation(formatting, date, loc)

		// parsing html in tiem without loc
		//returnDate, err = time.Parse("Jan 02 2006, 15:04:05 MST", date)

		if err != nil {
			log.Println(err)
		}
	}

	// Return of formatted date and time string to functions..
	return returnDate.UTC().Format(time.RFC3339)
}

// TibiaDataHTMLRemoveLinebreaks func
func TibiaDataHTMLRemoveLinebreaks(data string) string {
	return strings.ReplaceAll(data, "\n", "")
}

var removeUrlRegex = regexp.MustCompile(`<a.*>(.*)<\/a>`)

// TibiaDataRemoveURLs func
func TibiaDataRemoveURLs(data string) string {
	// prepare return value
	var returnData string

	result := removeUrlRegex.FindAllStringSubmatch(data, -1)

	if len(result) > 0 {
		returnData = result[0][1]
	} else {
		returnData = ""
	}

	return returnData
}

// TibiaDataStringWorldFormatToTitle func
func TibiaDataStringWorldFormatToTitle(world string) string {
	return cases.Title(language.English).String(world)
}

// TibiaDataQueryEscapeString func - encode string to be correct formatted
func TibiaDataQueryEscapeString(data string) string {
	// switching "+" to " "
	data = strings.ReplaceAll(data, "+", " ")

	// returning with QueryEscape function
	return url.QueryEscape(data)
}

// TibiaDataDate func
func TibiaDataDate(date string) string {
	// removing weird spacing and comma
	date = TibiaDataSanitizeStrings(strings.ReplaceAll(date, ",", ""))

	// var time parser
	var tmpDate time.Time
	var err error

	// date formats to parse
	dateFormats := map[string][]string{
		"YearMonthDay": {"January 2 2006", "Jan 02 2006"},
		"YearMonth":    {"January 2006", "Jan 2006", "2006-01", "01/06"},
	}

	for _, layout := range dateFormats["YearMonthDay"] {
		tmpDate, err = time.Parse(layout, date)
		if err == nil {
			// If the parse succeeds, format the date as "2006-01-02"
			return tmpDate.UTC().Format("2006-01-02")
		}
	}

	for _, layout := range dateFormats["YearMonth"] {
		tmpDate, err = time.Parse(layout, date)
		if err == nil {
			// If the parse succeeds, format the date as "2006-01"
			return tmpDate.Format("2006-01")
		}
	}

	return tmpDate.UTC().Format("2006-01-02")
}

// TibiaDataStringToInteger converts a string to an int
func TibiaDataStringToInteger(data string) int {
	str := strings.ReplaceAll(data, ",", "")

	returnData, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("[warning] TibiaDataStringToInteger: couldn't convert string into int. error: %s", err)
	}

	return returnData
}

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

// RemoveHtmlTag replaces all HTML tags with an empty string
//
// Algo from:
// https://stackoverflow.com/questions/55036156/how-to-replace-all-html-tag-with-empty-string-in-golang
func RemoveHtmlTag(s string) string {
	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}

// TibiaDataSanitizeEscapedString func - run unescape string on string
func TibiaDataSanitizeEscapedString(data string) string {
	return html.UnescapeString(data)
}

// TibiaDataSanitizeDoubleQuoteString func - replaces double quotes to single quotes in strings
func TibiaDataSanitizeDoubleQuoteString(data string) string {
	return strings.ReplaceAll(data, "\"", "'")
}

// TibiaDataSanitizeStrings func - replacing various encoded strings to pure html
func TibiaDataSanitizeStrings(data string) string {
	// replaces weird \u00A0 string to real space
	data = strings.ReplaceAll(data, "\u00A0", " ")
	// replaces weird \u0026 string to amp (&)
	data = strings.ReplaceAll(data, "\u0026", "&")
	// returning string unescaped
	return TibiaDataSanitizeEscapedString(data)
}

// TibiaDataSanitize0026String replaces \u0026#39; with '
func TibiaDataSanitize0026String(data string) string {
	return strings.ReplaceAll(data, "\u0026#39;", "'")
}

// isEnvExist func - check if environment var is set and not empty
func isEnvExist(key string) bool {
	data, ok := os.LookupEnv(key)
	return ok && data != ""
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultVal
}

// getEnvAsBool func - read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// TibiaDataConvertValuesWithK func - convert price strings that contain k, kk or more to 3x0
func TibiaDataConvertValuesWithK(data string) int {
	return TibiaDataStringToInteger(strings.ReplaceAll(data, "k", "") + strings.Repeat("000", strings.Count(data, "k")))
}

// TibiaDataVocationValidator func - return valid vocation string and vocation id
func TibiaDataVocationValidator(vocation string) (string, string) {
	// defining return vars
	var vocationid string

	switch strings.ToLower(vocation) {
	case "none":
		vocationid = "1"
		vocation = "none"
	case "knight", "knights":
		vocationid = "2"
		vocation = "knights"
	case "paladin", "paladins":
		vocationid = "3"
		vocation = "paladins"
	case "sorcerer", "sorcerers":
		vocationid = "4"
		vocation = "sorcerers"
	case "druid", "druids":
		vocationid = "5"
		vocation = "druids"
	case "monk", "monks":
		vocationid = "6"
		vocation = "monks"
	default:
		vocationid = "0"
		vocation = "all"
	}

	// returning vars
	return vocation, vocationid
}

// TibiaDataGetNewsCategory func - extract news category by newsicon
func TibiaDataGetNewsCategory(data string) string {
	switch {
	case strings.Contains(data, "newsicon_cipsoft"):
		return "cipsoft"
	case strings.Contains(data, "newsicon_community"):
		return "community"
	case strings.Contains(data, "newsicon_development"):
		return "development"
	case strings.Contains(data, "newsicon_support"):
		return "support"
	case strings.Contains(data, "newsicon_technical"):
		return "technical"
	default:
		return "unknown"
	}
}

// TibiaDataGetNewsType func - extract news type
func TibiaDataGetNewsType(data string) string {
	switch data {
	case "News Ticker":
		return "ticker"
	case "Featured Article":
		return "article"
	case "News":
		return "news"
	default:
		return "unknown"
	}
}
