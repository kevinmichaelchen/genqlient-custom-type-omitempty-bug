package main

import (
	"context"
	"fmt"
	khan "github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"
	"github.com/kevinmichaelchen/genqlient-custom-type-omitempty-bug/internal/client/graphql"
	"github.com/kevinmichaelchen/genqlient-custom-type-omitempty-bug/pkg/postgres"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	d := doer{c: http.DefaultClient}
	client := khan.NewClient("http://localhost:8080/v1/graphql", d)

	id := uuid.New().String()
	now := time.Now()
	name := "Kevin"

	_, err := graphql.CreatePerson(ctx, client, &graphql.PersonInsertInput{
		Id:   &id,
		Name: &name,
		ValidTimeRange: &postgres.TimeRange{
			Start: now,
			End:   now.Add(time.Hour),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success!")
}

type doer struct {
	c *http.Client
}

func (d doer) Do(req *http.Request) (*http.Response, error) {
	//req.Header.Set("x-hasura-admin-secret", "myadminsecretkey")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwczovL2hhc3VyYS5pby9qd3QvY2xhaW1zIjp7IngtaGFzdXJhLWFsbG93ZWQtcm9sZXMiOlsic2VydmljZSJdLCJ4LWhhc3VyYS1kZWZhdWx0LXJvbGUiOiJzZXJ2aWNlIn0sImlhdCI6MTY4MTkyMDg4OSwic3ViIjoic2VydmljZSJ9.TEjxWgfMvJXBFaxr0qxJnPKrrkfnO-6m80moTNliLPE")
	return d.c.Do(req)
}
