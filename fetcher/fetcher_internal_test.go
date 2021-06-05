package fetcher

import (
	"apptica/model"

	"encoding/json"
	"reflect"
	"testing"
)

func TestFetcher_prepareDataTest(t *testing.T) {
	dataJSON := `{"2":{"1":{"2021-05-20":34,"2021-05-21":38,"2021-05-22":34,"2021-05-23":34,"2021-05-24":34,"2021-05-25":37,"2021-05-26":36,"2021-05-27":37}},"23":{"3":{"2021-05-20":49,"2021-05-21":51,"2021-05-22":52,"2021-05-23":49,"2021-05-24":49,"2021-05-25":50,"2021-05-26":49,"2021-05-27":49},"1":{"2021-05-20":6,"2021-05-21":8,"2021-05-22":6,"2021-05-23":6,"2021-05-24":6,"2021-05-25":6,"2021-05-26":6,"2021-05-27":6}}}`
	data := make(map[model.Category]map[model.Subcategory]map[model.Date]model.Position)

	err := json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		t.Fatalf("failed to parse JSON: %s", err)
	}

	expected := map[model.Category]map[model.Date]model.Position{
		"2": {
			"2021-05-20": 34,
			"2021-05-21": 38,
			"2021-05-22": 34,
			"2021-05-23": 34,
			"2021-05-24": 34,
			"2021-05-25": 37,
			"2021-05-26": 36,
			"2021-05-27": 37,
		},
		"23": {
			"2021-05-20": 6,
			"2021-05-21": 8,
			"2021-05-22": 6,
			"2021-05-23": 6,
			"2021-05-24": 6,
			"2021-05-25": 6,
			"2021-05-26": 6,
			"2021-05-27": 6,
		},
	}

	actual := prepareData(data)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("failed to prepare data, %v != %v", actual, expected)
	}
}
