package result

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsMarkdownTable(t *testing.T) {
  results := &Result{}
  tablularResults := tabularResult{
    {
      "id",
      "user",
      "created_date",
    },
    {
      1,
      "someone",
      time.Date(2022, 11, 07, 0, 0, 0, 0, time.UTC),
    },
    {
      2,
      "someone-else",
      time.Date(2022, 11, 8, 0, 0, 0, 0, time.UTC),
    },
  }
  results.tabularResults = tablularResults
  expectedOuput := `| id | user | created_date |
| ---:| --- | --- |
| 1 | someone | 2022-11-07 00:00:00 +0000 UTC |
| 2 | someone-else | 2022-11-08 00:00:00 +0000 UTC |`
  actualOutput := results.AsMarkdownTable()

  assert.Equal(t, expectedOuput, actualOutput)
}