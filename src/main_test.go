package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	const (
		expectedUserAgent     = "TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source)"
		expectedUserAgentHost = "TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source; +https://unittest.example.com)"
	)

	TibiaDataHost = ""
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	assert.Equal(t, expectedUserAgent, TibiaDataUserAgent)

	TibiaDataHost = "unittest.example.com"
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	assert.Equal(t, expectedUserAgentHost, TibiaDataUserAgent)
}

func TestTibiaDataInitializer(t *testing.T) {
	assert := assert.New(t)

	TibiaDataHost = "unittest.example.com"

	// Call the function to be tested
	TibiaDataInitializer()

	// Check that the variables have been set correctly
	assert.Equal("open-source", TibiaDataBuildEdition)
	assert.Equal("https", TibiaDataProtocol)
	assert.Equal("unittest.example.com", TibiaDataHost)
}

func TestTibiaDataInitializerFansiteToken(t *testing.T) {
	prevToken := TibiaFansiteToken
	prevEnabled := TibiaFansiteAPI
	t.Cleanup(func() {
		os.Unsetenv("TIBIA_FANSITEAPI_TOKEN")
		TibiaFansiteToken = prevToken
		TibiaFansiteAPI = prevEnabled
	})

	t.Run("valid token enables fansiteapi", func(t *testing.T) {
		assert := assert.New(t)

		validToken := testJWT(t, tibiaFansiteTokenAlgorithm, map[string]any{
			"nameid": "TibiaData",
			"role":   tibiaFansiteTokenRole,
			"exp":    time.Now().Add(time.Hour).Unix(),
		})
		os.Setenv("TIBIA_FANSITEAPI_TOKEN", validToken)

		TibiaDataInitializer()

		assert.True(TibiaFansiteAPI)
		assert.Equal(validToken, TibiaFansiteToken)
	})

	t.Run("invalid token disables fansiteapi", func(t *testing.T) {
		assert := assert.New(t)

		os.Setenv("TIBIA_FANSITEAPI_TOKEN", "not-a-jwt")

		TibiaDataInitializer()

		assert.False(TibiaFansiteAPI)
		assert.Empty(TibiaFansiteToken)
	})
}
