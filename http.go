package sportmonks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL         = "https://api.sportmonks.com/v3"
	bookmakersURI          = "/odds/bookmakers"
	coachesURI             = "/football/coaches"
	commentariesFixtureURI = "/commentaries/fixture"
	continentsURI          = "/core/continents"
	countriesURI           = "/core/countries"
	fixturesURI            = "/football/fixtures"
	fixturesDateURI        = "/football/fixtures/date"
	fixturesBetweenURI     = "/football/fixtures/between"
	fixturesMultiURI       = "/football/fixtures/multi"
	headToHeadURI          = "/football/fixtures/head-to-head"
	leaguesURI             = "/football/leagues"
	oddsFixtureURI         = "/football/odds/pre-match/fixtures"
	playersURI             = "/football/players"
	roundsURI              = "/football/rounds"
	roundsSeasonURI        = "/football/rounds/seasons"
	seasonsURI             = "/football/seasons"
	stagesURI              = "/football/stages"
	stagesSeasonURI        = "/football/stages/seasons"
	teamSquadURI           = "/football/squads/teams"
	teamSeasonSquadURI     = "/football/squads/seasons"
	teamsURI               = "/football/teams"
	teamsSeasonURI         = "/football/teams/seasons"
	topScorersSeasonURI    = "/football/topscorers/seasons"
	tvStationsURI          = "/football/tv-stations/fixtures"
	venuesURI              = "/football/venues"
	venuesSeasonURI        = "/football/venues/seasons"
)

// HTTPClient is a HTTP request builder and sender.
type HTTPClient struct {
	HTTPClient *http.Client
	BaseURL    string
	Key        string
}

// NewDefaultHTTPClient creates a new Client with default settings. A key is required to instantiate the Client.
func NewDefaultHTTPClient(key string) *HTTPClient {
	return &HTTPClient{
		HTTPClient: &http.Client{},
		BaseURL:    defaultBaseURL,
		Key:        key,
	}
}

// SetHTTPClient provides functionality to override the default HTTPClient property.
func (c *HTTPClient) SetHTTPClient(h *http.Client) {
	if h != nil {
		c.HTTPClient = h
	}
}

// SetBaseURL provides functionality to override the default BaseURL property.
func (c *HTTPClient) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *HTTPClient) getResource(ctx context.Context, url string, query url.Values, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+url, nil)

	if err != nil {
		return err
	}

	query.Set("api_token", c.Key)

	req.URL.RawQuery = query.Encode()

	return c.do(ctx, req, response)
}

func (c *HTTPClient) do(ctx context.Context, req *http.Request, intf interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = checkStatusCode(resp); err != nil {
		return err
	}

	return parseJSONResponseBody(resp.Body, intf)
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		err := new(ErrRateLimit)

		e := parseJSONResponseBody(resp.Body, &err)

		if e != nil {
			return e
		}

		return err
	}

	if resp.StatusCode != http.StatusOK {
		err := new(ErrBadStatusCode)

		e := parseJSONResponseBody(resp.Body, &err)

		if e != nil {
			return e
		}

		return err
	}

	return nil
}

func parseJSONResponseBody(body io.ReadCloser, response interface{}) error {
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return err
	}

	return nil
}

func formatFilters(query *url.Values, filters map[string][]int) {
	for k, v := range filters {
		query.Set(k, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]"))
	}
}
