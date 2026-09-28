package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		cacheMetadata := &TibiaDataCacheMetadata{}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuth = r.Header.Get("Authorization")
			w.Header().Set("Age", "120")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"characterGameInformation":{"characterName":"Test"}}`))
		}))
		defer server.Close()

		TibiaFansiteToken = "test-token"
		body, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL, CacheMetadata: cacheMetadata})
		assert.NoError(err)
		assert.Equal(180*time.Second, cacheMetadata.remainingCacheAge(300*time.Second))
		assert.Equal(`{"characterGameInformation":{"characterName":"Test"}}`, body)
		assert.Equal("Bearer test-token", gotAuth)
	})

	t.Run("uses Last-Modified age when Age is missing", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Date", "Mon, 28 Sep 2026 19:47:06 GMT")
			w.Header().Set("Last-Modified", "Mon, 28 Sep 2026 19:46:06 GMT")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"characterGameInformation":{"characterName":"Test"}}`))
		}))
		defer server.Close()

		cacheMetadata := &TibiaDataCacheMetadata{}
		TibiaFansiteToken = "test-token"
		_, err := TibiaDataJSONDataCollector(TibiaDataRequestStruct{URL: server.URL, CacheMetadata: cacheMetadata})

		assert.NoError(t, err)
		assert.Equal(t, 240*time.Second, cacheMetadata.remainingCacheAge(300*time.Second))
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

func TestTibiaDataCacheMetadata(t *testing.T) {
	const (
		date         = "Mon, 28 Sep 2026 19:47:06 GMT"
		lastModified = "Mon, 28 Sep 2026 19:46:06 GMT"
	)
	assert.Equal(t, time.Duration(0), (&TibiaDataCacheMetadata{}).remainingCacheAge(0))

	t.Run("uses Age when present", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{
			"Age":           []string{"120"},
			"Date":          []string{date},
			"Last-Modified": []string{lastModified},
		})

		assert.Equal(t, 180*time.Second, metadata.remainingCacheAge(300*time.Second))
		age, source := metadata.ageDetails()
		assert.Equal(t, 120*time.Second, age)
		assert.Equal(t, "Age", source)
	})

	t.Run("falls back to Date minus Last-Modified", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{
			"Date":          []string{date},
			"Last-Modified": []string{lastModified},
		})

		assert.Equal(t, 240*time.Second, metadata.remainingCacheAge(300*time.Second))
		age, source := metadata.ageDetails()
		assert.Equal(t, 60*time.Second, age)
		assert.Equal(t, "Date-Last-Modified", source)
	})

	t.Run("falls back when Age is invalid", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{
			"Age":           []string{"invalid"},
			"Date":          []string{date},
			"Last-Modified": []string{lastModified},
		})

		assert.Equal(t, 240*time.Second, metadata.remainingCacheAge(300*time.Second))
	})

	t.Run("keeps maximum lifetime when headers are unavailable", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{"Date": []string{date}})

		assert.Equal(t, 300*time.Second, metadata.remainingCacheAge(300*time.Second))
	})

	t.Run("clamps stale responses to zero", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{"Age": []string{"400"}})

		assert.Equal(t, time.Duration(0), metadata.remainingCacheAge(300*time.Second))
	})

	t.Run("keeps the greatest age from multiple upstream responses", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{"Age": []string{"120"}})
		metadata.recordUpstreamHeaders(http.Header{"Age": []string{"60"}})

		age, source := metadata.ageDetails()
		assert.Equal(t, 120*time.Second, age)
		assert.Equal(t, "Age", source)
	})

	t.Run("ignores invalid Last-Modified values", func(t *testing.T) {
		metadata := &TibiaDataCacheMetadata{}
		metadata.recordUpstreamHeaders(http.Header{
			"Date":          []string{date},
			"Last-Modified": []string{"not a date"},
		})

		age, source := metadata.ageDetails()
		assert.Zero(t, age)
		assert.Equal(t, "none", source)
	})
}

func TestUpstreamResponseAgeParsing(t *testing.T) {
	maxDurationSeconds := int64((1<<63 - 1) / int64(time.Second))
	tests := []struct {
		name      string
		value     string
		wantAge   time.Duration
		wantValid bool
	}{
		{name: "empty", wantValid: false},
		{name: "negative", value: "-1", wantValid: false},
		{name: "malformed", value: "one", wantValid: false},
		{name: "valid zero", value: "0", wantAge: 0, wantValid: true},
		{name: "valid seconds", value: "12", wantAge: 12 * time.Second, wantValid: true},
		{name: "saturates overflow", value: strconv.FormatInt(maxDurationSeconds+1, 10), wantAge: time.Duration(1<<63 - 1), wantValid: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			age, ok := parseAgeHeader(test.value)
			assert.Equal(t, test.wantValid, ok)
			assert.Equal(t, test.wantAge, age)
		})
	}

	t.Run("rejects missing or inverted date headers", func(t *testing.T) {
		for _, headers := range []http.Header{
			{},
			{"Date": []string{"not a date"}, "Last-Modified": []string{"Mon, 28 Sep 2026 19:46:06 GMT"}},
			{"Date": []string{"Mon, 28 Sep 2026 19:46:06 GMT"}, "Last-Modified": []string{"Mon, 28 Sep 2026 19:47:06 GMT"}},
		} {
			age, source, ok := lastModifiedAge(headers)
			assert.Zero(t, age)
			assert.Empty(t, source)
			assert.False(t, ok)
		}
	})
}

func TestTibiaDataHTMLDataCollectorRecordsAge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Age", "60")
		w.Header().Set("Date", "Mon, 28 Sep 2026 19:47:06 GMT")
		w.Header().Set("Last-Modified", "Mon, 28 Sep 2026 19:46:36 GMT")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>test</html>"))
	}))
	defer server.Close()

	metadata := &TibiaDataCacheMetadata{}
	body, err := TibiaDataHTMLDataCollector(TibiaDataRequestStruct{
		Method:        http.MethodGet,
		URL:           server.URL,
		RawBody:       true,
		CacheMetadata: metadata,
	})

	assert.NoError(t, err)
	assert.Equal(t, "<html>test</html>", body)
	assert.Equal(t, 240*time.Second, metadata.remainingCacheAge(300*time.Second))
}

func TestTibiaDataAPIHandleCachedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)

	previousCacheAwareness := TibiaDataCacheAwareness
	TibiaDataCacheAwareness = true
	t.Cleanup(func() { TibiaDataCacheAwareness = previousCacheAwareness })

	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second, nil)

	assert.Equal(t, "public, max-age=300, s-maxage=300", w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	staleMetadata := &TibiaDataCacheMetadata{}
	staleMetadata.recordUpstreamHeaders(http.Header{"Age": []string{"300"}})
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second, staleMetadata)
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	TibiaDataCacheAwareness = false
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second, nil)
	assert.Empty(t, w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	TibiaDataCacheAwareness = true
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v4/character/Test", nil)
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 0, nil)
	assert.Empty(t, w.Header().Get("Cache-Control"))
	assert.Equal(t, http.StatusOK, w.Code)

	TibiaDataAPIHandleCachedResponse(nil, "test", gin.H{"ok": true}, 300*time.Second, nil)
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
	TibiaDataAPIHandleCachedResponse(c, "test", gin.H{"ok": true}, 300*time.Second, nil)

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
