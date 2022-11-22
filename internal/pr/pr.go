package pr

import (
	"fmt"
  //"net/url"
	"sort"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/yardbirdsax/gh-stats/internal/result"
)

func MyReviews() (*result.Results, error) {
	filter := "is:pr reviewed-by:@me"
  return getIssueCount(filter)
}

func TeamReviews(orgName string, teamName string) (*result.Results, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	teamResponse := &teamMemberResponse{}
	err = client.Get(fmt.Sprintf("/orgs/%s/teams/%s/members", orgName, teamName), teamResponse)
	if err != nil {
		return nil, err
	}
	teamMembers := []string{}
	for _, member := range *teamResponse {
		teamMembers = append(teamMembers, fmt.Sprintf("reviewed-by:%s", member.Login))
	}

	filter := fmt.Sprintf("is:pr (%s) user:%s", strings.Join(teamMembers, " "), orgName)
	return getIssueCount(filter)
}

func getClient() (api.RESTClient, error) {
	token, _ := auth.TokenForHost("github.com")
	return gh.RESTClient(&api.ClientOptions{
		AuthToken: token,
	})
}

func getIssueCount(filter string) (*result.Results, error) {
	var results *result.Results

	client, err := getClient()
	if err != nil {
		return results, err
	}

	response := issueSearchResponse{}
	err = client.Get(fmt.Sprintf("search/issues?q=%s", filter), &response)
	if err != nil {
		return results, err
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

	results, err = result.FromSlice(columnNames, data)
	return results, err
}

type issueSearchResponse struct {
	TotalCount        int               `json:"total_count"`
	IncompleteResults bool              `json:"incomplete_results"`
	Items             []issueSearchItem `json:"items"`
}

type issueSearchItem struct {
	CreatedAt time.Time `json:"created_at"`
}

type teamMemberResponse []teamMember

type teamMember struct {
	Login string `json:"login"`
}
