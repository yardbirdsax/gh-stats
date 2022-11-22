package pr

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/yardbirdsax/gh-stats/internal/result"
)

type groupByField string

const (
	groupByFieldCreatedAt groupByField = "CreatedAt"
)

func MyReviews(startDate time.Time, endDate time.Time, groupByField string) (*result.Results, error) {
	filter := fmt.Sprintf("is:pr reviewed-by:@me created:%s..%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
  return getIssueCount(filter, groupByField)
}

func TeamReviews(orgName string, teamName string, startDate time.Time, endDate time.Time, groupByField string) (*result.Results, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	teamResponse := &teamMemberResponse{}
	err = client.Get(fmt.Sprintf("orgs/%s/teams/%s/members", orgName, teamName), teamResponse)
	if err != nil {
		return nil, err
	}
	teamMembers := []string{}
	for _, member := range *teamResponse {
		teamMembers = append(teamMembers, fmt.Sprintf("reviewed-by:%s", member.Login))
	}

	filter := fmt.Sprintf("is:pr %s user:%s created:%s..%s", strings.Join(teamMembers, " "), orgName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return getIssueCount(filter, groupByField)
}

func getClient() (api.RESTClient, error) {
	token, _ := auth.TokenForHost("github.com")
	return gh.RESTClient(&api.ClientOptions{
		AuthToken: token,
	})
}

func getIssueCount(filter string, groupByField string) (*result.Results, error) {
	var results *result.Results

	client, err := getClient()
	if err != nil {
		return results, err
	}

	responses := []issueSearchResponse{}
	totalItemCount := 0
	page := 1
	sanitizedFilter := url.QueryEscape(filter)
	for {
		response := issueSearchResponse{}
		err = client.Get(fmt.Sprintf("search/issues?q=%s&page=%d&per_page=100", sanitizedFilter, page), &response)
		if err != nil {
			return results, err
		}
		responses = append(responses, response)
		totalItemCount += len(response.Items)
		if totalItemCount >= response.TotalCount {
			break
		}
		page++
	}

	columnNames := []interface{}{
		strings.ToLower(groupByField),
		"count",
	}

	mapData := make(map[interface{}]int)
	for _, response := range responses {
		for _, i := range response.Items {
			field := reflect.ValueOf(i).FieldByName(groupByField)
			if field == (reflect.Value{}) {
				return results, fmt.Errorf("field with name '%s' does not exist", groupByField)
			}
			switch field.Type().String() {
			case "string":
				mapData[field.String()] += 1
			case "time.Time":
				roundedDate := field.Interface().(time.Time).Truncate(24 * time.Hour)
				mapData[roundedDate] += 1
			default:
				return results, fmt.Errorf("type of group by field (%s) is not currently supported", field.Type().String())
			}
		}
	}

	data := make([][]interface{}, 0, len(mapData))
	for k, v := range mapData {
		data = append(data, []interface{}{k, v})
	}
	sort.Slice(data, func(i, j int) bool {
		firstValue := data[i][0]
		secondValue := data[j][0]
		switch reflect.TypeOf(firstValue).Name() {
		case "string":
			firstValueString := firstValue.(string)
			secondValueString := secondValue.(string)
			return firstValueString < secondValueString
		case "Time":
			firstValueTime := firstValue.(time.Time)
			secondValueTime := secondValue.(time.Time)
			return firstValueTime.Before(secondValueTime)
		default:
			return false
		}
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
	RepositoryURL string `json:"repository_url"`
}

type teamMemberResponse []teamMember

type teamMember struct {
	Login string `json:"login"`
}
