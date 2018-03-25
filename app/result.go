package app

import (
	"cloud.google.com/go/bigquery"
)

type Result struct {
	RowCount  int
	JobStatus *bigquery.JobStatus
	Schema    bigquery.Schema
	ValueList []bigquery.Value
}

func (r *Result) IsOK() bool {
	return r.RowCount == 0 && r.JobStatus.Err() == nil
}
