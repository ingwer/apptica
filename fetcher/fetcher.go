package fetcher

import (
	"apptica/model"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Fetcher interface {
	FetchTopPositions(ctx context.Context, appID, countryCode int, dateFrom, dateTo model.Date) (map[model.Category]map[model.Date]model.Position, error)
}

type fetcher struct {
	endpoint   string
	httpClient *http.Client
	token      string
}

type apiResponse struct {
	StatusCode int                                                                    `json:"status_code"`
	Message    string                                                                 `json:"message"`
	Data       map[model.Category]map[model.Subcategory]map[model.Date]model.Position `json:"data"`
}

var (
	errFailedToDoHTTPRequest = errors.New("failed to do http request")
	errBadResponseCode       = errors.New("response status code is not Ok")
	errFailedToReadResponse  = errors.New("failed to read response body")
	errBadJSONResponse       = errors.New("failed to unmarshal json")
)

func (f fetcher) FetchTopPositions(ctx context.Context, appID, countryCode int, dateFrom, dateTo model.Date) (map[model.Category]map[model.Date]model.Position, error) {
	url := fmt.Sprintf(
		"%s/%d/%d?date_from=%s&date_to=%s&%s",
		f.endpoint,
		appID,
		countryCode,
		dateFrom,
		dateTo,
		f.token,
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errFailedToDoHTTPRequest
	}

	response, err := f.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, errFailedToDoHTTPRequest
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errFailedToReadResponse
	}

	if response.StatusCode != http.StatusOK {
		return nil, errBadResponseCode
	}

	r := &apiResponse{}
	err = json.Unmarshal(responseBody, &r)

	if err != nil {
		return nil, errBadJSONResponse
	}

	if r.StatusCode != http.StatusOK {
		return nil, errBadResponseCode
	}

	return prepareData(r.Data), nil
}

func prepareData(data map[model.Category]map[model.Subcategory]map[model.Date]model.Position) map[model.Category]map[model.Date]model.Position {
	result := make(map[model.Category]map[model.Date]model.Position)

	for category, subcatStat := range data {
		for _, dateToPosition := range subcatStat {
			for date, position := range dateToPosition {
				if stat, ok := result[category]; ok {
					if currentPosition, ok := stat[date]; (ok && position < currentPosition) || !ok {
						stat[date] = position
					}
				} else {
					result[category] = make(map[model.Date]model.Position)
					result[category][date] = position
				}
			}
		}
	}

	return result
}

type Option func(*fetcher)

func WithHTTPClient(client *http.Client) Option {
	return func(fetcher *fetcher) {
		fetcher.httpClient = client
	}
}

func WithEndpoint(endpoint string) Option {
	return func(fetcher *fetcher) {
		fetcher.endpoint = endpoint
	}
}

func WithToken(token string) Option {
	return func(fetcher *fetcher) {
		fetcher.token = token
	}
}

func New(options ...Option) Fetcher {
	f := &fetcher{}

	for _, option := range options {
		option(f)
	}

	return f
}
