package app

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"log"
)

type App struct {
	Client *bigquery.Client
}

func NewApp(client *bigquery.Client) *App {
	return &App{
		Client: client,
	}
}

func (a *App) Run(ctx context.Context, query string) (*Result, error) {
	var err error

	q := a.Client.Query(query)

	var job *bigquery.Job
	if job, err = q.Run(ctx); err != nil {
		return nil, errors.WithStack(err)
	}

	var jobStatus *bigquery.JobStatus
	if jobStatus, err = job.Wait(ctx); err != nil {
		log.Printf("%+v", jobStatus)
		return nil, errors.WithStack(err)
	}

	if jobStatus.Err() != nil {
		return &Result{
			RowCount:  0,
			JobStatus: jobStatus,
		}, nil
	}

	var it *bigquery.RowIterator
	if it, err = job.Read(ctx); err != nil {
		return nil, errors.WithStack(err)
	}

	var values []bigquery.Value

	for {
		var v []bigquery.Value
		err = it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != iterator.Done && err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, v)
	}

	return &Result{
		RowCount:  len(values),
		JobStatus: jobStatus,
		Schema:    it.Schema,
		ValueList: values,
	}, nil
}
