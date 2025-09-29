package validation

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/tibiadata/tibiadata-api-go/src/tibiamapping"
)

// validator is a local struct representing some validation data
type validator struct {
	Worlds    []string   `json:"worlds"`
	Towns     []string   `json:"towns"`
	Houses    []House    `json:"houses"`
	Creatures []Creature `json:"creatures"`
	Spells    []Spell    `json:"spells"`
}

type Creature struct {
	Endpoint   string `json:"endpoint"`
	PluralName string `json:"plural_name"`
	Name       string `json:"name"`
}

type Spell struct {
	Name     string `json:"name"`
	Formula  string `json:"formula"`
	Endpoint string `json:"endpoint"`
}

// Houses represents a house
type House struct {
	ID   int    `json:"house_id"`
	Town string `json:"town"`
	Type string `json:"type"`
}

var (
	initiated bool             // initiated reports whether the validator has already been initiated
	val       = validator{}    // val is the local validator that will be read from to get the necessary data
	locker    = sync.RWMutex{} // locker is a locker to prevent InitiateValidator to be accessed concurrently
	sha256sum string           // sha256sum stores the sha256sum of the data.min.json file
	sha512sum string           // sha512sum stores the sha512sum of the data.min.json file

	smallestCreatureName, biggestCreatureName, smallestCreatureWord, biggestCreatureWord                                           string // smallest and biggest creature names and words
	smallestCreatureNameRuneCount, biggestCreatureNameRuneCount, smallestCreatureWordRuneCount, biggestCreatureWordRuneCount       int    // smallest and biggest creature names and words rune count
	smallestSpellNameOrFormula, biggestSpellNameOrFormula, smallestSpellWord, biggestSpellWord                                     string // smalles and biggest spell names or formulas and words
	smallestSpellNameOrFormulaRuneCount, biggestSpellNameOrFormulaRuneCount, smallestSpellWordRuneCount, biggestSpellWordRuneCount int    // smallest and biggest creature names or formulas and words rune count
)

// Initiate initiates the validator, this should be called on the init() func
func Initiate(TibiaDataUserAgent string) error {
	// Make sure InitiateValidator can not be called concurrently
	locker.Lock()
	defer locker.Unlock()

	// Check if the validator has already been initiated
	// as there is no need to initiate it twice
	if initiated {
		return ErrorAlreadyRunning
	}

	// Get the assets
	tibiaMapping, err := tibiamapping.Run(TibiaDataUserAgent)
	if err != nil {
		panic(err)
	}

	// Check if we got a nil struct
	if len(tibiaMapping.RawData) == 0 &&
		tibiaMapping.Sha256Sum == "" &&
		tibiaMapping.Sha512Sum == "" {
		return errors.New("tibia mapping struct is nil")
	}

	bytes := tibiaMapping.RawData

	sha256Fields := strings.Fields(tibiaMapping.Sha256Sum)
	sha256sum = sha256Fields[2]

	sha512Fields := strings.Fields(tibiaMapping.Sha512Sum)
	sha512sum = sha512Fields[2]

	// Check if the file is empty
	if len(bytes) == 0 {
		return errors.New("data.json file is empty")
	}

	// Unmarshal the json bytes into a go struct
	err = json.Unmarshal(bytes, &val)
	if err != nil {
		return err
	}

	// Set non changing vars
	setVars()

	// The validator is properly initiated
	initiated = true

	return nil
}

func setVars() {
	setCreaturesVars()
	setSpellsVars()
}

// setCreaturesVars sets creatures vars
// this only needs to be called once as it will never change during runtime
func setCreaturesVars() {
	if smallestCreatureName == "" {
		smallestName := val.Creatures[0].Name
		var smallestWord string

		for _, creature := range val.Creatures {
			if utf8.RuneCountInString(creature.Name) < utf8.RuneCountInString(smallestName) {
				smallestName = creature.Name
			}

			fields := strings.Fields(creature.Name)

			if utf8.RuneCountInString(smallestWord) == 0 {
				smallestWord = fields[0]
			}

			for _, str := range fields {
				if utf8.RuneCountInString(str) < utf8.RuneCountInString(smallestWord) {
					smallestWord = str
				}
			}

			if utf8.RuneCountInString(creature.Endpoint) < utf8.RuneCountInString(smallestWord) {
				smallestWord = creature.Endpoint
			}
		}

		smallestCreatureName = smallestName
		smallestCreatureNameRuneCount = utf8.RuneCountInString(smallestName)
		smallestCreatureWord = smallestWord
		smallestCreatureWordRuneCount = utf8.RuneCountInString(smallestWord)
	}

	if biggestCreatureName == "" {
		biggestName := val.Creatures[0].PluralName
		var biggestWord string

		for _, creature := range val.Creatures {
			if utf8.RuneCountInString(creature.Name) > utf8.RuneCountInString(biggestName) {
				biggestName = creature.PluralName
			}

			fields := strings.Fields(creature.Name)

			if utf8.RuneCountInString(biggestWord) == 0 {
				biggestWord = fields[0]
			}

			for _, str := range fields {
				if utf8.RuneCountInString(str) > utf8.RuneCountInString(biggestWord) {
					biggestWord = str
				}
			}

			if utf8.RuneCountInString(creature.Endpoint) > utf8.RuneCountInString(biggestWord) {
				biggestWord = creature.Endpoint
			}
		}

		biggestCreatureName = biggestName
		biggestCreatureNameRuneCount = utf8.RuneCountInString(biggestName)
		biggestCreatureWord = biggestWord
		biggestCreatureWordRuneCount = utf8.RuneCountInString(biggestWord)
	}
}

// setSpellsVarss sets spells vars
// this only needs to be called once as it will never change during runtime
func setSpellsVars() {
	if smallestSpellNameOrFormula == "" {
		smallestName := val.Spells[0].Name
		var smallestWord string

		for _, spell := range val.Spells {
			if len(spell.Name) < utf8.RuneCountInString(smallestName) {
				smallestName = spell.Name
			}

			if len(spell.Formula) < utf8.RuneCountInString(smallestName) {
				smallestName = spell.Formula
			}

			fieldsName := strings.Fields(spell.Name)

			if utf8.RuneCountInString(smallestWord) == 0 {
				smallestWord = fieldsName[0]
			}

			for _, str := range fieldsName {
				if utf8.RuneCountInString(str) < utf8.RuneCountInString(smallestWord) {
					smallestWord = str
				}
			}

			fieldsFormula := strings.Fields(spell.Formula)

			for _, str := range fieldsFormula {
				if utf8.RuneCountInString(str) < utf8.RuneCountInString(smallestWord) {
					smallestWord = str
				}
			}
		}

		smallestSpellNameOrFormula = smallestName
		smallestSpellNameOrFormulaRuneCount = utf8.RuneCountInString(smallestName)
		smallestSpellWord = smallestWord
		smallestSpellWordRuneCount = utf8.RuneCountInString(smallestWord)
	}

	if biggestSpellNameOrFormula == "" {
		biggestName := val.Spells[0].Name
		var biggestWord string

		for _, spell := range val.Spells {
			if len(spell.Name) > utf8.RuneCountInString(biggestName) {
				biggestName = spell.Name
			}

			if len(spell.Formula) > utf8.RuneCountInString(biggestName) {
				biggestName = spell.Formula
			}

			fieldsName := strings.Fields(spell.Name)

			if utf8.RuneCountInString(biggestWord) == 0 {
				biggestWord = fieldsName[0]
			}

			for _, str := range fieldsName {
				if utf8.RuneCountInString(str) > utf8.RuneCountInString(biggestWord) {
					biggestWord = str
				}
			}

			fieldsFormula := strings.Fields(spell.Formula)

			for _, str := range fieldsFormula {
				if utf8.RuneCountInString(str) > utf8.RuneCountInString(biggestWord) {
					biggestWord = str
				}
			}

			if utf8.RuneCountInString(spell.Endpoint) > utf8.RuneCountInString(biggestWord) {
				biggestWord = spell.Endpoint
			}
		}

		biggestSpellNameOrFormula = biggestName
		biggestSpellNameOrFormulaRuneCount = utf8.RuneCountInString(biggestName)
		biggestSpellWord = biggestWord
		biggestSpellWordRuneCount = utf8.RuneCountInString(biggestWord)
	}
}
