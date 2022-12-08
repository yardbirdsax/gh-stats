package pr

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yardbirdsax/gh-stats/internal/result"
)

func TestMyReviews(t *testing.T) {
  defer gock.Off()
  gock.Observe(gock.DumpRequest)
  expectedStartDateTime := time.Now().Add(-1 * 24 * time.Hour)
  expectedStartDateTimeString := expectedStartDateTime.Format("2006-01-02")
  expectedEndDateTime := time.Now()
  expectedEndDateTimeString := expectedEndDateTime.Format("2006-01-02")
  expectedParamEscaped := regexp.QuoteMeta(fmt.Sprintf("is:pr reviewed-by:@me created:%s..%s", expectedStartDateTimeString, expectedEndDateTimeString))
  gock.New("https://api.github.com").Get("search/issues").MatchParam("q", expectedParamEscaped).Reply(200).JSON(
    map[string]interface{}{
      "total_count": 2,
      "incomplete_results": false,
      "items": []map[string]interface{}{
        {
          "created_at": "2022-11-02T20:20:20Z",
        },
        {
          "created_at": "2022-11-01T20:20:20Z",
        },
      },
    },
  )
  expectedResults := &result.Results{
    {"createdat", "count"},
    {time.Date(2022,11,1, 0, 0, 0, 0, time.UTC),  1},
    {time.Date(2022,11,2, 0, 0, 0, 0, time.UTC), 1},
  }

  actualResults, err := MyReviews(expectedStartDateTime, expectedEndDateTime, "CreatedAt")

  require.NoError(t, err)
  assert.Equal(t, expectedResults, actualResults)
  assert.False(t, gock.IsPending(), "gock has pending requests")
  assert.False(t, gock.HasUnmatchedRequest(), "gock has unmatched requests")
}

func TestTeamReviews(t *testing.T) {
  defer gock.Off()
  gock.Observe(gock.DumpRequest)
  expectedOrgName := "org"
  expectedTeamname := "team"
  gock.New("https://api.github.com").Get(fmt.Sprintf("orgs/%s/teams/%s/members", expectedOrgName, expectedTeamname)).Reply(200).JSON(
    []map[string]interface{}{
      {
        "login": "user1",
      },
      {
        "login": "user2",
      },
    },
  )
  expectedStartDateTime := time.Now().Add(-1 * 24 * time.Hour)
  expectedStartDateTimeString := expectedStartDateTime.Format("2006-01-02")
  expectedEndDateTime := time.Now()
  expectedEndDateTimeString := expectedEndDateTime.Format("2006-01-02")
  expectedParamEscaped := regexp.QuoteMeta(fmt.Sprintf("is:pr reviewed-by:user1 reviewed-by:user2 user:org created:%s..%s", expectedStartDateTimeString, expectedEndDateTimeString))
  gock.New("https://api.github.com").Get("search/issues").MatchParam("q", expectedParamEscaped).MatchParam("page", "1").Reply(200).JSON(
    map[string]interface{}{
      "total_count": 4,
      "incomplete_results": false,
      "items": []map[string]interface{}{
        {
          "created_at": "2022-11-04T20:20:20Z",
        },
        {
          "created_at": "2022-11-03T20:20:20Z",
        },
      },
    },
  )
  gock.New("https://api.github.com").Get("search/issues").MatchParam("q", expectedParamEscaped).MatchParam("page", "2").Reply(200).JSON(
    map[string]interface{}{
      "total_count": 4,
      "incomplete_results": false,
      "items": []map[string]interface{}{
        {
          "created_at": "2022-11-02T20:20:20Z",
        },
        {
          "created_at": "2022-11-01T20:20:20Z",
        },
      },
    },
  )
  expectedResults := &result.Results{
    {"createdat", "count"},
    {time.Date(2022,11,1, 0, 0, 0, 0, time.UTC),  1},
    {time.Date(2022,11,2, 0, 0, 0, 0, time.UTC), 1},
    {time.Date(2022,11,3, 0, 0, 0, 0, time.UTC), 1},
    {time.Date(2022,11,4, 0, 0, 0, 0, time.UTC), 1},
  }

  actualResults, err := TeamReviews(expectedOrgName, expectedTeamname, expectedStartDateTime, expectedEndDateTime, "CreatedAt")

  require.NoError(t, err)
  assert.Equal(t, expectedResults, actualResults)
  assert.False(t, gock.IsPending(), "gock has pending requests")
  assert.False(t, gock.HasUnmatchedRequest(), "gock has unmatched requests")
}
