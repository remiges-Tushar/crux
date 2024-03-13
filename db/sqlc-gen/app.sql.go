// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: app.sql

package sqlc

import (
	"context"
)

const appNew = `-- name: AppNew :many
INSERT INTO
    app (
        realm, shortname, shortnamelc, longname, setby
    )
VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING id, realm, shortname, shortnamelc, longname, setby, setat
`

type AppNewParams struct {
	Realm       string `json:"realm"`
	Shortname   string `json:"shortname"`
	Shortnamelc string `json:"shortnamelc"`
	Longname    string `json:"longname"`
	Setby       string `json:"setby"`
}

func (q *Queries) AppNew(ctx context.Context, arg AppNewParams) ([]App, error) {
	rows, err := q.db.Query(ctx, appNew,
		arg.Realm,
		arg.Shortname,
		arg.Shortnamelc,
		arg.Longname,
		arg.Setby,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []App
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.Realm,
			&i.Shortname,
			&i.Shortnamelc,
			&i.Longname,
			&i.Setby,
			&i.Setat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAppName = `-- name: GetAppName :one
select count(1) FROM app WHERE shortnamelc = $1 AND realm = $2
`

type GetAppNameParams struct {
	Shortnamelc string `json:"shortnamelc"`
	Realm       string `json:"realm"`
}

func (q *Queries) GetAppName(ctx context.Context, arg GetAppNameParams) (int64, error) {
	row := q.db.QueryRow(ctx, getAppName, arg.Shortnamelc, arg.Realm)
	var count int64
	err := row.Scan(&count)
	return count, err
}
