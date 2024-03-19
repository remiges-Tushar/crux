// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	ActivateRecord(ctx context.Context, arg ActivateRecordParams) error
	AddWFNewInstances(ctx context.Context, arg AddWFNewInstancesParams) ([]Wfinstance, error)
	AllRuleset(ctx context.Context) ([]Ruleset, error)
	AllSchemas(ctx context.Context) ([]Schema, error)
	AppDelete(ctx context.Context, arg AppDeleteParams) error
	AppExist(ctx context.Context, app string) (int32, error)
	AppNew(ctx context.Context, arg AppNewParams) ([]App, error)
	AppUpdate(ctx context.Context, arg AppUpdateParams) error
	CloneRecordInConfigBySliceID(ctx context.Context, arg CloneRecordInConfigBySliceIDParams) (pgconn.CommandTag, error)
	CloneRecordInRealmSliceBySliceID(ctx context.Context, arg CloneRecordInRealmSliceBySliceIDParams) (int32, error)
	CloneRecordInRulesetBySliceID(ctx context.Context, arg CloneRecordInRulesetBySliceIDParams) (pgconn.CommandTag, error)
	CloneRecordInSchemaBySliceID(ctx context.Context, arg CloneRecordInSchemaBySliceIDParams) (pgconn.CommandTag, error)
	DeactivateRecord(ctx context.Context, arg DeactivateRecordParams) error
	DeleteCapGranForApp(ctx context.Context, arg DeleteCapGranForAppParams) error
	DeleteWFInstanceListByParents(ctx context.Context, arg DeleteWFInstanceListByParentsParams) ([]Wfinstance, error)
	DeleteWfinstanceByID(ctx context.Context, arg DeleteWfinstanceByIDParams) ([]Wfinstance, error)
	GetApp(ctx context.Context, arg GetAppParams) (string, error)
	GetAppList(ctx context.Context, realm string) ([]GetAppListRow, error)
	GetAppName(ctx context.Context, arg GetAppNameParams) ([]App, error)
	GetCapGrantForApp(ctx context.Context, arg GetCapGrantForAppParams) ([]Capgrant, error)
	GetClass(ctx context.Context, arg GetClassParams) (string, error)
	GetRealmSliceListByRealm(ctx context.Context, realm string) ([]GetRealmSliceListByRealmRow, error)
	GetSchemaWithLock(ctx context.Context, arg GetSchemaWithLockParams) (GetSchemaWithLockRow, error)
	GetWFActiveStatus(ctx context.Context, arg GetWFActiveStatusParams) (pgtype.Bool, error)
	GetWFINstance(ctx context.Context, arg GetWFINstanceParams) (int64, error)
	GetWFInstanceList(ctx context.Context, arg GetWFInstanceListParams) ([]Wfinstance, error)
	GetWFInstanceListByParents(ctx context.Context, id []int32) ([]Wfinstance, error)
	GetWFInternalStatus(ctx context.Context, arg GetWFInternalStatusParams) (bool, error)
	GetWorkflow(ctx context.Context, step string) ([]GetWorkflowRow, error)
	InsertNewRecordInRealmSlice(ctx context.Context, arg InsertNewRecordInRealmSliceParams) (int32, error)
	RealmSliceActivate(ctx context.Context, arg RealmSliceActivateParams) (Realmslice, error)
	RealmSliceAppsList(ctx context.Context, id int32) ([]RealmSliceAppsListRow, error)
	RealmSliceDeactivate(ctx context.Context, arg RealmSliceDeactivateParams) (Realmslice, error)
	RealmSlicePurge(ctx context.Context, realm string) (pgconn.CommandTag, error)
	RulesetRowLock(ctx context.Context, arg RulesetRowLockParams) (Ruleset, error)
	SchemaDelete(ctx context.Context, id int32) (int32, error)
	SchemaGet(ctx context.Context, arg SchemaGetParams) ([]SchemaGetRow, error)
	SchemaNew(ctx context.Context, arg SchemaNewParams) (int32, error)
	SchemaUpdate(ctx context.Context, arg SchemaUpdateParams) error
	UserActivate(ctx context.Context, arg UserActivateParams) (Capgrant, error)
	UserDeactivate(ctx context.Context, arg UserDeactivateParams) (Capgrant, error)
	WfPatternSchemaGet(ctx context.Context, arg WfPatternSchemaGetParams) ([]byte, error)
	WfSchemaGet(ctx context.Context, arg WfSchemaGetParams) (Schema, error)
	WfSchemaList(ctx context.Context, arg WfSchemaListParams) ([]WfSchemaListRow, error)
	Wfschemadelete(ctx context.Context, arg WfschemadeleteParams) error
	Wfschemaget(ctx context.Context, arg WfschemagetParams) (WfschemagetRow, error)
	WorkFlowNew(ctx context.Context, arg WorkFlowNewParams) error
	WorkFlowUpdate(ctx context.Context, arg WorkFlowUpdateParams) (pgconn.CommandTag, error)
	WorkflowDelete(ctx context.Context, arg WorkflowDeleteParams) (pgconn.CommandTag, error)
	WorkflowList(ctx context.Context, arg WorkflowListParams) ([]WorkflowListRow, error)
	Workflowget(ctx context.Context, arg WorkflowgetParams) (WorkflowgetRow, error)
}

var _ Querier = (*Queries)(nil)
