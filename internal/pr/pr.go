package pr

import (
	"sort"
	"time"

	"github.com/cli/go-gh"
	"github.com/yardbirdsax/gh-stats/internal/result"
)

func MyReviews() (*result.Results, error) {
	var myReviews *result.Results

	client, err := gh.RESTClient(nil)
	if err != nil {
		return myReviews, err
	}

	response := issueSearchResponse{}
	err = client.Get("search/issues?q=is:pr+reviewed-by:@me", &response)
	if err != nil {
		return myReviews, err
	}

  columnNames := []interface{}{
    "created date",
    "count",
  }

  mapData := make(map[time.Time]int)
	for _, i := range response.Items {
		roundedDate := i.CreatedAt.Truncate(24 * time.Hour)
    mapData[roundedDate] += 1
	}

  data := make([][]interface{}, 0, len(mapData))
  for k, v := range mapData {
    data = append(data, []interface{}{k, v})
  }
  sort.Slice(data, func(i, j int) bool {
    firstValue := data[i][0].(time.Time)
    secondValue := data[j][0].(time.Time)
    return firstValue.Before(secondValue)
  })

	myReviews, err = result.FromSlice(columnNames, data)
  return myReviews, err
}

type issueSearchResponse struct {
	TotalCount        int               `json:"total_count"`
	IncompleteResults bool              `json:"incomplete_results"`
	Items             []issueSearchItem `json:"items"`
}

type issueSearchItem struct {
	CreatedAt time.Time `json:"created_at"`
}
