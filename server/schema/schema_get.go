package schema

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/types"
)

type schemaGetReq struct {
	Slice *int32  `json:"slice" validate:"required"`
	App   *string `json:"app" validate:"required,alpha"`
	Class *string `json:"class" validate:"required,alpha"`
}

// SchemaGet will be responsible for processing the /wfschemaget request that comes through as a POST
func SchemaGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaGet request received")

	// var response schemaGetResp
	var request schemaGetReq
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

	valError := wscutils.WscValidate(request, getValsForSchemaGetReqError)
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}

	dbResponse, err := s.Database.(*sqlc.Queries).Wfschemaget(c, sqlc.WfschemagetParams{
		Slice: *request.Slice,
		App:   *request.App,
		Class: *request.Class,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(types.RECORD_NOT_EXIST, nil)}))
		lh.Debug0().LogActivity("failed to get data from DB:", err.Error)
		return
	}

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": dbResponse}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(dbResponse))
}

func getValsForSchemaGetReqError(err validator.FieldError) []string {
	// validationErrorVals := types.GetErrorValidationMapByAPIName("SchemaGet")
	return types.CommonValidation(err)
}