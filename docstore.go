package docstore

import (
	"context"
	"fmt"
	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/awsdynamodb"
	"net/url"
	"strconv"
)

func OpenCollection(ctx context.Context, uri string) (*docstore.Collection, error) {

	var col *docstore.Collection

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse '%s', %w", uri, err)
	}

	if u.Scheme == "awsdynamodb" {

		table := u.Host

		q := u.Query()

		partition_key := q.Get("partition-key")
		region := q.Get("region")
		local := q.Get("local")
		credentials := q.Get("credentials")
		q_allow_scans := q.Get("allow-scans")

		cl_uri := fmt.Sprintf("dynamodb://?region=%s&credentials=%s&local=%s", region, credentials, local)

		cl, err := aa_dynamodb.NewClientWithURI(ctx, cl_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create DynamoDB client, %w", err)
		}

		col_opts := &awsdynamodb.Options{}

		if q_allow_scans != "" {

			allow, err := strconv.ParseBool(q_allow_scans)

			if err != nil {
				return nil, fmt.Errorf("Failed to parse ?allow-scans= parameter, %w", err)
			}

			col_opts.AllowScans = allow
		}

		c, err := awsdynamodb.OpenCollection(cl, table, partition_key, "", col_opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to open collection, %w", err)
		}

		col = c

	} else {

		c, err := docstore.OpenCollection(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to open collection, %w", err)
		}

		col = c
	}

	return col, nil
}
