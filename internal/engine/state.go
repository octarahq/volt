package engine

import (
	"reflect"
	"regexp"
	"volt/internal/types"
)

type EngineState struct {
	Vars map[string]string
}

func NewEngineState(initialVars map[string]string) *EngineState {
	state := &EngineState{
		Vars: make(map[string]string),
	}
	if initialVars != nil {
		for k, v := range initialVars {
			state.Vars[k] = v
		}
	}
	return state
}

func (s *EngineState) GetVar(name string) string {
	return s.Vars[name]
}

func (s *EngineState) SetVar(name, value string) {
	s.Vars[name] = value
}

func (s *EngineState) Interpolate(text string) string {
	if text == "" {
		return text
	}
	re := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		varName := re.FindStringSubmatch(match)[1]
		if val, ok := s.Vars[varName]; ok {
			return val
		}
		return ""
	})
}

func (s *EngineState) InterpolateStep(step *types.Step) {
	interpolateStruct(s, reflect.ValueOf(step).Elem())
}

func interpolateStruct(state *EngineState, v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(state.Interpolate(field.String()))
		case reflect.Struct:
			interpolateStruct(state, field)
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				interpolateStruct(state, field.Elem())
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					elem.SetString(state.Interpolate(elem.String()))
				}
			}
		case reflect.Map:
			if field.Type().Key().Kind() == reflect.String && field.Type().Elem().Kind() == reflect.String {
				iter := field.MapRange()
				for iter.Next() {
					k := iter.Key()
					vVal := iter.Value()
					newVal := state.Interpolate(vVal.String())
					field.SetMapIndex(k, reflect.ValueOf(newVal))
				}
			}
		}
	}
}
