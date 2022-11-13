// Package result deals with generic ways to return and format the results of API calls.
package result

import (
  "github.com/jedib0t/go-pretty/v6/table"
)

// This is a set of results that can be represented as rows of data.
type tabularResult [][]interface{}

type Result struct {
  tabularResults tabularResult
}

// AsMarkdownTable renders tabular results as a Markdown table.
func (r *Result) AsMarkdownTable() string {
  var renderedString string
  t := table.NewWriter()

  // Construct the table header
  firstRow := r.tabularResults[0]
  tableHeaderRow := table.Row{}
  for _, k := range firstRow {
    tableHeaderRow = append(tableHeaderRow, k)
  }
  t.AppendHeader(tableHeaderRow)

  for _, result := range r.tabularResults[1:] {
    tableRow := table.Row{}
    for _, colVal := range result {
      tableRow = append(tableRow, colVal)
    }
    t.AppendRow(tableRow)
  }

  renderedString = t.RenderMarkdown()
  return renderedString
}
