package validation

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/projectbarks/cimap"
)

const (
	TIBIADATA_API_TESTING = "TibiaData-API-Testing"
	TIBIADATA_JSON_ERROR  = "Wanted json error: %s, got %s instead"
)

var (
	// characterNameRegex is used to check if the character name provided is valid
	// nowadays Tibia only accepts a-zA-Z, but we have to consider old names
	characterNameRegex = regexp.MustCompile(`[^\s'\p{L}\-\.\,\+]`)

	// creatureAndSpellNameRegex is used to check if the creature or spell name provided is valid
	creatureAndSpellNameRegex = regexp.MustCompile(`[^\s'a-zA-Z]`)

	// guildNameRegex is used to check if the guild name provided is valid
	guildNameRegex = regexp.MustCompile(`[^\sa-zA-Z]`)

	// validVocations stores all valid tibia vocations
	validVocations = func() *cimap.CaseInsensitiveMap[bool] {
		m := cimap.New[bool](12)
		for _, v := range []string{"none", "knight", "knights", "paladin", "paladins", "sorcerer", "sorcerers", "druid", "druids", "monk", "monks", "all"} {
			m.Add(v, true)
		}
		return m
	}()
)

// IsRestrictionMode reports whether the restriction mode is enabled
// Check if error == nil to see whether the restriction mode is enabled or not
func IsRestrictionMode(mode bool) error {
	if mode {
		return ErrorRestrictionMode
	}

	return nil
}

// IsNewsIDValid reports wheter the provided int represents a valid news ID
// Check if error == nil to see whether the ID is valid or not
func IsNewsIDValid(ID int) error {
	// If the ID is 0 or a negative number it is invalid
	if ID <= 0 {
		return ErrorInvalidNewsID
	}

	return nil
}

// IsVocationValid reports wheter the provided string represents a valid vocation
// Check if error == nil to see whether the vocation is valid or not
func IsVocationValid(vocation string) error {
	if _, ok := validVocations.Get(vocation); ok {
		return nil
	}

	return ErrorVocationDoesNotExist
}

// IsCharacterNameValid reports wheter the provided string represents a valid character name
// Check if error == nil to see whether the name is valid or not
func IsCharacterNameValid(name string) error {
	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return ErrorCharacterNameEmpty
	case lenName < MinRunesAllowedInACharacterName: // Name is too small
		return ErrorCharacterNameTooSmall
	case lenName > MaxRunesAllowedInACharacterName: // Name is too big
		return ErrorCharacterNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return ErrorCharacterNameIsOnlyWhiteSpace
	}

	// Check if any word in the name has a length > 14
	strs := strings.FieldsFunc(name, isValidCharacterNameSeparator)
	for _, str := range strs {
		if utf8.RuneCountInString(str) > MaxRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooBig
		}

		if utf8.RuneCountInString(str) < MinRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := characterNameRegex.MatchString(name)
	if matched {
		return ErrorCharacterNameInvalid
	}

	return nil
}

func isValidCharacterNameSeparator(r rune) bool {
	switch r {
	case '+', '-', '\'':
		return true
	}
	return unicode.IsSpace(r)
}

// IsGuildNameValid reports wheter the provided string represents a valid guild name
// Check if error == nil to see whether the name is valid or not
func IsGuildNameValid(name string) error {
	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return ErrorGuildNameEmpty
	case lenName < MinRunesAllowedInAGuildName: // Name is too small
		return ErrorGuildNameTooSmall
	case lenName > MaxRunesAllowedInAGuildName: // Name is too big
		return ErrorGuildNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return ErrorGuildNameIsOnlyWhiteSpace
	}

	// Check if any word in the name has a length > 14
	strs := strings.Fields(name)
	for _, str := range strs {
		if utf8.RuneCountInString(str) > MaxRunesAllowedInAGuildNameWord {
			return ErrorGuildWordTooBig
		}

		if utf8.RuneCountInString(str) < MinRunesAllowedInAGuildNameWord {
			return ErrorGuildWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := guildNameRegex.MatchString(name)
	if matched {
		return ErrorGuildNameInvalid
	}

	return nil
}

// IsCreatureNameValid reports wheter the provided string represents a valid creature name
// Check if error == nil to see whether the creature is valid or not
// It will also return the creature endpoint
func IsCreatureNameValid(name string) (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return "", ErrorCreatureNameEmpty
	case lenName < smallestCreatureNameRuneCount: // Name is too small
		return "", ErrorCreatureNameTooSmall
	case lenName > biggestCreatureNameRuneCount: // Name is too big
		return "", ErrorCreatureNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return "", ErrorCreatureNameIsOnlyWhiteSpace
	}

	// Check words length
	strs := strings.Fields(name)
	for _, str := range strs {
		utfCount := utf8.RuneCountInString(str)

		if utfCount > biggestCreatureWordRuneCount {
			return "", ErrorCreatureWordTooBig
		}

		if utfCount < smallestCreatureWordRuneCount {
			return "", ErrorCreatureWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := creatureAndSpellNameRegex.MatchString(name)
	if matched {
		return "", ErrorCreatureNameInvalid
	}

	// Check if creature exists
	if endpoint, ok := creatureLookup.Get(name); ok {
		return endpoint, nil
	}

	return "", ErrorCreatureNotFound
}

// IsSpellNameOrFormulaValid reports wheter the provided string represents a valid spell name or formula
// Check if error == nil to see whether the creature is valid or not
// It will also return the spell endpoint
func IsSpellNameOrFormulaValid(name string) (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return "", ErrorSpellNameEmpty
	case lenName < smallestSpellNameOrFormulaRuneCount: // Name is too small
		return "", ErrorSpellNameTooSmall
	case lenName > biggestSpellNameOrFormulaRuneCount: // Name is too big
		return "", ErrorSpellNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return "", ErrorSpellNameIsOnlyWhiteSpace
	}

	// Check words length
	strs := strings.Fields(name)
	for _, str := range strs {
		utfCount := utf8.RuneCountInString(str)

		if utfCount > biggestSpellWordRuneCount {
			return "", ErrorSpellWordTooBig
		}

		if utfCount < smallestSpellWordRuneCount {
			return "", ErrorSpellWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := creatureAndSpellNameRegex.MatchString(name)
	if matched {
		return "", ErrorSpellNameInvalid
	}

	// Check if spell exists
	if endpoint, ok := spellLookup.Get(name); ok {
		return endpoint, nil
	}

	return "", ErrorSpellNotFound
}

// GetWorlds returns a list of all existing worlds
func GetWorlds() ([]string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Worlds, nil
}

// WorldExists reports whether the specified world exists
// This function is case insensitive
func WorldExists(world string) (bool, error) {
	// Check if the validator has been initiated
	if !initiated {
		return false, ErrorValidatorNotInitiated
	}

	// Try to find the world
	_, found := worldLookup.Get(world)
	return found, nil
}

// GetTowns returns a list of all existing towns
func GetTowns() ([]string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Towns, nil
}

// TowndExists reports whether the specified town exists
// This function is case insensitive
func TownExists(town string) (bool, error) {
	// Check if the validator has been initiated
	if !initiated {
		return false, ErrorValidatorNotInitiated
	}

	// Replace + with spaces in the town name
	town = strings.ReplaceAll(town, "+", " ")

	// Try to find the town
	_, found := townLookup.Get(town)
	return found, nil
}

// GetHouses returns a slice of all houses
func GetHouses() ([]House, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Houses, nil
}

// GetHouseRaw returns a house by it's ID, independently
// of what town the house is from
// This function will return a nil house AND a nil error
// if the specified ID doesn't exist
func GetHouseRaw(houseID int) (House, error) {
	// Check if the validator has been initiated
	if !initiated {
		return House{}, ErrorValidatorNotInitiated
	}

	// Try to find the house
	for _, h := range val.Houses {
		if h.ID == houseID {
			return h, nil
		}
	}

	return House{}, nil
}

// HouseExistsRaw reports whether a house exits, independently
// of what town the house is from
func HouseExistsRaw(houseID int) (bool, error) {
	house, err := GetHouseRaw(houseID)
	return house != House{}, err
}

// GetHouseInTown returns a house by it's ID and town
// This function will return a nil house AND a nil error
// if the specified ID doesn't exist in the specified town
// or if the specified town doesn't exist
func GetHouseInTown(houseID int, town string) (House, error) {
	// We don't need to check if the validator has been initiated
	// because TownExists will already check that for us
	townExists, err := TownExists(town)
	if err != nil {
		return House{}, err
	}

	// Town doesn't exist
	if !townExists {
		return House{}, nil
	}

	// Try to find the house
	for _, h := range val.Houses {
		if h.ID == houseID && strings.EqualFold(h.Town, town) {
			return h, nil
		}
	}

	return House{}, nil
}

// HouseExistsInTown reports whether a house exits in the specified town
// This function will return false AND a nil error if the specified ID
// doesn't exist in the specified town or if the specified town doesn't exist
func HouseExistsInTown(houseID int, town string) (bool, error) {
	house, err := GetHouseInTown(houseID, town)
	return house != House{}, err
}

// GetCreatures returns a list of all existing creatures
func GetCreatures() ([]Creature, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Creatures, nil
}
