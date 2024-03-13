// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: realmSlicieManagement.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const cloneRecordInConfigBySliceID = `-- name: CloneRecordInConfigBySliceID :execresult
INSERT INTO
    config (
        realm, slice, name, descr, val, ver, setby
    )
SELECT realm, $2, name, descr, val, ver, $3
FROM config
WHERE
    config.slice = $1
`

type CloneRecordInConfigBySliceIDParams struct {
	Slice   int32  `json:"slice"`
	Slice_2 int32  `json:"slice_2"`
	Setby   string `json:"setby"`
}

func (q *Queries) CloneRecordInConfigBySliceID(ctx context.Context, arg CloneRecordInConfigBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInConfigBySliceID, arg.Slice, arg.Slice_2, arg.Setby)
}

const cloneRecordInRealmSliceBySliceID = `-- name: CloneRecordInRealmSliceBySliceID :one
INSERT INTO
    realmslice (
        realm, descr, active, activateat, deactivateat
    )
SELECT
    realm,
    COALESCE(descr, $3::text),
    true,
    activateat,
    deactivateat
FROM realmslice
WHERE
    realmslice.id = $1
    AND realmslice.realm = $2
RETURNING
    realmslice.id
`

type CloneRecordInRealmSliceBySliceIDParams struct {
	ID    int32       `json:"id"`
	Realm string      `json:"realm"`
	Descr pgtype.Text `json:"descr"`
}

func (q *Queries) CloneRecordInRealmSliceBySliceID(ctx context.Context, arg CloneRecordInRealmSliceBySliceIDParams) (int32, error) {
	row := q.db.QueryRow(ctx, cloneRecordInRealmSliceBySliceID, arg.ID, arg.Realm, arg.Descr)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const cloneRecordInRulesetBySliceID = `-- name: CloneRecordInRulesetBySliceID :execresult
INSERT INTO
    ruleset (
        realm, slice, app, brwf, class, setname, schemaid, is_active, is_internal, ruleset, createdby
    )
SELECT
    realm,
    $2,
    app,
    brwf,
    class,
    setname,
    schemaid,
    is_active,
    is_internal,
    ruleset,
    $3
FROM ruleset
WHERE
    ruleset.slice = $1
    AND (
        $4::text [] is null
        OR app = any ($4::text [])
    )
`

type CloneRecordInRulesetBySliceIDParams struct {
	Slice     int32    `json:"slice"`
	Slice_2   int32    `json:"slice_2"`
	Createdby string   `json:"createdby"`
	App       []string `json:"app"`
}

func (q *Queries) CloneRecordInRulesetBySliceID(ctx context.Context, arg CloneRecordInRulesetBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInRulesetBySliceID,
		arg.Slice,
		arg.Slice_2,
		arg.Createdby,
		arg.App,
	)
}

const cloneRecordInSchemaBySliceID = `-- name: CloneRecordInSchemaBySliceID :execresult
INSERT INTO
    schema (
        realm, slice, app, brwf, class, patternschema, actionschema, createdby
    )
SELECT
    realm,
    $2,
    app,
    brwf,
    class,
    patternschema,
    actionschema,
    $3
FROM schema
WHERE
    schema.slice = $1
    AND (
        $4::text [] is null
        OR app = any ($4::text [])
    )
`

type CloneRecordInSchemaBySliceIDParams struct {
	Slice     int32    `json:"slice"`
	Slice_2   int32    `json:"slice_2"`
	Createdby string   `json:"createdby"`
	App       []string `json:"app"`
}

func (q *Queries) CloneRecordInSchemaBySliceID(ctx context.Context, arg CloneRecordInSchemaBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInSchemaBySliceID,
		arg.Slice,
		arg.Slice_2,
		arg.Createdby,
		arg.App,
	)
}

const insertNewRecordInRealmSlice = `-- name: InsertNewRecordInRealmSlice :one
INSERT INTO
    realmslice (
        realm, descr, active
    )
VALUES (
    $1, $2, true
)
RETURNING
   realmslice.id
`

type InsertNewRecordInRealmSliceParams struct {
	Realm string `json:"realm"`
	Descr string `json:"descr"`
}

func (q *Queries) InsertNewRecordInRealmSlice(ctx context.Context, arg InsertNewRecordInRealmSliceParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertNewRecordInRealmSlice, arg.Realm, arg.Descr)
	var id int32
	err := row.Scan(&id)
	return id, err
}