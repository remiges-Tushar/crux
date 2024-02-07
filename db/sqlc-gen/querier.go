// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddWFNewInstace(ctx context.Context, arg AddWFNewInstaceParams) (string, error)
	GetApp(ctx context.Context, arg GetAppParams) (string, error)
	GetClass(ctx context.Context, arg GetClassParams) (string, error)
	GetWFActiveStatus(ctx context.Context, arg GetWFActiveStatusParams) (pgtype.Bool, error)
	GetWFINstance(ctx context.Context, arg GetWFINstanceParams) ([]Wfinstance, error)
	GetWFInternalStatus(ctx context.Context, arg GetWFInternalStatusParams) (bool, error)
	SchemaDelete(ctx context.Context, id int32) (int32, error)
	SchemaGet(ctx context.Context, arg SchemaGetParams) ([]SchemaGetRow, error)
	SchemaList(ctx context.Context) ([]SchemaListRow, error)
	SchemaListByApp(ctx context.Context, app string) ([]SchemaListByAppRow, error)
	SchemaListByAppAndClass(ctx context.Context, arg SchemaListByAppAndClassParams) ([]SchemaListByAppAndClassRow, error)
	SchemaListByAppAndSlice(ctx context.Context, arg SchemaListByAppAndSliceParams) ([]SchemaListByAppAndSliceRow, error)
	SchemaListByClass(ctx context.Context, class string) ([]SchemaListByClassRow, error)
	SchemaListByClassAndSlice(ctx context.Context, arg SchemaListByClassAndSliceParams) ([]SchemaListByClassAndSliceRow, error)
	SchemaListBySlice(ctx context.Context, slice int32) ([]SchemaListBySliceRow, error)
	SchemaNew(ctx context.Context, arg SchemaNewParams) (int32, error)
	SchemaUpdate(ctx context.Context, arg SchemaUpdateParams) (int32, error)
	UpdateSchemaWithLock(ctx context.Context, arg UpdateSchemaWithLockParams) (UpdateSchemaWithLockRow, error)
	WfPatternSchemaGet(ctx context.Context, arg WfPatternSchemaGetParams) ([]byte, error)
	WfSchemaGet(ctx context.Context, arg WfSchemaGetParams) (Schema, error)
	Wfschemadelete(ctx context.Context, arg WfschemadeleteParams) error
	Wfschemaget(ctx context.Context, arg WfschemagetParams) (WfschemagetRow, error)
	WorkFlowNew(ctx context.Context, arg WorkFlowNewParams) (int32, error)
	Workflowget(ctx context.Context, arg WorkflowgetParams) (WorkflowgetRow, error)
}

var _ Querier = (*Queries)(nil)
