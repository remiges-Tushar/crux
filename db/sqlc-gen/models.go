// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"database/sql"
	"encoding/json"
	"time"
)

type App struct {
	ID          int32     `json:"id"`
	Realm       string    `json:"realm"`
	Shortname   string    `json:"shortname"`
	Shortnamelc string    `json:"shortnamelc"`
	Longname    string    `json:"longname"`
	Setby       string    `json:"setby"`
	Setat       time.Time `json:"setat"`
}

type Capgrant struct {
	ID        int32          `json:"id"`
	Realm     int32          `json:"realm"`
	User      string         `json:"user"`
	App       sql.NullString `json:"app"`
	Cap       string         `json:"cap"`
	From      sql.NullTime   `json:"from"`
	To        sql.NullTime   `json:"to"`
	Setat     time.Time      `json:"setat"`
	Setby     string         `json:"setby"`
	Isdeleted sql.NullBool   `json:"isdeleted"`
}

type Config struct {
	Realm int32          `json:"realm"`
	Slice int32          `json:"slice"`
	Name  string         `json:"name"`
	Descr string         `json:"descr"`
	Val   sql.NullString `json:"val"`
	Ver   sql.NullInt32  `json:"ver"`
	Setby string         `json:"setby"`
	Setat time.Time      `json:"setat"`
}

type Deactivated struct {
	ID      int32          `json:"id"`
	Realm   string         `json:"realm"`
	User    sql.NullString `json:"user"`
	Deactby string         `json:"deactby"`
	Deactat time.Time      `json:"deactat"`
}

type Realm struct {
	ID          int32           `json:"id"`
	Shortname   string          `json:"shortname"`
	Shortnamelc string          `json:"shortnamelc"`
	Longname    string          `json:"longname"`
	Setby       string          `json:"setby"`
	Setat       time.Time       `json:"setat"`
	Payload     json.RawMessage `json:"payload"`
}

type Realmslice struct {
	ID           int32        `json:"id"`
	Realm        string       `json:"realm"`
	Descr        string       `json:"descr"`
	Active       bool         `json:"active"`
	Activateat   sql.NullTime `json:"activateat"`
	Deactivateat sql.NullTime `json:"deactivateat"`
}

type Ruleset struct {
	ID         int32           `json:"id"`
	Realm      int32           `json:"realm"`
	Slice      int32           `json:"slice"`
	App        string          `json:"app"`
	Brwf       string          `json:"brwf"`
	Class      string          `json:"class"`
	Setname    string          `json:"setname"`
	Schemaid   int32           `json:"schemaid"`
	IsActive   sql.NullBool    `json:"is_active"`
	IsInternal bool            `json:"is_internal"`
	Ruleset    json.RawMessage `json:"ruleset"`
	Createdat  time.Time       `json:"createdat"`
	Createdby  string          `json:"createdby"`
	Editedat   time.Time       `json:"editedat"`
	Editedby   string          `json:"editedby"`
}

type Schema struct {
	ID            int32           `json:"id"`
	Realm         int32           `json:"realm"`
	Slice         int32           `json:"slice"`
	App           string          `json:"app"`
	Brwf          string          `json:"brwf"`
	Class         string          `json:"class"`
	Patternschema json.RawMessage `json:"patternschema"`
	Actionschema  json.RawMessage `json:"actionschema"`
	Createdat     time.Time       `json:"createdat"`
	Createdby     string          `json:"createdby"`
	Editedat      time.Time       `json:"editedat"`
	Editedby      string          `json:"editedby"`
}

type Stepworkflow struct {
	Slice    int32          `json:"slice"`
	App      sql.NullString `json:"app"`
	Step     string         `json:"step"`
	Workflow string         `json:"workflow"`
}

type Wfinstance struct {
	ID       int32         `json:"id"`
	Entityid int32         `json:"entityid"`
	Slice    int32         `json:"slice"`
	App      string        `json:"app"`
	Class    string        `json:"class"`
	Workflow string        `json:"workflow"`
	Step     string        `json:"step"`
	Loggedat time.Time     `json:"loggedat"`
	Doneat   sql.NullTime  `json:"doneat"`
	Nextstep string        `json:"nextstep"`
	Parent   sql.NullInt32 `json:"parent"`
}