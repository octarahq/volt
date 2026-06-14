package utils

import (
	"errors"
	"fmt"
	"os"
	"time"
	"volt/internal/types"
)

func CheckScript(script *types.VoltScript) {
	if errs := ValidateScript(script); len(errs) > 0 {
		fmt.Printf("✘ File is invalid. %d error(s) detected :\n", len(errs))
		for _, e := range errs {
			fmt.Printf("  • %v\n", e)
		}
		os.Exit(1)
	}
}

func CheckScriptConsole(script *types.VoltScript) {
	CheckScript(script)
	fmt.Printf("✔ This file is correct, no error has been detected.\n")
}

func ValidateScript(script *types.VoltScript) []error {
	var errs []error
	if script.Name == "" {
		errs = append(errs, errors.New("missing parameter: 'name' is required"))
	}
	if script.Config.SlowMo != "" {
		if _, err := time.ParseDuration(script.Config.SlowMo); err != nil {
			errs = append(errs, fmt.Errorf("invalid config.slow_mo: %v", err))
		}
	}
	if script.Config.Timeout != "" {
		if _, err := time.ParseDuration(script.Config.Timeout); err != nil {
			errs = append(errs, fmt.Errorf("invalid config.timeout: %v", err))
		}
	}
	errs = append(errs, validateSteps(script.Steps, "steps")...)
	return errs
}

func validateSteps(steps []types.Step, contextPath string) []error {
	var errs []error
	for i, step := range steps {
		path := fmt.Sprintf("%s[%d]", contextPath, i)
		if step.Action == "" {
			errs = append(errs, fmt.Errorf("%s: 'action' is required", path))
			continue
		}
		switch step.Action {
		case "navigate":
			if step.URL == "" {
				errs = append(errs, fmt.Errorf("%s: navigate action requires 'url'", path))
			}
		case "click", "hover", "check", "uncheck", "scroll":
			if step.Selector == "" {
				errs = append(errs, fmt.Errorf("%s: %s action requires 'selector'", path, step.Action))
			}
		case "type":
			if step.Selector == "" || step.Value == "" {
				errs = append(errs, fmt.Errorf("%s: type action requires 'selector' and 'value'", path))
			}
		case "select":
			if step.Selector == "" || step.Value == "" {
				errs = append(errs, fmt.Errorf("%s: select action requires 'selector' and 'value'", path))
			}
		case "upload":
			if step.Selector == "" || step.File == "" {
				errs = append(errs, fmt.Errorf("%s: upload action requires 'selector' and 'file'", path))
			}
		case "if":
			if step.Condition == "" {
				errs = append(errs, fmt.Errorf("%s: if action requires 'condition'", path))
			}
			if len(step.Then) == 0 {
				errs = append(errs, fmt.Errorf("%s: if action requires 'then' steps", path))
			}
			errs = append(errs, validateSteps(step.Then, path+".then")...)
			if len(step.Else) > 0 {
				errs = append(errs, validateSteps(step.Else, path+".else")...)
			}
		case "loop":
			if step.Loop == nil || step.Loop.Index == "" {
				errs = append(errs, fmt.Errorf("%s: loop action requires 'loop.index'", path))
			}
			if len(step.Do) == 0 {
				errs = append(errs, fmt.Errorf("%s: loop action requires 'do' steps", path))
			}
			errs = append(errs, validateSteps(step.Do, path+".do")...)
		case "for_each":
			if len(step.ForEach) == 0 || step.Iterator == "" {
				errs = append(errs, fmt.Errorf("%s: for_each action requires 'for_each' list and 'iterator'", path))
			}
			if len(step.Do) == 0 {
				errs = append(errs, fmt.Errorf("%s: for_each action requires 'do' steps", path))
			}
			errs = append(errs, validateSteps(step.Do, path+".do")...)
		}
	}
	return errs
}
