package schema

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

type updateSchema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0,lt=15"`
	App           string              `json:"App" validate:"required,alpha,lt=15"`
	Class         string              `json:"class" validate:"required,lowercase,lt=15"`
	PatternSchema []PatternSchema     `json:"patternSchema,omitempty"`
	ActionSchema  crux.ActionSchema_t `json:"actionSchema,omitempty"`
}

func SchemaUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaUpdate()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
	// 	return
	// }
	capForUpdate := []string{"schema"}
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForUpdate,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req updateSchema

	err = wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query paramaeters to struct:")
		return
	}
	// if req.PatternSchema != nil {
	// 	newPatternSchema := convertPatternSchema(*req.PatternSchema)
	// }
	newPatternSchema := convertPatternSchema(req.PatternSchema)
	schema := crux.Schema_t{
		Class:         req.Class,
		PatternSchema: newPatternSchema,
		ActionSchema:  req.ActionSchema,
		NChecked:      0,
	}

	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(schema)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	patternSchema, err := json.Marshal(newPatternSchema)
	if err != nil {
		patternSchema := "patternSchema"
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &patternSchema)}))
		return
	}

	actionSchema, err := json.Marshal(req.ActionSchema)
	if err != nil {
		actionSchema := "actionSchema"
		l.Debug1().LogDebug("Error while marshaling actionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &actionSchema)}))
		return
	}

	tx, err := connpool.Begin(c)
	if err != nil {
		l.Info().Error(err).Log("Error while beginning transaction")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)
	getSchema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		RealmName:   realmName,
		ID:          req.Slice,
		Class:       req.Class,
		Shortnamelc: strings.ToLower(req.App),
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while locking schema to get old value")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
		RealmName:     realmName,
		Slice:         req.Slice,
		Class:         req.Class,
		App:           strings.ToLower(req.App),
		Brwf:          sqlc.BrwfEnumW,
		Patternschema: patternSchema,
		Actionschema:  actionSchema,
		Editedby:      pgtype.Text{String: userID, Valid: true},
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while updating schema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	// if patternSchema != nil && actionSchema == nil {
	// 	err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
	// 		Realm:         realmName,
	// 		Slice:         req.Slice,
	// 		Class:         req.Class,
	// 		App:           req.App,
	// 		Brwf:          sqlc.BrwfEnumW,
	// 		Patternschema: patternSchema,
	// 		Actionschema:  actionSchema,
	// 		Editedby:      pgtype.Text{String: userID},
	// 	})
	// 	if err != nil {
	// 		tx.Rollback(c)

	// 		errmsg := db.HandleDatabaseError(err)
	// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
	// 		return
	// 	}
	// } else if actionSchema != nil && patternSchema == nil {
	// 	err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
	// 		Realm:         realmName,
	// 		Slice:         sh.Slice,
	// 		Class:         sh.Class,
	// 		App:           sh.App,
	// 		Brwf:          sqlc.BrwfEnumW,
	// 		Patternschema: schema.Patternschema,
	// 		Actionschema:  actionSchema,
	// 		Editedby:      pgtype.Text{String: userID},
	// 	})
	// 	if err != nil {
	// 		tx.Rollback(c)

	// 		errmsg := db.HandleDatabaseError(err)
	// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
	// 		return
	// 	}
	// } else {
	// 	err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
	// 		Realm:         realmName,
	// 		Slice:         sh.Slice,
	// 		Class:         sh.Class,
	// 		App:           sh.App,
	// 		Brwf:          sqlc.BrwfEnumW,
	// 		Patternschema: schema.Patternschema,
	// 		Actionschema:  schema.Actionschema,
	// 		Editedby:      pgtype.Text{String: userID},
	// 	})
	// 	if err != nil {
	// 		tx.Rollback(c)

	// 		errmsg := db.HandleDatabaseError(err)
	// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
	// 		return
	// 	}

	// }

	if err := tx.Commit(c); err != nil {

		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	dclog := l.WithClass("schema").WithInstanceId(string(getSchema.ID))
	dclog.LogDataChange("Updated schema", logharbour.ChangeInfo{
		Entity: "schema",
		Op:     "Update",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "patternSchema",
				OldVal: string(getSchema.Patternschema),
				NewVal: newPatternSchema},
			{
				Field:  "actionSchema",
				OldVal: string(getSchema.Actionschema),
				NewVal: req.ActionSchema},
		},
	})

	// err = schemaUpdateWithTX(c, query, connpool, l, req)
	// if err != nil {
	// 	l.LogActivity("Error while Updating schema", err.Error())
	// 	errmsg := db.HandleDatabaseError(err)
	// 	wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
	// 	return
	// }
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of SchemaUpdate()")
}

// func schemaUpdateWithTX(c context.Context, query *sqlc.Queries, connpool *pgxpool.Pool, l *logharbour.Logger, sh updateSchema) error {
// 	patternSchema, err := json.Marshal(sh.PatternSchema)
// 	if err != nil {
// 		l.Debug1().Error(err).Log("Error while marshaling patternSchema")
// 		return err
// 	}

// 	actionSchema, err := json.Marshal(sh.ActionSchema)
// 	if err != nil {
// 		l.Debug1().Error(err).Log("Error while marshaling actionSchema")
// 		return err
// 	}

// 	tx, err := connpool.Begin(c)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(c)
// 	qtx := query.WithTx(tx)
// 	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
// 		Realm: realmName,
// 		Slice: sh.Slice,
// 		Class: sh.Class,
// 		App:   sh.App,
// 	})
// 	if err != nil {
// 		tx.Rollback(c)
// 		return err
// 	}
// 	if patternSchema != nil && actionSchema == nil {
// 		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
// 			Realm:         realmName,
// 			Slice:         sh.Slice,
// 			Class:         sh.Class,
// 			App:           sh.App,
// 			Brwf:          sqlc.BrwfEnumW,
// 			Patternschema: schema.Patternschema,
// 			Actionschema:  actionSchema,
// 			Editedby:      pgtype.Text{String: userID},
// 		})
// 		if err != nil {
// 			tx.Rollback(c)
// 			return err
// 		}
// 	} else if actionSchema != nil && patternSchema == nil {
// 		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
// 			Realm:         realmName,
// 			Slice:         sh.Slice,
// 			Class:         sh.Class,
// 			App:           sh.App,
// 			Brwf:          sqlc.BrwfEnumW,
// 			Patternschema: schema.Patternschema,
// 			Actionschema:  actionSchema,
// 			Editedby:      pgtype.Text{String: userID},
// 		})
// 		if err != nil {
// 			tx.Rollback(c)
// 			return err
// 		}
// 	} else {
// 		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
// 			Realm:         realmName,
// 			Slice:         sh.Slice,
// 			Class:         sh.Class,
// 			App:           sh.App,
// 			Brwf:          sqlc.BrwfEnumW,
// 			Patternschema: schema.Patternschema,
// 			Actionschema:  schema.Actionschema,
// 			Editedby:      pgtype.Text{String: userID},
// 		})
// 		if err != nil {
// 			tx.Rollback(c)
// 			return err
// 		}

// 	}

// 	if err := tx.Commit(c); err != nil {
// 		return err
// 	}
// 	dclog := l.WithClass("schema").WithInstanceId(string(schema.ID))
// 	dclog.LogDataChange("Updated schema", logharbour.ChangeInfo{
// 		Entity: "schema",
// 		Op:     "Update",
// 		Changes: []logharbour.ChangeDetail{
// 			{
// 				Field:  "patternSchema",
// 				OldVal: string(schema.Patternschema),
// 				NewVal: patternSchema},
// 			{
// 				Field:  "actionSchema",
// 				OldVal: string(schema.Actionschema),
// 				NewVal: sh.ActionSchema},
// 		},
// 	})

// 	return nil
// }

// func customValidationErrorsForUpdate(sh updateSchema) []wscutils.ErrorMessage {
// 	var validationErrors []wscutils.ErrorMessage
// 	if sh.PatternSchema != nil && sh.ActionSchema == nil {
// 		patternSchemaError := verifyPatternSchemaUpdate(sh.PatternSchema)
// 		validationErrors = append(validationErrors, patternSchemaError...)
// 	} else if sh.ActionSchema != nil && sh.PatternSchema == nil {
// 		actionSchemaError := verifyActionSchemaUpdate(sh.ActionSchema)
// 		validationErrors = append(validationErrors, actionSchemaError...)
// 	} else if sh.PatternSchema == nil && sh.ActionSchema == nil {
// 		fieldName := fmt.Sprintln("PatternSchema/ActionSchema")
// 		vErr := wscutils.BuildErrorMessage(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne, &fieldName)
// 		validationErrors = append(validationErrors, vErr)
// 	} else {
// 		patternSchemaError := verifyPatternSchemaUpdate(sh.PatternSchema)
// 		validationErrors = append(validationErrors, patternSchemaError...)
// 		actionSchemaError := verifyActionSchemaUpdate(sh.ActionSchema)
// 		validationErrors = append(validationErrors, actionSchemaError...)
// 	}

// 	return validationErrors
// }
// func verifyPatternSchemaUpdate(ps *patternSchema) []wscutils.ErrorMessage {
// 	var validationErrors []wscutils.ErrorMessage

// 	for i, attrSchema := range ps.Attr {
// 		i++
// 		if !re.MatchString(attrSchema.Name) {
// 			fieldName := fmt.Sprintf("attrSchema[%d].Name", i)
// 			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, attrSchema.Name)
// 			validationErrors = append(validationErrors, vErr)
// 		}
// 		if !validTypes[attrSchema.ValType] {
// 			fieldName := fmt.Sprintf("attrSchema[%d].ValType", i)
// 			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, attrSchema.ValType)
// 			validationErrors = append(validationErrors, vErr)
// 		}
// 		if attrSchema.ValType == "enum" && len(attrSchema.Vals) == 0 {
// 			fieldName := fmt.Sprintf("attrSchema[%d].Vals", i)
// 			vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
// 			validationErrors = append(validationErrors, vErr)
// 		}
// 	}
// 	return validationErrors
// }

// func verifyActionSchemaUpdate(as *actionSchema) []wscutils.ErrorMessage {
// 	var validationErrors []wscutils.ErrorMessage
// 	re := regexp.MustCompile(cruxIDRegExp)

// 	for i, task := range as.Tasks {
// 		if !re.MatchString(task) {
// 			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
// 			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, task)
// 			validationErrors = append(validationErrors, vErr)
// 		}
// 	}
// 	for i, propName := range as.Properties {
// 		if !re.MatchString(propName) {
// 			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
// 			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, propName)
// 			validationErrors = append(validationErrors, vErr)
// 		}
// 	}
// 	return validationErrors
// }
