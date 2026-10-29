package validation

import (
	"errors"
)

// Error represents a validation error
type Error struct {
	error
}

var (
	////////////////////
	// Server Errors //
	//////////////////

	// ErrorAlreadyRunning will be sent when InitiateValidator() is called but the validator is already running
	// Code: 10
	ErrorAlreadyRunning = Error{errors.New("validator has already been initiated on this session")}

	// ErrorValidatorNotInitiated will be sent when a validation func is called but the validator has not been initiated
	// Code: 11
	ErrorValidatorNotInitiated = Error{errors.New("validator func called but the validator has not been initiated")}

	////////////////////
	/// User Errors ///
	//////////////////

	// ErrorStringCanNotBeConvertedToInt will be sent if the request needs to be converted to an int but can't
	// Code: 9001
	ErrorStringCanNotBeConvertedToInt = Error{errors.New("the provided string can not be converted to an integer")}

	// ErrorRestrictionMode will be sent if the request contains a page that is not available due to restriction mode
	// Code: 9002
	ErrorRestrictionMode = Error{errors.New("the provided page is not available due to restriction mode")}

	// ErrorCharacterNameEmpty will be sent if the request contains an empty character name
	// Code: 10001
	ErrorCharacterNameEmpty = Error{errors.New("the provided character name is an empty string")}

	// ErrorCharacterNameTooSmall will be sent if the request contains a character name of length < MinRunesAllowedInACharacterName
	// Code: 10002
	ErrorCharacterNameTooSmall = Error{errors.New("the provided character name is too small")}

	// ErrorCharacterNameInvalid will be sent if the request contains an invalid character name
	// Code: 10003
	ErrorCharacterNameInvalid = Error{errors.New("the provided character name is invalid")}

	// ErrorCharacterNameIsOnlyWhiteSpace will be sent if the request contains a name that consists of only whitespaces
	// Code: 10004
	ErrorCharacterNameIsOnlyWhiteSpace = Error{errors.New("the provided character name consists only of whitespaces")}

	// ErrorCharacterNameTooBig will be sent if the request contains a character name of length > MaxRunesAllowedInACharacterName
	// Code: 10005
	ErrorCharacterNameTooBig = Error{errors.New("the provided character name is too big")}

	// ErrorCharacterWordTooBig will be sent if the request contains a word with length > MaxRunesAllowedInACharacterNameWord in the character name
	// Code: 10006
	ErrorCharacterWordTooBig = Error{errors.New("the provided character name has a word too big")}

	// ErrorCharacterWordTooSmall will be sent if the request contains a word with length < MinRunesAllowedInACharacterNameWord in the character name
	// Code: 10007
	ErrorCharacterWordTooSmall = Error{errors.New("the provided character name has a word too small")}

	// ErrorInvalidNewsID will be sent if the request contains an invalid news ID
	// Code: 11001
	ErrorInvalidNewsID = Error{errors.New("the provided news id is invalid")}

	// ErrorWorldDoesNotExist will be sent if the request contains a world that does not exist
	// Code: 11002
	ErrorWorldDoesNotExist = Error{errors.New("the provided world does not exist")}

	// ErrorVocationDoesNotExist will be sent if the request contains a vocation that does not exist
	// Code: 11003
	ErrorVocationDoesNotExist = Error{errors.New("the provided vocation does not exist")}

	// ErrorHighscoreCategoryDoesNotExist will be sent if the request contains a highscore catregory that does not exist
	// Code: 11004
	ErrorHighscoreCategoryDoesNotExist = Error{errors.New("the provided highscore category does not exist")}

	// ErrorHouseDoesNotExist will be sent if the request contains a house that does not exist
	// Code: 11005
	ErrorHouseDoesNotExist = Error{errors.New("the provided house does not exist")}

	// ErrorTownDoesNotExist will be sent if the request contains a town that does not exist
	// Code: 11006
	ErrorTownDoesNotExist = Error{errors.New("the provided town does not exist")}

	// ErrorHighscorePageInvalid will be sent if the page is not valid
	// Code: 11007
	ErrorHighscorePageInvalid = Error{errors.New("the provided page does not exist or is invalid")}

	// ErrorHighscorePageTooBig
	// Code: 11008
	ErrorHighscorePageTooBig = Error{errors.New("the provided page is larger than max amount of pages")}

	// ErrorLeaderboardRotationInvalid will be sent if the rotation is neither a number nor current
	// Code: 11009
	ErrorLeaderboardRotationInvalid = Error{errors.New("the provided leaderboard rotation is invalid")}

	// ErrorLeaderboardRotationDoesNotExist will be sent if the rotation is not listed on tibia.com
	// Code: 11010
	ErrorLeaderboardRotationDoesNotExist = Error{errors.New("the provided leaderboard rotation does not exist")}

	// ErrorCreatureNameEmpty will be sent if the request contains an empty creature name
	// Code: 12001
	ErrorCreatureNameEmpty = Error{errors.New("the provided creature name is an empty string")}

	// ErrorCreatureNameTooSmall will be sent if the request contains a creature name of length < smallestCreatureName
	// Code: 12002
	ErrorCreatureNameTooSmall = Error{errors.New("the provided creature name is too smal")}

	// ErrorCreatureNameInvalid will be sent if the request contains an invalid creature name
	// Code: 12003
	ErrorCreatureNameInvalid = Error{errors.New("the provided creature name is invalid")}

	// ErrorCreatureNameIsOnlyWhiteSpace will be sent if the request contains a name that consists of only whitespaces
	// Code: 12004
	ErrorCreatureNameIsOnlyWhiteSpace = Error{errors.New("the provided creature name consists only of whitespaces")}

	// ErrorCreatureNameTooBig will be sent if the request contains a creature name of length > biggestCreatureNameRuneCount
	// Code: 12005
	ErrorCreatureNameTooBig = Error{errors.New("the provided creature name is too big")}

	// ErrorCreatureWordTooBig will be sent if the request contains a word with length > biggestCreatureWordRuneCount in the creature name
	// Code: 12006
	ErrorCreatureWordTooBig = Error{errors.New("the provided creature name has a word too big")}

	// ErrorCreatureWordTooSmall will be sent if the request contains a word with length < smallestCreatureWordRuneCount in the creature name
	// Code: 12007
	ErrorCreatureWordTooSmall = Error{errors.New("the provided creature name has a word too small")}

	// ErrorSpellNameEmpty will be sent if the request contains an empty spell name
	// Code: 13001
	ErrorSpellNameEmpty = Error{errors.New("the provided spell name is an empty string")}

	// ErrorSpellNameTooSmall will be sent if the request contains a spell name of length < smallestSpellNameOrFormulaRuneCount
	// Code: 13002
	ErrorSpellNameTooSmall = Error{errors.New("the provided spell name is too smal")}

	// ErrorSpellNameInvalid will be sent if the request contains an invalid spell name
	// Code: 13003
	ErrorSpellNameInvalid = Error{errors.New("the provided spell name is invalid")}

	// ErrorSpellNameIsOnlyWhiteSpace will be sent if the request contains a name that consists of only whitespaces
	// Code: 13004
	ErrorSpellNameIsOnlyWhiteSpace = Error{errors.New("the provided spell name consists only of whitespaces")}

	// ErrorSpellNameTooBig will be sent if the request contains a spell name of length > biggestSpellNameOrFormulaRuneCount
	// Code: 13005
	ErrorSpellNameTooBig = Error{errors.New("the provided spell name is too big")}

	// ErrorSpellWordTooBig will be sent if the request contains a word with length > biggestSpellWordRuneCount in the spell name
	// Code: 13006
	ErrorSpellWordTooBig = Error{errors.New("the provided spell name has a word too big")}

	// ErrorSpellWordTooSmall will be sent if the request contains a word with length < smallestSpellWordRuneCount in the creature name
	// Code: 13007
	ErrorSpellWordTooSmall = Error{errors.New("the provided spell name has a word too small")}

	// ErrorGuildNameEmpty will be sent if the request contains an empty guild name
	// Code: 14001
	ErrorGuildNameEmpty = Error{errors.New("the provided guild name is an empty string")}

	// ErrorGuildNameTooSmall will be sent if the request contains a Guild name of length < MinRunesAllowedInAGuildName
	// Code: 14002
	ErrorGuildNameTooSmall = Error{errors.New("the provided guild name is too small")}

	// ErrorGuildNameInvalid will be sent if the request contains an invalid guild name
	// Code: 14003
	ErrorGuildNameInvalid = Error{errors.New("the provided guild name is invalid")}

	// ErrorGuildNameIsOnlyWhiteSpace will be sent if the request contains a name that consists of only whitespaces
	// Code: 14004
	ErrorGuildNameIsOnlyWhiteSpace = Error{errors.New("the provided guild name consists only of whitespaces")}

	// ErrorGuildNameTooBig will be sent if the request contains a guild name of length > MaxRunesAllowedInAGuildName
	// Code: 14005
	ErrorGuildNameTooBig = Error{errors.New("the provided guild name is too big")}

	// ErrorGuildWordTooBig will be sent if the request contains a word with length > MaxRunesAllowedInAGuildNameWord in the guild name
	// Code: 14006
	ErrorGuildWordTooBig = Error{errors.New("the provided guild name has a word too big")}

	// ErrorGuildWordTooSmall will be sent if the request contains a word with length < MinRunesAllowedInAGuildNameWord in the guild name
	// Code: 14007
	ErrorGuildWordTooSmall = Error{errors.New("the provided guild name has a word too smal")}

	///////////////////
	// Tibia Errors //
	/////////////////

	// ErrorCharacterNotFound will be sent if the requested character does not exist
	// Code: 20001
	ErrorCharacterNotFound = Error{errors.New("could not find character")}

	// ErrorCreatureNotFound will be sent if the requested creature does not exist
	// Code: 20002
	ErrorCreatureNotFound = Error{errors.New("could not find creature")}

	// ErrorSpellNotFound will be sent if the requested spell does not exist
	// Code: 20003
	ErrorSpellNotFound = Error{errors.New("could not find spell")}

	// ErrorGuildNotFound will be sent if the requested guild does not exist
	// Code: 20004
	ErrorGuildNotFound = Error{errors.New("could not find guild")}

	// ErrorMaintenanceMode will be sent if there is ongoing maintenance
	// Code: 20005
	ErrorMaintenanceMode = Error{errors.New("maintenance mode active")}

	// ErrStatusForbidden will be sent if tibia sent us a 403 response.
	// This usually happens when we are rate limited.
	// Code: 20006
	ErrStatusForbidden = Error{errors.New("got status forbidden from tibia.com")}

	// ErrStatusForbidden will be sent if tibia sent us a 302 response, but it
	// is not in MaintenanceMode. Because if it were, we would throw a
	// ErrorMaintenanceMode.
	// Code: 20007
	ErrStatusFound = Error{errors.New("got status found from tibia.com")}

	// ErrStatusUnknown will be sent a HTTP request we are not expecting.
	// Code: 20008
	ErrStatusUnknown = Error{errors.New("got unknown status from tibia.com")}
)

// Code will return the code of the error
func (e Error) Code() int {
	switch e {
	case ErrorAlreadyRunning:
		return 10
	case ErrorValidatorNotInitiated:
		return 11
	case ErrorStringCanNotBeConvertedToInt:
		return 9001
	case ErrorRestrictionMode:
		return 9002
	case ErrorCharacterNameEmpty:
		return 10001
	case ErrorCharacterNameTooSmall:
		return 10002
	case ErrorCharacterNameInvalid:
		return 10003
	case ErrorCharacterNameIsOnlyWhiteSpace:
		return 10004
	case ErrorCharacterNameTooBig:
		return 10005
	case ErrorCharacterWordTooBig:
		return 10006
	case ErrorCharacterWordTooSmall:
		return 10007
	case ErrorInvalidNewsID:
		return 11001
	case ErrorWorldDoesNotExist:
		return 11002
	case ErrorVocationDoesNotExist:
		return 11003
	case ErrorHighscoreCategoryDoesNotExist:
		return 11004
	case ErrorHouseDoesNotExist:
		return 11005
	case ErrorTownDoesNotExist:
		return 11006
	case ErrorHighscorePageInvalid:
		return 11007
	case ErrorHighscorePageTooBig:
		return 11008
	case ErrorLeaderboardRotationInvalid:
		return 11009
	case ErrorLeaderboardRotationDoesNotExist:
		return 11010
	case ErrorCreatureNameEmpty:
		return 12001
	case ErrorCreatureNameTooSmall:
		return 12002
	case ErrorCreatureNameInvalid:
		return 12003
	case ErrorCreatureNameIsOnlyWhiteSpace:
		return 12004
	case ErrorCreatureNameTooBig:
		return 12005
	case ErrorCreatureWordTooBig:
		return 12006
	case ErrorCreatureWordTooSmall:
		return 12007
	case ErrorSpellNameEmpty:
		return 13001
	case ErrorSpellNameTooSmall:
		return 13002
	case ErrorSpellNameInvalid:
		return 13003
	case ErrorSpellNameIsOnlyWhiteSpace:
		return 13004
	case ErrorSpellNameTooBig:
		return 13005
	case ErrorSpellWordTooBig:
		return 13006
	case ErrorSpellWordTooSmall:
		return 13007
	case ErrorGuildNameEmpty:
		return 14001
	case ErrorGuildNameTooSmall:
		return 14002
	case ErrorGuildNameInvalid:
		return 14003
	case ErrorGuildNameIsOnlyWhiteSpace:
		return 14004
	case ErrorGuildNameTooBig:
		return 14005
	case ErrorGuildWordTooBig:
		return 14006
	case ErrorGuildWordTooSmall:
		return 14007
	case ErrorCharacterNotFound:
		return 20001
	case ErrorCreatureNotFound:
		return 20002
	case ErrorSpellNotFound:
		return 20003
	case ErrorGuildNotFound:
		return 20004
	case ErrorMaintenanceMode:
		return 20005
	case ErrStatusForbidden:
		return 20006
	case ErrStatusFound:
		return 20007
	case ErrStatusUnknown:
		return 20008
	default:
		return 0
	}
}
