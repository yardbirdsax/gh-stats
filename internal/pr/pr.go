package pr

import (
  "time"

  "github.com/cli/go-gh"
)

func MyReviews() (map[time.Time]int, error) {
  var myReviews map[time.Time]int = map[time.Time]int{}

  client, err := gh.RESTClient(nil)
  if err != nil {
    return myReviews, err
  }

  response := issueSearchResponse{}
  err = client.Get("search/issues?q=is:pr+reviewed-by:@me", &response)
  if err != nil {
    return myReviews, err
  }

  for _, i := range response.Items {
    roundedDate := i.CreatedAt.Truncate(24 * time.Hour)
    myReviews[roundedDate] += 1
  }

  return myReviews, nil
}

type issueSearchResponse struct {
  TotalCount int `json:"total_count"`
  IncompleteResults bool `json:"incomplete_results"`
  Items []issueSearchItem `json:"items"`
}

type issueSearchItem struct {
  CreatedAt time.Time `json:"created_at"`
}
