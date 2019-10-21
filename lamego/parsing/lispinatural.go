package parsing

import (
	"errors"

	// Reflection is used to avoid massive amounts of redundant
	// code to process each LisPI type. To achieve self-hosting
	// it will likely be necessary to generate all the cases,
	// since some languages do not support reflection.
	"reflect"

	"github.com/KernelDeimos/LaME/lamego/model"
)

type SyntaxFrontendLisPINatural struct{}

// TODO: Reports go error type; this needs to be changed before LaME can be self-hosting
func (this SyntaxFrontendLisPINatural) Process(script string) ([]model.SequenceableInstruction, error) {
	tokens, err := ParseListSimple(script)
	if err != nil {
		return nil, err
	}

	output := []model.SequenceableInstruction{}

	for _, statementAsInterface := range tokens {
		switch statement := statementAsInterface.(type) {
		case []interface{}:
			lis, err := this.maybeProcessSequenceable(statement)
			if err != nil {
				return nil, err
			}
			output = append(output, lis)
		default:
			return nil, errors.New("Found string '" + statementAsInterface.(string) + "' when expecting list")
		}
	}

	return output, nil
}

func (this SyntaxFrontendLisPINatural) maybeProcessSequenceable(
	statement []interface{},
) (model.SequenceableInstruction, error) {
	if len(statement) < 1 {
		return nil, errors.New("Found blank list when expecting statement")
	}
	name, ok := statement[0].(string)
	if !ok {
		return nil, errors.New("First token must be a string")
	}

	return this.processSequenceable(name, statement[1:])
}

func (this SyntaxFrontendLisPINatural) maybeProcessExpression(
	expression []interface{},
) (model.ExpressionInstruction, error) {
	if len(expression) < 1 {
		return nil, errors.New("Found blank list when expecting expression")
	}
	name, ok := expression[0].(string)
	if !ok {
		return nil, errors.New("First token must be a string")
	}

	return this.processExpression(name, expression[1:])
}

func (this SyntaxFrontendLisPINatural) processSequenceable(
	name string, args []interface{},
) (model.SequenceableInstruction, error) {
	validSequenceables := map[string]model.SequenceableInstruction{
		"return": model.Return{},
		"iset":   model.ISet{},
	}

	output, recognized := validSequenceables[name]
	if !recognized {
		return nil, errors.New("Expression not recognized:" + name)
	}

	outI, err := this.reflectListToLisPI(reflect.TypeOf(output), args)
	if err != nil {
		return nil, err
	}

	output = reflect.ValueOf(outI).Elem().Interface().(model.SequenceableInstruction)

	return output, nil
}

func (this SyntaxFrontendLisPINatural) processExpression(
	name string, args []interface{},
) (model.ExpressionInstruction, error) {
	validExpressions := map[string]model.ExpressionInstruction{
		"iget": model.IGet{},
	}

	output, recognized := validExpressions[name]
	if !recognized {
		return nil, errors.New("Expression not recognized:" + name)
	}

	outI, err := this.reflectListToLisPI(reflect.TypeOf(output), args)
	if err != nil {
		return nil, err
	}
	output = reflect.ValueOf(outI).Elem().Interface().(model.ExpressionInstruction)

	return output, nil
}

func (this SyntaxFrontendLisPINatural) reflectListToLisPI(
	t reflect.Type, args []interface{},
) (interface{}, error) {
	if len(args) != t.NumField() {
		return nil, errors.New("Wrong number of fields for a " + t.Name())
	}
	var output interface{}
	output = (reflect.New(t).Interface()).(interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		switch field.Type.Name() {
		case "ExpressionInstruction":
			exprList, ok := (args[i]).([]interface{})
			if !ok {
				return nil, errors.New("Found string '" + args[i].(string) + "' when expecting list")
			}
			expr, err := this.maybeProcessExpression(exprList)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(output).Elem().Field(i).Set(reflect.ValueOf(expr))
		case "string":
			strtoken, ok := (args[i]).(string)
			if !ok {
				return nil, errors.New("found list token when expecting string")
			}
			reflect.ValueOf(output).Elem().Field(i).Set(reflect.ValueOf(strtoken))
		case "bool":
			strtoken, ok := (args[i]).(string)
			if !ok {
				return nil, errors.New("found list token when expecting string")
			}
			reflect.ValueOf(output).Elem().Field(i).Set(reflect.ValueOf(strtoken == "true"))
		default:
			// This should never happen
			return nil, errors.New("Unrecognized field type: " + field.Type.Name())
		}
	}
	return output, nil
}
