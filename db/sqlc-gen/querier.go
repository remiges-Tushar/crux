// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
	"encoding/json"
)

type Querier interface {
	WfPatternSchemaGet(ctx context.Context, arg WfPatternSchemaGetParams) (json.RawMessage, error)
	Wfschemadelete(ctx context.Context, arg WfschemadeleteParams) error
	Wfschemaget(ctx context.Context, arg WfschemagetParams) (WfschemagetRow, error)
	Workflowget(ctx context.Context, arg WorkflowgetParams) (WorkflowgetRow, error)
}

var _ Querier = (*Queries)(nil)