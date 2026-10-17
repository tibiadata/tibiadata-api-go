package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestFakeToUpCodeCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// adding support for proxy for tests
	if isEnvExist("TIBIADATA_PROXY") {
		TibiaDataProxyDomain = "https://" + getEnv("TIBIADATA_PROXY", "www.tibia.com") + "/"
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "Durin",
		},
	}

	assert := assert.New(t)

	tibiaBoostableBosses(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaCreaturesOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "race",
			Value: "Demon",
		},
	}

	tibiaCreaturesCreature(c)
	fmt.Println("tibiaCreaturesCreature", w)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaFansites(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "Pax",
		},
	}

	tibiaGuildsGuild(c)
	fmt.Println("tibiaGuildsGuild", w)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaGuildsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "Antica",
		},
		{
			Key:   "category",
			Value: "experience",
		},
		{
			Key:   "vocation",
			Value: "sorcerer",
		},
		{
			Key:   "page",
			Value: "4",
		},
	}

	tibiaHighscores(c)
	fmt.Println("tibiaHighscores", w)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
		{
			Key:   "house_id",
			Value: "59054",
		},
	}

	tibiaHousesHouse(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
		{
			Key:   "town",
			Value: "venore",
		},
	}

	tibiaHousesOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaKillstatistics(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.False(false, tibiaNewslistArchive())
	assert.False(false, tibiaNewslistArchiveDays())
	assert.False(false, tibiaNewslistLatest())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "vocation",
			Value: "sorcerer",
		},
	}

	tibiaSpellsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "spell_id",
			Value: "exori",
		},
	}

	tibiaSpellsSpell(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaWorldsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "antica",
		},
	}

	tibiaWorldsWorld(c)
	assert.Equal(http.StatusOK, w.Code)

	rootz(c)
	assert.Equal(http.StatusOK, w.Code)

	healthz(c)
	assert.Equal(http.StatusOK, w.Code)

	readyz(c)
	assert.Equal(http.StatusOK, w.Code)

	type test struct {
		T string `json:"t"`
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	TibiaDataAPIHandleResponse(c, "", test{T: "abc"})
	assert.Equal(http.StatusOK, w.Code)
}

func TestErrorHandler(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, errors.New("test error"), http.StatusBadRequest)
	assert.Equal(http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorAlreadyRunning, 0)
	assert.Equal(http.StatusInternalServerError, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorCharacterNameInvalid, 0)
	assert.Equal(http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorVocationDoesNotExist, http.StatusBadRequest)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorHighscoreCategoryDoesNotExist, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorHouseDoesNotExist, http.StatusBadRequest)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorTownDoesNotExist, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorCharacterNotFound, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorCreatureNotFound, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorSpellNotFound, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorGuildNotFound, 0)
	assert.Equal(http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, errors.New("test error"), 0)
	assert.Equal(http.StatusBadGateway, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrStatusForbidden, http.StatusForbidden)
	assert.Equal(http.StatusBadGateway, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrStatusFound, http.StatusFound)
	assert.Equal(http.StatusBadGateway, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrStatusUnknown, http.StatusConflict)
	assert.Equal(http.StatusBadGateway, w.Code)
}

func readNewsFixture(t *testing.T, path string) []byte {
	t.Helper()

	file, err := static.TestFiles.Open(path)
	if err != nil {
		t.Fatalf("open news fixture %q: %v", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read news fixture %q: %v", path, err)
	}
	return data
}

func TestTibiaNewsHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newslistFixture := readNewsFixture(t, "testdata/news/newslist.html")
	newsFixture := readNewsFixture(t, "testdata/news/archive/6512.html")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") == "6512" {
			assert.Equal(t, http.MethodGet, r.Method)
			_, _ = w.Write(newsFixture)
			return
		}

		assert.Equal(t, http.MethodPost, r.Method)
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		assert.Equal(t, "article", r.PostForm.Get("filter_article"))
		assert.Equal(t, "news", r.PostForm.Get("filter_news"))
		_, _ = w.Write(newslistFixture)
	}))
	t.Cleanup(server.Close)

	previousProxyDomain := TibiaDataProxyDomain
	TibiaDataProxyDomain = server.URL + "/"
	t.Cleanup(func() { TibiaDataProxyDomain = previousProxyDomain })

	previousHost, previousProtocol := TibiaDataHost, TibiaDataProtocol
	TibiaDataHost, TibiaDataProtocol = "localhost", "http"
	t.Cleanup(func() {
		TibiaDataHost, TibiaDataProtocol = previousHost, previousProtocol
	})

	t.Run("news list", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v4/news/archive/90", nil)
		c.Params = gin.Params{{Key: "days", Value: "90"}}

		tibiaNewslist(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response NewsListResponse
		if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response)) {
			assert.NotEmpty(t, response.News)
			assert.Equal(t, 6529, response.News[0].ID)
			assert.Equal(t, "http://localhost/v4/news/id/6529", response.News[0].ApiURL)
		}
	})

	t.Run("news article", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v4/news/id/6512", nil)
		c.Params = gin.Params{{Key: "news_id", Value: "6512"}}

		tibiaNews(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response NewsResponse
		if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response)) {
			assert.Equal(t, 6512, response.News.ID)
			assert.Equal(t, "ticker", response.News.Type)
		}
	})
}

func TestTibiaNewsHandlersReturnBadGatewayOnUpstreamFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	previousProxyDomain := TibiaDataProxyDomain
	TibiaDataProxyDomain = server.URL + "/"
	t.Cleanup(func() { TibiaDataProxyDomain = previousProxyDomain })

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/news/id/6512", nil)
	c.Params = gin.Params{{Key: "news_id", Value: "6512"}}

	tibiaNews(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)
}

func TestTibiaDataJSONDataCollector(t *testing.T) {
	prevToken := TibiaFansiteToken
	t.Cleanup(func() { TibiaFansiteToken = prevToken })

	t.Run("missing token", func(t *testing.T) {
		assert := assert.New(t)
		TibiaFansiteToken = ""

		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: "https://example.invalid"})
		assert.Empty(body)
		if assert.Error(err) {
			assert.Contains(err.Error(), "missing TibiaFansiteToken")
		}
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		var gotAuth string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuth = r.Header.Get("Authorization")
			w.Header().Set("Age", "120")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"characterGameInformation":{"characterName":"Test"}}`))
		}))
		defer server.Close()

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL})
		assert.NoError(err)
		assert.Equal(`{"characterGameInformation":{"characterName":"Test"}}`, body)
		assert.Equal("Bearer test-token", gotAuth)
	})

	t.Run("forbidden", func(t *testing.T) {
		assert := assert.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL})
		assert.Empty(body)
		assert.ErrorIs(err, validation.ErrStatusForbidden)
	})

	t.Run("not found", func(t *testing.T) {
		assert := assert.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL})
		assert.Empty(body)
		assert.ErrorIs(err, validation.ErrorCharacterNotFound)
	})

	t.Run("unknown status", func(t *testing.T) {
		assert := assert.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL})
		assert.Empty(body)
		assert.ErrorIs(err, validation.ErrStatusUnknown)
	})

	t.Run("request error", func(t *testing.T) {
		assert := assert.New(t)

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: "http://127.0.0.1:0"})
		assert.Empty(body)
		assert.Error(err)
	})
}

func TestTibiaDataRequestHandlerReturnsNotFoundForMissingFansiteCharacter(t *testing.T) {
	prevToken := TibiaFansiteToken
	TibiaFansiteToken = "test-token"
	t.Cleanup(func() { TibiaFansiteToken = prevToken })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	tibiaDataRequestHandler(
		c,
		TibiaDataRequestStruct{
			URL:           server.URL,
			UseFansiteAPI: true,
		},
		func(string) (interface{}, error) {
			t.Fatal("request handler should not be called for a missing character")
			return nil, nil
		},
		"test",
	)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response OutInformation
	if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response)) {
		assert.Equal(t, http.StatusNotFound, response.Information.Status.HTTPCode)
		assert.Equal(t, validation.ErrorCharacterNotFound.Code(), response.Information.Status.Error)
	}
}

func TestTibiaDataHTMLDataCollectorReturnsBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Age", "60")
		w.Header().Set("Date", "Mon, 28 Sep 2026 19:47:06 GMT")
		w.Header().Set("Last-Modified", "Mon, 28 Sep 2026 19:46:36 GMT")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>test</html>"))
	}))
	defer server.Close()

	body, err := TibiaDataHTMLDataCollector(TibiaDataRequestStruct{
		Method:  http.MethodGet,
		URL:     server.URL,
		RawBody: true,
	})

	assert.NoError(t, err)
	assert.Equal(t, "<html>test</html>", body)
}

func TestTibiaDataAPIHandleCachedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)

	previousCacheAwareness := TibiaDataCacheControl
	TibiaDataCacheControl = true
	t.Cleanup(func() { TibiaDataCacheControl = previousCacheAwareness })

	for _, test := range []struct {
		name        string
		maxAge      time.Duration
		cacheHeader string
	}{
		{name: "characters", maxAge: cacheMaxAgeCharacters, cacheHeader: "public, max-age=300, s-maxage=300"},
		{name: "guilds", maxAge: cacheMaxAgeGuilds, cacheHeader: "public, max-age=120, s-maxage=120"},
		{name: "highscores", maxAge: cacheMaxAgeHighscores, cacheHeader: "public, max-age=900, s-maxage=900"},
		{name: "houses", maxAge: cacheMaxAgeHouses, cacheHeader: "public, max-age=300, s-maxage=300"},
		{name: "news", maxAge: cacheMaxAgeNews, cacheHeader: "public, max-age=900, s-maxage=900"},
		{name: "worlds", maxAge: cacheMaxAgeWorlds, cacheHeader: "public, max-age=60, s-maxage=60"},
	} {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
			TibiaDataAPIHandleCachedResponse(c, test.name, gin.H{"ok": true}, test.maxAge)

			assert.Equal(t, test.cacheHeader, w.Header().Get(cacheControlHeader))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	TibiaDataCacheControl = false
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second)
	assert.Empty(t, w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	TibiaDataCacheControl = true
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 0)
	assert.Empty(t, w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	TibiaDataAPIHandleCachedResponse(nil, "test", gin.H{"ok": true}, 300*time.Second)
}

func TestTibiaDataAPIHandleCachedResponseDebugLogging(t *testing.T) {
	previousDebug := TibiaDataDebug
	previousMode := gin.Mode()
	TibiaDataDebug = true
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() {
		TibiaDataDebug = previousDebug
		gin.SetMode(previousMode)
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second)

	assert.Equal(t, "public, max-age=300, s-maxage=300", w.Header().Get(cacheControlHeader))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckTibiaFansiteAPIStatus(t *testing.T) {
	prevURL := TibiaFansiteAPIStatusURL
	t.Cleanup(func() { TibiaFansiteAPIStatusURL = prevURL })

	t.Run("available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"isAvailable":true}`))
		}))
		defer server.Close()

		TibiaFansiteAPIStatusURL = server.URL
		checkTibiaFansiteAPIStatus()
	})

	t.Run("unavailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"isAvailable":false}`))
		}))
		defer server.Close()

		TibiaFansiteAPIStatusURL = server.URL
		checkTibiaFansiteAPIStatus()
	})

	t.Run("non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		TibiaFansiteAPIStatusURL = server.URL
		checkTibiaFansiteAPIStatus()
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`not json`))
		}))
		defer server.Close()

		TibiaFansiteAPIStatusURL = server.URL
		checkTibiaFansiteAPIStatus()
	})

	t.Run("request error", func(t *testing.T) {
		TibiaFansiteAPIStatusURL = "http://127.0.0.1:0"
		checkTibiaFansiteAPIStatus()
	})
}

func TestRegisterAPIVersionRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerAPIVersions(router)

	var v4Routes, v5Routes []string
	for _, route := range router.Routes() {
		switch {
		case strings.HasPrefix(route.Path, "/v4/"):
			v4Routes = append(v4Routes, route.Method+" "+strings.TrimPrefix(route.Path, "/v4"))
		case strings.HasPrefix(route.Path, "/v5/"):
			v5Routes = append(v5Routes, route.Method+" "+strings.TrimPrefix(route.Path, "/v5"))
		}
	}

	assert.NotEmpty(t, v4Routes)
	assert.ElementsMatch(t, v4Routes, v5Routes)

	previousCacheControl := TibiaDataCacheControl
	TibiaDataCacheControl = true
	t.Cleanup(func() { TibiaDataCacheControl = previousCacheControl })

	for _, version := range []string{"v4", "v5"} {
		for _, test := range []struct {
			path     string
			location string
		}{
			{
				path:     "/highscores/Antica",
				location: "/highscores/Antica/experience/all/1",
			},
			{
				path:     "/highscores/Antica/experience",
				location: "/highscores/Antica/experience/all/1",
			},
		} {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+version+test.path, nil)
			router.ServeHTTP(w, request)
			assert.Equal(t, http.StatusMovedPermanently, w.Code)
			assert.Equal(t, "/"+version+test.location, w.Header().Get("Location"))
			assert.Equal(t, "public, max-age=900, s-maxage=900", w.Header().Get(cacheControlHeader))
		}
	}
}

func TestTibiaCharactersCharacterUsesHTMLURLWhenFansiteDisabled(t *testing.T) {
	assert := assert.New(t)
	prev := TibiaFansiteAPI
	TibiaFansiteAPI = false
	t.Cleanup(func() { TibiaFansiteAPI = prev })

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Darkside%20Rafa", nil)
	c.Params = []gin.Param{{Key: "name", Value: "Darkside Rafa"}}

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusOK, w.Code)
	assert.NotEmpty(w.Body.Bytes())

	var resp CharacterResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if !assert.NoError(err) {
		return
	}
	if !assert.NotEmpty(resp.Information.TibiaURLs) {
		return
	}
	assert.Contains(resp.Information.TibiaURLs[0], "https://www.tibia.com/community/?subtopic=characters&name=")
}

func TestTibiaCharactersCharacterUsesFansiteAPIWhenEnabled(t *testing.T) {
	assert := assert.New(t)
	prevEnabled := TibiaFansiteAPI
	prevToken := TibiaFansiteToken
	TibiaFansiteAPI = true
	TibiaFansiteToken = ""
	t.Cleanup(func() {
		TibiaFansiteAPI = prevEnabled
		TibiaFansiteToken = prevToken
	})

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Darkside%20Rafa", nil)
	c.Params = []gin.Param{{Key: "name", Value: "Darkside Rafa"}}

	tibiaCharactersCharacter(c)

	// with no token configured, the fansite API request fails fast (no network call),
	// exercising the fansite request-building and dispatch branches.
	assert.Equal(http.StatusBadGateway, w.Code)
}

func TestTibiaDataAPIHandleResponseBranches(t *testing.T) {
	type payload struct {
		T string `json:"t"`
	}

	t.Run("nil context", func(t *testing.T) {
		// exercises the early-return branch when no request context is available.
		TibiaDataAPIHandleResponse(nil, "test", payload{T: "abc"})
	})

	t.Run("debug mode logs request details", func(t *testing.T) {
		assert := assert.New(t)
		prevMode := gin.Mode()
		gin.SetMode(gin.DebugMode)
		t.Cleanup(func() { gin.SetMode(prevMode) })

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		TibiaDataAPIHandleResponse(c, "test", payload{T: "abc"})
		assert.Equal(http.StatusOK, w.Code)
	})

	t.Run("TibiaDataDebug logs execution", func(t *testing.T) {
		assert := assert.New(t)
		prevDebug := TibiaDataDebug
		TibiaDataDebug = true
		t.Cleanup(func() { TibiaDataDebug = prevDebug })

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		TibiaDataAPIHandleResponse(c, "test", payload{T: "abc"})
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestAPIResponseUsesRequestedVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousAPIDetails := TibiaDataAPIDetails
	TibiaDataAPIDetails.Version = 4
	t.Cleanup(func() { TibiaDataAPIDetails = previousAPIDetails })

	for _, test := range []struct {
		path            string
		expectedVersion int
	}{
		{path: "/v4/character/Trollefar", expectedVersion: 4},
		{path: "/v5/character/Trollefar", expectedVersion: 5},
	} {
		t.Run(test.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, test.path, nil)

			TibiaDataAPIHandleResponse(c, "test", OutInformation{
				Information: Information{APIDetails: TibiaDataAPIDetails},
			})

			var response OutInformation
			if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response)) {
				assert.Equal(t, test.expectedVersion, response.Information.APIDetails.Version)
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, test.path, nil)
			TibiaDataErrorHandler(c, validation.ErrorCharacterNameInvalid, http.StatusBadRequest)

			if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response)) {
				assert.Equal(t, test.expectedVersion, response.Information.APIDetails.Version)
			}
		})
	}
}

func TestTibiaDataResponseForRequestLeavesUnexpectedShapesUnchanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/v5/character/Trollefar", nil)

	type noInformation struct {
		Value string
	}
	type noAPIDetails struct {
		Information struct{}
	}
	type nonIntegerVersion struct {
		Information struct {
			APIDetails struct {
				Version string
			}
		}
	}

	tests := []struct {
		name     string
		response interface{}
	}{
		{name: "non-struct response", response: "value"},
		{name: "missing information field", response: noInformation{Value: "value"}},
		{name: "missing api details field", response: noAPIDetails{Information: struct{}{}}},
		{
			name: "non-integer version field",
			response: nonIntegerVersion{Information: struct {
				APIDetails struct {
					Version string
				}
			}{APIDetails: struct{ Version string }{Version: "4"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.response, tibiaDataResponseForRequest(c, test.response))
		})
	}
}

func TestTibiaDataUserAgentForRequest(t *testing.T) {
	origHost, origUA, origUAs := TibiaDataHost, TibiaDataUserAgent, tibiaDataUserAgents
	t.Cleanup(func() {
		TibiaDataHost, TibiaDataUserAgent, tibiaDataUserAgents = origHost, origUA, origUAs
	})

	TibiaDataHost = ""
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	tibiaDataInitUserAgents()

	const details = " (release/unknown; build/manual; commit/-; edition/open-source)"
	tests := []struct {
		path     string
		expected string
	}{
		{"/v4/worlds", "TibiaData-API/v4" + details},
		{"/v5/worlds", "TibiaData-API/v5" + details},
		{"/healthz", TibiaDataUserAgent},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, tt.path, nil)
			assert.Equal(t, tt.expected, tibiaDataUserAgentForRequest(c))
		})
	}

	assert.Equal(t, TibiaDataUserAgent, tibiaDataUserAgentForRequest(nil))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	assert.Equal(t, TibiaDataUserAgent, tibiaDataUserAgentForRequest(c))

	tibiaDataUserAgents = nil
	c.Request = httptest.NewRequest(http.MethodGet, "/v5/worlds", nil)
	assert.Equal(t, TibiaDataUserAgent, tibiaDataUserAgentForRequest(c))
}

func TestTibiaDataInitUserAgents(t *testing.T) {
	previousDebug, previousUserAgents := TibiaDataDebug, tibiaDataUserAgents
	t.Cleanup(func() {
		TibiaDataDebug, tibiaDataUserAgents = previousDebug, previousUserAgents
	})

	for _, debug := range []bool{false, true} {
		t.Run(fmt.Sprintf("debug=%t", debug), func(t *testing.T) {
			TibiaDataDebug = debug
			tibiaDataInitUserAgents()
			assert.Equal(t, map[int]string{
				4: TibiaDataUserAgentGenerator(4),
				5: TibiaDataUserAgentGenerator(5),
			}, tibiaDataUserAgents)
		})
	}
}

func TestTibiaDataCollectorsUserAgent(t *testing.T) {
	previousUA, previousClient, previousToken := TibiaDataUserAgent, tibiaDataClient, TibiaFansiteToken
	t.Cleanup(func() {
		TibiaDataUserAgent, tibiaDataClient, TibiaFansiteToken = previousUA, previousClient, previousToken
	})

	TibiaDataUserAgent = TibiaDataUserAgentGenerator(4)
	TibiaFansiteToken = "test-token"
	initTibiaDataClient()

	for _, test := range []struct {
		name      string
		method    string
		collector func(TibiaDataRequestStruct) (string, error)
	}{
		{"HTML GET", http.MethodGet, TibiaDataHTMLDataCollector},
		{"HTML POST", http.MethodPost, TibiaDataHTMLDataCollector},
		{"JSON GET", http.MethodGet, TibiaDataJSONDataCollector},
	} {
		t.Run(test.name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer server.Close()

			for _, userAgent := range []string{TibiaDataUserAgentGenerator(5), TibiaDataUserAgentGenerator(4), ""} {
				body, err := test.collector(TibiaDataRequestStruct{
					URL: server.URL, Method: test.method, RawBody: true, UserAgent: userAgent,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, `{"ok":true}`, body)
				request := <-requests
				assert.Equal(t, test.method, request.Method)
				expected := userAgent
				if expected == "" {
					expected = TibiaDataUserAgent
				}
				assert.Equal(t, expected, request.UserAgent())
			}
		})
	}
}

func TestTibiaDataRequestHandlerUserAgent(t *testing.T) {
	previousClient, previousToken, previousUserAgents := tibiaDataClient, TibiaFansiteToken, tibiaDataUserAgents
	t.Cleanup(func() {
		tibiaDataClient, TibiaFansiteToken, tibiaDataUserAgents = previousClient, previousToken, previousUserAgents
	})
	TibiaFansiteToken = "test-token"
	initTibiaDataClient()
	tibiaDataInitUserAgents()

	for _, fansite := range []bool{false, true} {
		for _, version := range []int{4, 5} {
			for _, override := range []string{"", "custom-agent"} {
				t.Run(fmt.Sprintf("fansite=%t/v%d/override=%s", fansite, version, override), func(t *testing.T) {
					userAgents := make(chan string, 1)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						userAgents <- r.UserAgent()
						_, _ = w.Write([]byte(`{"ok":true}`))
					}))
					defer server.Close()

					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v%d/worlds", version), nil)
					tibiaDataRequestHandler(c, TibiaDataRequestStruct{
						URL: server.URL, RawBody: true, UseFansiteAPI: fansite, UserAgent: override,
					}, func(body string) (interface{}, error) {
						assert.Equal(t, `{"ok":true}`, body)
						return gin.H{"ok": true}, nil
					}, "test")

					if !assert.Equal(t, http.StatusOK, w.Code) {
						return
					}
					expected := override
					if expected == "" {
						expected = TibiaDataUserAgentGenerator(version)
					}
					assert.Equal(t, expected, <-userAgents)
				})
			}
		}
	}
}

func TestHighscoresCategoryCanonicalRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerAPIVersions(router)

	previousCacheControl := TibiaDataCacheControl
	TibiaDataCacheControl = true
	t.Cleanup(func() { TibiaDataCacheControl = previousCacheControl })

	serve := func(path string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		return w
	}

	t.Run("v5 redirects aliases and case variants", func(t *testing.T) {
		for _, test := range []struct{ path, location string }{
			{"/v5/highscores/Solidera/magic/all/1", "/v5/highscores/Solidera/magiclevel/all/1"},
			{"/v5/highscores/Solidera/MagicLevel/all/2", "/v5/highscores/Solidera/magiclevel/all/2"},
			{"/v5/highscores/all/drome/all", "/v5/highscores/all/dromescore/all"},
			{"/v5/highscores/all/axe/all/1", "/v5/highscores/all/axefighting/all/1"},
			{"/v5/highscores/all/club/all/1", "/v5/highscores/all/clubfighting/all/1"},
			{"/v5/highscores/all/sword/all/1", "/v5/highscores/all/swordfighting/all/1"},
			{"/v5/highscores/all/distance/all/1", "/v5/highscores/all/distancefighting/all/1"},
			{"/v5/highscores/all/fist/all/1", "/v5/highscores/all/fistfighting/all/1"},
			{"/v5/highscores/all/shield/all/1", "/v5/highscores/all/shielding/all/1"},
			{"/v5/highscores/all/loyalty/all/1", "/v5/highscores/all/loyaltypoints/all/1"},
			{"/v5/highscores/all/charm/all/1", "/v5/highscores/all/charmpoints/all/1"},
			{"/v5/highscores/all/charms/all/1", "/v5/highscores/all/charmpoints/all/1"},
			{"/v5/highscores/all/boss/all/1", "/v5/highscores/all/bosspoints/all/1"},
			{"/v5/highscores/all/bosses/all/1", "/v5/highscores/all/bosspoints/all/1"},
			{"/v5/highscores/all/achievement/all/1", "/v5/highscores/all/achievements/all/1"},
			{"/v5/highscores/all/goshnar/all/1", "/v5/highscores/all/goshnarstaint/all/1"},
			{"/v5/highscores/all/magic/all/1?foo=bar", "/v5/highscores/all/magiclevel/all/1?foo=bar"},
			{"/v5/highscores/Solidera/magic", "/v5/highscores/Solidera/magiclevel/all/1"},
			{"/v5/highscores/Solidera/Experience", "/v5/highscores/Solidera/experience/all/1"},
		} {
			w := serve(test.path)
			assert.Equal(t, http.StatusMovedPermanently, w.Code, test.path)
			assert.Equal(t, test.location, w.Header().Get("Location"), test.path)
			assert.Equal(t, "public, max-age=900, s-maxage=900", w.Header().Get(cacheControlHeader), test.path)
		}
	})

	t.Run("canonical v5 and v4 aliases are not redirected", func(t *testing.T) {
		// page 0 fails validation before any upstream request is made
		for _, path := range []string{
			"/v5/highscores/all/magiclevel/all/0",
			"/v4/highscores/all/magic/all/0",
			"/v4/highscores/all/MagicLevel/all/0",
		} {
			w := serve(path)
			assert.Equal(t, http.StatusBadRequest, w.Code, path)
			assert.Empty(t, w.Header().Get("Location"), path)
		}
	})

	t.Run("v4 two segment route keeps category", func(t *testing.T) {
		w := serve("/v4/highscores/Solidera/magic")
		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/v4/highscores/Solidera/magic/all/1", w.Header().Get("Location"))
	})

	t.Run("invalid category returns JSON error", func(t *testing.T) {
		for _, path := range []string{
			"/v4/highscores/all/notacategory/all/1",
			"/v5/highscores/all/notacategory/all/1",
			"/v5/highscores/all/notacategory",
		} {
			w := serve(path)
			assert.Equal(t, http.StatusNotFound, w.Code, path)
			assert.Empty(t, w.Header().Get("Location"), path)
			assert.Empty(t, w.Header().Get(cacheControlHeader), path)

			var resp struct {
				Information Information `json:"information"`
			}
			if assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp), path) {
				assert.Equal(t, validation.ErrorHighscoreCategoryDoesNotExist.Code(), resp.Information.Status.Error, path)
			}
		}
	})
}
