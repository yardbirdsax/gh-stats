package pr

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyReviews(t *testing.T) {
  defer gock.Off()
  gock.Observe(gock.DumpRequest)
  gock.New("https://api.github.com").Get("search/issues").MatchParam("q", "is:pr reviewed-by:@me").Reply(200).JSON(
    map[string]interface{}{
      "total_count": 10,
      "incomplete_results": false,
      "items": []map[string]interface{}{
        {
          "created_at": "2022-11-01T20:20:20Z",
        },
        {
          "created_at": "2022-11-02T20:20:20Z",
        },
      },
    },
  )
  expectedResults := map[time.Time]int{
    time.Date(2022,11,1, 0, 0, 0, 0, time.UTC): 1,
    time.Date(2022,11,2, 0, 0, 0, 0, time.UTC): 1,
  }

  actualResults, err := MyReviews()

  require.NoError(t, err)
  assert.Equal(t, expectedResults, actualResults)
  assert.False(t, gock.IsPending(), "gock has pending requests")
  assert.False(t, gock.HasUnmatchedRequest(), "gock has unmatched requests")
}