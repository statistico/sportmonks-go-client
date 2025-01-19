package sportmonks

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type Fixture struct {
	ID                  int            `json:"id"`
	SportID             int            `json:"sport_id"`
	LeagueID            int            `json:"league_id"`
	SeasonID            int            `json:"season_id"`
	StageID             int            `json:"stage_id"`
	GroupID             *int           `json:"group_id"`
	AggregateID         *int           `json:"aggregate_id"`
	RoundID             int            `json:"round_id"`
	StateID             int            `json:"state_id"`
	VenueID             *int           `json:"venue_id"`
	Name                string         `json:"name"`
	StartingAt          string         `json:"starting_at"`
	ResultInfo          string         `json:"result_info"`
	Leg                 string         `json:"leg"`
	Details             *string        `json:"details"`
	Length              int            `json:"length"`
	Placeholder         bool           `json:"placeholder"`
	HasOdds             bool           `json:"has_odds"`
	HasPremiumOdds      bool           `json:"has_premium_odds"`
	StartingAtTimestamp int64          `json:"starting_at_timestamp"`
	Round               *Round         `json:"round,omitempty"`
	Stage               *Stage         `json:"stage,omitempty"`
	League              *League        `json:"league,omitempty"`
	Season              *Season        `json:"season,omitempty"`
	Coaches             []Coach        `json:"coaches,omitempty"`
	Venues              *Venue         `json:"venues,omitempty"`
	FixtureState        *FixtureState  `json:"state,omitempty"`
	WeatherReport       *WeatherReport `json:"weatherReport,omitempty"`
	Lineups             []LineupPlayer `json:"lineups"`
	Events              []FixtureEvent `json:"events,omitempty"`
	Statistics          []FixtureStat  `json:"statistics,omitempty"`
	Scores              []Score        `json:"scores,omitempty"`
	Formations          []Formation    `json:"formations,omitempty"`
	Participants        []Team         `json:"participants,omitempty"`
}

// FixtureByID fetches a Fixture resource by ID. Use the includes slice of string to enrich the response data.
func (c *HTTPClient) FixtureByID(ctx context.Context, id int, includes []string, filters map[string][]int) (*Fixture, *ResponseDetails, error) {
	path := fmt.Sprintf(fixturesURI+"/%d", id)

	values := url.Values{
		"include": {strings.Join(includes, ";")},
	}

	formatFilters(&values, filters)

	response := struct {
		Data         *Fixture       `json:"data"`
		Subscription []Subscription `json:"subscription"`
		RateLimit    RateLimit      `json:"rate_limit"`
		TimeZone     string         `json:"timezone"`
	}{}

	err := c.getResource(ctx, path, values, &response)

	if err != nil {
		return nil, nil, err
	}

	return response.Data, &ResponseDetails{
		Subscription: response.Subscription,
		RateLimit:    response.RateLimit,
		TimeZone:     response.TimeZone,
	}, err
}

// FixturesByID fetches multiple Fixture resources by their IDS. Use the includes slice of string to enrich the response data.
func (c *HTTPClient) FixturesByID(ctx context.Context, ids []int, includes []string, filters map[string][]int) ([]Fixture, *ResponseDetails, error) {
	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	path := fmt.Sprintf(fixturesMultiURI+"/%s", str)

	return multipleFixtureResponse(ctx, c, path, includes, filters)
}

// FixturesByDate fetches multiple Fixture resources for a given date. Use the includes slice of string to enrich the response data.
func (c *HTTPClient) FixturesByDate(ctx context.Context, date time.Time, includes []string, filters map[string][]int) ([]Fixture, *ResponseDetails, error) {
	path := fmt.Sprintf(fixturesDateURI + "/" + date.Format(dateFormat))

	return multipleFixtureResponse(ctx, c, path, includes, filters)
}

// FixturesBetween fetches multiple Fixture resources for between two dates. Use the includes slice of string to enrich the response data.
func (c *HTTPClient) FixturesBetween(ctx context.Context, from, to time.Time, includes []string, filters map[string][]int) ([]Fixture, *ResponseDetails, error) {
	path := fmt.Sprintf(fixturesBetweenURI+"/%s/%s", from.Format(dateFormat), to.Format(dateFormat))

	return multipleFixtureResponse(ctx, c, path, includes, filters)
}

// FixturesBetweenForTeam fetches multiple Fixture resources for between two dates for a given team ID. Use the includes slice of string
// to enrich the response data.
func (c *HTTPClient) FixturesBetweenForTeam(ctx context.Context, from, to time.Time, teamID int, includes []string, filters map[string][]int) ([]Fixture, *ResponseDetails, error) {
	path := fmt.Sprintf(fixturesBetweenURI+"/%s/%s/%d", from.Format(dateFormat), to.Format(dateFormat), teamID)

	return multipleFixtureResponse(ctx, c, path, includes, filters)
}

// HeadToHead fetches multiple Fixture resources of results between two teams. Use the includes slice of string to enrich
// the response data.
func (c *HTTPClient) HeadToHead(ctx context.Context, idOne, idTwo int, includes []string) ([]Fixture, *ResponseDetails, error) {
	path := fmt.Sprintf(headToHeadURI+"/%d/%d", idOne, idTwo)

	return multipleFixtureResponse(ctx, c, path, includes, map[string][]int{})
}

func multipleFixtureResponse(ctx context.Context, client *HTTPClient, path string, includes []string, filters map[string][]int) ([]Fixture, *ResponseDetails, error) {
	values := url.Values{
		"include": {strings.Join(includes, ";")},
	}

	formatFilters(&values, filters)

	response := struct {
		Data         []Fixture      `json:"data"`
		Pagination   *Pagination    `json:"pagination"`
		Subscription []Subscription `json:"subscription"`
		RateLimit    RateLimit      `json:"rate_limit"`
		TimeZone     string         `json:"timezone"`
	}{}

	err := client.getResource(ctx, path, values, &response)

	if err != nil {
		return nil, nil, err
	}

	return response.Data, &ResponseDetails{
		Pagination:   response.Pagination,
		Subscription: response.Subscription,
		RateLimit:    response.RateLimit,
		TimeZone:     response.TimeZone,
	}, err
}
