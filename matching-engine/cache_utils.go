package crux

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"


	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/remiges-tech/crux/db/sqlc-gen"
)


func lockCache() {
	cacheLock.Lock()
}

func unlockCache() {
	cacheLock.Unlock()
}

func NewProvider(cfg string) sqlc.Querier {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg)
	if err != nil {
		log.Fatal("error connecting db")
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database")
	return sqlc.NewQuerierWithTX(db)
}

func AddReferencesToRuleSetCache() {
	for realmKey, perRealm := range RulesetCache {
		for _, perApp := range perRealm {
			for sliceKey, perSlice := range perApp {
				for _, rulesets := range perSlice.BRRulesets {
					for _, rule := range rulesets {
						for idx := range rule.Rules {
							processSubRule(&rule.Rules[idx], realmKey, sliceKey)
						}
					}
				}
				for _, rulesets := range perSlice.Workflows {
					for _, rule := range rulesets {
						for idx := range rule.Rules {
							processSubRule(&rule.Rules[idx], realmKey, sliceKey)
						}
					}
				}
			}
		}
	}
}

func processSubRule(subRule *Rule_t, realmKey Realm_t, sliceKey Slice_t) {
	if subRule.RuleActions.ThenCall != "" {
		referRuleset := searchAndAddReferences(subRule.RuleActions.ThenCall, RulesetCache, realmKey, sliceKey, "thencall")
		if referRuleset != nil {
			subRule.RuleActions.References = append(subRule.RuleActions.References, referRuleset)
		}
	}
	if subRule.RuleActions.ElseCall != "" {
		referRuleset := searchAndAddReferences(subRule.RuleActions.ElseCall, RulesetCache, realmKey, sliceKey, "elsecall")
		if referRuleset != nil {
			subRule.RuleActions.References = append(subRule.RuleActions.References, referRuleset)
		}
	}
}

func searchAndAddReferences(targetSetName string, cache map[Realm_t]PerRealm_t, realmKey Realm_t,
	sliceKey Slice_t, calltype string) *Ruleset_t {
	for _, perApp := range cache[realmKey] {
		for otherSliceKey, perSlice := range perApp {
			if otherSliceKey == sliceKey {
				continue
			}
			for _, existingRulesets := range perSlice.BRRulesets {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
						return existingRule

					}
				}
			}
			for _, existingRulesets := range perSlice.Workflows {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
						return existingRule
					}
				}
			}
		}
	}
	return nil
}

func PrintAllRuleSetCache() {

	for realmKey, perRealm := range RulesetCache {
		fmt.Println("Realm:", realmKey)
		for appKey, perApp := range perRealm {
			fmt.Println("\tApp:", appKey)
			for sliceKey, perSlice := range perApp {
				fmt.Println("\t\tSlice:", sliceKey)
				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)

				// Print BRRulesets

				for className, BRRulesets := range perSlice.BRRulesets {
					fmt.Println("\t\t\tBRRulesets - Class:", className)
					for _, rule := range BRRulesets {
						for _, t := range rule.Rules {
							fmt.Println("\t\t\t\tRulePatterns:", t.RulePatterns)
							fmt.Println("\t\t\t\tRuleActions:", t.RuleActions)
							fmt.Println("\t\t\t\tNMatched:", t.NMatched)
							fmt.Println("\t\t\t\tNFailed:", t.NFailed)

							for _, refrule := range t.RuleActions.References {
								for _, z := range refrule.Rules {
									fmt.Println("\t\t\t\t\tReferenced Rule:")
									fmt.Println("\t\t\t\t\t\tRulePatterns:", z.RulePatterns)
									fmt.Println("\t\t\t\t\t\tRuleActions:", z.RuleActions)
									fmt.Println("\t\t\t\t\t\tNMatched:", z.NMatched)
									fmt.Println("\t\t\t\t\t\tNFailed:", z.NFailed)
								}
							}
						}
					}
				}

				// Print Workflows
				for className, workflows := range perSlice.Workflows {
					fmt.Println("\t\t\tWorkflows - Class:", className)
					for _, workflow := range workflows {
						for _, t := range workflow.Rules {
							fmt.Println("\t\t\t\tRulePatterns:", t.RulePatterns)
							fmt.Println("\t\t\t\tRuleActions:", t.RuleActions)
							fmt.Println("\t\t\t\tNMatched:", t.NMatched)
							fmt.Println("\t\t\t\tNFailed:", t.NFailed)

							for _, refrule := range t.RuleActions.References {
								for _, z := range refrule.Rules {
									fmt.Println("\t\t\t\t\tReferenced Rule:")
									fmt.Println("\t\t\t\t\t\tRulePatterns:", z.RulePatterns)
									fmt.Println("\t\t\t\t\t\tRuleActions:", z.RuleActions)
									fmt.Println("\t\t\t\t\t\tNMatched:", z.NMatched)
									fmt.Println("\t\t\t\t\t\tNFailed:", z.NFailed)
								}
							}
						}
					}
				}
			}
		}
	}
}
func PrintAllSchemaCache() {

	for realmKey, perRealm := range SchemaCache {
		fmt.Println("Realm:", realmKey)
		for appKey, perApp := range perRealm {
			fmt.Println("\tApp:", appKey)
			for sliceKey, perSlice := range perApp {
				fmt.Println("\t\tSlice:", sliceKey)
				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)
				for className, schema := range perSlice.BRSchema {
					fmt.Println("\t\t\tBRSchema - Class:", className)
					//for _, schema := range schemas {
					fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
					fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
					fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
					//}
				}
				for className, schema := range perSlice.WFSchema {
					fmt.Println("\t\t\tWFSchema - Class:", className)
					//for _, schema := range schemas {
					fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
					fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
					fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
					//}
				}

			}
		}
	}

}

func containsField(value interface{}, fieldName string, t *testing.T) bool {

	switch v := value.(type) {

	case []byte:

		var raw json.RawMessage
		if err := json.Unmarshal(v, &raw); err != nil {
			fmt.Println("Error unmarshalling actual pattern:", err, v)
			return false
		}

		var data map[string]interface{}
		if err := json.Unmarshal(raw, &data); err != nil {
			var arrayData []interface{}
			if err := json.Unmarshal(raw, &arrayData); err != nil {
				fmt.Println("Error unmarshalling actual pattern:", err, v)
				return false
			}
			for _, element := range arrayData {
				if containsFieldName(element, fieldName) {
					return true
				}
			}
		}
		for _, value := range data {
			if containsFieldName(value, fieldName) {
				return true
			}
		}
	case map[string]interface{}:

		for key := range v {
			if key == fieldName {
				return true
			}
		}

	case []interface{}:
		for _, item := range v {
			if containsField(item, fieldName, t) {
				return true
			}
		}
	case string:
		return v == fieldName
	}
	return false
}

func containsFieldName(value interface{}, fieldName string) bool {

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if key.Kind() == reflect.String && key.String() == fieldName {
				return true
			}
			if nestedValue := v.MapIndex(key).Interface(); containsFieldName(nestedValue, fieldName) {
				return true
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if nestedValue := v.Index(i).Interface(); containsFieldName(nestedValue, fieldName) {
				return true
			}
		}
	case reflect.String:
		return value.(string) == fieldName
	}
	return false
}

func retrieveRuleSchemasFromCache(realm string, app string, class string, slice int) (*Schema_t, error) {
	realmKey := Realm_t(realm)

	perRealm, realmExists := SchemaCache[realmKey]
	if !realmExists {

		return nil, errors.New("schema Realmkey not match")
	}

	appKey := App_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {

		return nil, errors.New("schema AppKey not match")
	}

	sliceKey := Slice_t(slice)

	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {

		return nil, errors.New("schema Slice key not match")
	}

	brSchemas, brExists := perSlice.BRSchema[ClassName_t(class)]
	if brExists {

		return &brSchemas, nil
	}

	wfSchemas, wfExists := perSlice.WFSchema[ClassName_t(class)]
	if wfExists {
		return &wfSchemas, nil
	}

	return nil, nil
}
func convertAttrValue(entityAttrVal string, valType ValType_t) any {

	var entityAttrValConv any
	var err error
	switch valType {
	case ValBool_t:
		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
	case ValInt_t:
		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
	case ValFloat_t:
		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
	case ValString_t, ValEnum_t:
		entityAttrValConv = entityAttrVal
	case ValTimestamp_t:
		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
	}
	if err != nil {
		return err
	}
	return entityAttrValConv
}

func RetrieveRuleSetsFromCache(realm string, app string, class string, slice int) ([]*Ruleset_t, error) {
	realmKey := Realm_t(realm)

	perRealm, realmExists := RulesetCache[realmKey]
	if !realmExists {
		return nil, errors.New("ruleset realmkey not match")
	}

	appKey := App_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {
		return nil, errors.New("ruleset appKey not match")
	}

	sliceKey := Slice_t(slice)
	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {
		return nil, errors.New("ruleset slice key not match")
	}

	var ruleSets []*Ruleset_t

	for _, brRulesets := range perSlice.BRRulesets {
		ruleSets = append(ruleSets, brRulesets...)
	}

	for _, wfRulesets := range perSlice.Workflows {
		ruleSets = append(ruleSets, wfRulesets...)
	}

	return ruleSets, nil
}

func retriveRuleSchemasAndRuleSetsFromCache(realm string, app string, class string, slice string) (*Schema_t, []*Ruleset_t) {
	s, _ := strconv.Atoi(slice)

	ruleSchemas, err := retrieveRuleSchemasFromCache(realm, app, class, s)
	if err != nil {
		log.Printf("Failed to retrieveRuleSchemasFromCache: %v", err)
	}

	ruleSets, err := RetrieveRuleSetsFromCache(realm, app, class, s)
	if err != nil {
		log.Printf("Failed to RetrieveRuleSetsFromCache: %v", err)
	}

	return ruleSchemas, ruleSets
}

func RetrieveRuleSetsByNameFromCache(realm string, app string, class string, slice string, wfname string) *Ruleset_t {
	s, _ := strconv.Atoi(slice)
	// ruleSets, _ := RetrieveRuleSetsFromCache(realm, app, class, s)
	requiredWf := GetWorkflowFromCacheWithName(Realm_t(realm), App_t(app), Slice_t(s), ClassName_t(class), wfname)
	// TODO: check if requiredWf is nil and if not, return error
	return requiredWf
}

func printStats(statsData rulesetStats_t) {
	for realm, perRealm := range statsData {
		for app, perApp := range perRealm {
			for slice, perSlice := range perApp {
				fmt.Printf("Realm: %v, App: %v, Slice: %v\n", realm, app, slice)
				fmt.Printf("loadedAt: %v\n", perSlice.loadedAt)

				// Print stats for BRSchema
				for className, schema := range perSlice.BRSchema {
					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
				}

				// Print stats for BRRulesets
				for className, rulesets := range perSlice.BRRulesets {
					for _, ruleset := range rulesets {
						fmt.Printf("Class: %v, nCalled: %v\n", className, ruleset.nCalled)
						for _, rule := range ruleset.rulesStats {
							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
						}
					}
				}

				// Print stats for WFSchema
				for className, schema := range perSlice.WFSchema {
					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
				}

				// Print stats for Workflows
				for className, workflows := range perSlice.Workflows {
					for _, workflow := range workflows {
						fmt.Printf("Class: %v, nCalled: %v\n", className, workflow.nCalled)
						for _, rule := range workflow.rulesStats {
							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
						}
					}
				}
			}
		}
	}
}
func RetrieveWorkflowRulesetFromCache(realm string, app string, class string, slice int) ([]*Ruleset_t, error) {
	realmKey := Realm_t(realm)

	perRealm, realmExists := RulesetCache[realmKey]
	if !realmExists {
		return nil, errors.New("ruleset realmkey not match")
	}

	appKey := App_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {
		return nil, errors.New("ruleset appKey not match")
	}

	sliceKey := Slice_t(slice)
	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {
		return nil, errors.New("ruleset slice key not match")
	}

	// var ruleSets []*Ruleset_t

	// for _, wfRulesets := range perSlice.Workflows[ClassName_t(class)] {
	// 	ruleSets = append(ruleSets, wfRulesets)
	// }
	// ruleSets = append(ruleSets, perSlice.Workflows[ClassName_t(class)]...)
	ruleSets := perSlice.Workflows[ClassName_t(class)]
	// ruleSets, sliceExists := perApp[sliceKey].Workflows[ClassName_t(class)]
	// if !sliceExists {
	// 	return nil, errors.New("ruleset slice key not match")
	// }

	return ruleSets, nil
}


func GetWorkflowFromCacheWithName(realm Realm_t, app App_t, slice Slice_t, class ClassName_t, wfname string) (w *Ruleset_t) {
	workflows := RulesetCache[realm][app][slice].Workflows[class]
	for _, w := range workflows {
		if w.SetName == wfname {
			return w
		}
	}
	return nil
}
