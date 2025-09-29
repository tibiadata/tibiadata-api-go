package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Debug stores some debug informations
type Debug struct {
	TibiaDataUserAgent                  string `json:"tibia_data_user_agent"`
	DataSha256Sum                       string `json:"data_sha_256_sum"`
	DataSha512Sum                       string `json:"data_sha_512_sum"`
	SmallestCreatureName                string `json:"smallest_creature_name"`
	BiggestCreatureName                 string `json:"biggest_creature_name"`
	SmallestCreatureWord                string `json:"smallest_creature_word"`
	BiggestCreatureWord                 string `json:"biggest_creature_word"`
	SmallestCreatureNameRuneCount       int    `json:"smallest_creature_name_rune_count"`
	BiggestCreatureNameRuneCount        int    `json:"biggest_creature_name_rune_count"`
	SmallestCreatureWordRuneCount       int    `json:"smallest_creature_word_rune_count"`
	BiggestCreatureWordRuneCount        int    `json:"biggest_creature_word_rune_count"`
	SmallestSpellNameOrFormula          string `json:"smallest_spell_name_or_formula"`
	BiggestSpellNameOrFormula           string `json:"biggest_spell_name_or_formula"`
	SmallestSpellWord                   string `json:"smallest_spell_word"`
	BiggestSpellWord                    string `json:"biggest_spell_word"`
	SmallestSpellNameOrFormulaRuneCount int    `json:"smallest_spell_name_or_formula_rune_count"`
	BiggestSpellNameOrFormulaRuneCount  int    `json:"biggest_spell_name_or_formula_rune_count"`
	SmallestSpellWordRuneCount          int    `json:"smallest_spell_word_rune_count"`
	BiggestSpellWordRuneCount           int    `json:"biggest_spell_word_rune_count"`
}

// TibiaDataRequestTraceLogger func - prints out trace information to log
func TibiaDataRequestTraceLogger(res *resty.Response, err error) {
	log.Println("TRACE RESTY",
		"\n~~~ TRACE INFO ~~~",
		"\nDNSLookup      :", res.Request.TraceInfo().DNSLookup,
		"\nConnTime       :", res.Request.TraceInfo().ConnTime,
		"\nTCPConnTime    :", res.Request.TraceInfo().TCPConnTime,
		"\nTLSHandshake   :", res.Request.TraceInfo().TLSHandshake,
		"\nServerTime     :", res.Request.TraceInfo().ServerTime,
		"\nResponseTime   :", res.Request.TraceInfo().ResponseTime,
		"\nTotalTime      :", res.Request.TraceInfo().TotalTime,
		"\nIsConnReused   :", res.Request.TraceInfo().IsConnReused,
		"\nIsConnWasIdle  :", res.Request.TraceInfo().IsConnWasIdle,
		"\nConnIdleTime   :", res.Request.TraceInfo().ConnIdleTime,
		"\nRequestAttempt :", res.Request.TraceInfo().RequestAttempt,
		"\nRemoteAddr     :", res.Request.TraceInfo().RemoteAddr.String(),
		"\nError          :", err,
		"\n==============================================================================")
}

// debugHandler returns some debug information
func debugHandler(c *gin.Context) {
	data := Information{
		APIDetails: TibiaDataAPIDetails,
		Timestamp:  TibiaDataDatetime(""),
		Status: Status{
			HTTPCode: http.StatusOK,
			Message:  "UP",
		},
	}

	debug := Debug{
		TibiaDataUserAgent: TibiaDataUserAgent,
	}

	// Shas
	sha256, err := validation.GetSha256Sum()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.DataSha256Sum = sha256

	sha512, err := validation.GetSha512Sum()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.DataSha512Sum = sha512

	// Creatures
	smallestCreatureName, err := validation.GetSmallestCreatureName()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestCreatureName = smallestCreatureName

	biggestCreatureName, err := validation.GetBiggestCreatureName()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestCreatureName = biggestCreatureName

	biggestCreatureWord, err := validation.GetBiggestCreatureWord()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestCreatureWord = biggestCreatureWord

	smallestCreatureWord, err := validation.GetSmallestCreatureWord()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestCreatureWord = smallestCreatureWord

	smallestCreatureNameRuneCount, err := validation.GetSmallestCreatureNameRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestCreatureNameRuneCount = smallestCreatureNameRuneCount

	biggestCreatureNameRuneCount, err := validation.GetBiggestCreatureNameRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestCreatureNameRuneCount = biggestCreatureNameRuneCount

	smallestCreatureWordRuneCount, err := validation.GetSmallestCreatureWordRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestCreatureWordRuneCount = smallestCreatureWordRuneCount

	biggestCreatureWordRuneCount, err := validation.GetBiggestCreatureWordRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestCreatureWordRuneCount = biggestCreatureWordRuneCount

	// Spells
	smallestSpellName, err := validation.GetSmallestSpellNameOrFormula()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestSpellNameOrFormula = smallestSpellName

	biggestSpellName, err := validation.GetBiggestSpellNameOrFormula()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestSpellNameOrFormula = biggestSpellName

	biggestSpellWord, err := validation.GetBiggestSpellWord()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestSpellWord = biggestSpellWord

	smallestSpellWord, err := validation.GetSmallestSpellWord()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestSpellWord = smallestSpellWord

	smallestSpellNameRuneCount, err := validation.GetSmallestSpellNameOrFormulaRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestSpellNameOrFormulaRuneCount = smallestSpellNameRuneCount

	biggestSpellNameRuneCount, err := validation.GetBiggestSpellNameOrFormulaRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestSpellNameOrFormulaRuneCount = biggestSpellNameRuneCount

	smallestSpellWordRuneCount, err := validation.GetSmallestSpellWordRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.SmallestSpellWordRuneCount = smallestSpellWordRuneCount

	biggestSpellWordRuneCount, err := validation.GetBiggestSpellWordRuneCount()
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusInternalServerError)
	}
	debug.BiggestSpellWordRuneCount = biggestSpellWordRuneCount

	var output DebugOutInformation
	output.Information = data
	output.Debug = debug

	c.JSON(http.StatusOK, output)
}
