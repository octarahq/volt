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
	if len(script.Config.Browsers) > 0 {
		correctBrowsers := []string{"chromium", "firefox", "webkit"}
		for i, b := range script.Config.Browsers {
			valid := false
			for _, cb := range correctBrowsers {
				if b == cb {
					valid = true
					break
				}
			}
			if !valid {
				errs = append(errs, fmt.Errorf("invalid config.browsers[%d]: %s (allowed: %v)", i, b, correctBrowsers))
			}
		}
	}
	errs = append(errs, validateSteps(script.Steps, "steps")...)
	return errs
}

func validateSteps(steps []types.Step, contextPath string) []error {
	var errs []error
	for i, step := range steps {
		path := fmt.Sprintf("%s[%d]", contextPath, i+1)
		if step.Action == "" {
			errs = append(errs, fmt.Errorf("%s: 'action' is required", path))
			continue
		}
		switch step.Action {
		case "navigate":
			if step.URL == "" {
				errs = append(errs, fmt.Errorf("%s: navigate action requires 'url'", path))
			}
		case "click", "hover", "check", "uncheck":
			if step.Selector == "" {
				errs = append(errs, fmt.Errorf("%s: %s action requires 'selector'", path, step.Action))
			}
		case "scroll":
			if step.Selector == "" && step.Value != "top" && step.Value != "bottom" {
				errs = append(errs, fmt.Errorf("%s: scroll action requires 'selector' or value 'top'/'bottom'", path))
			}
		case "log":
			if step.Message == "" {
				errs = append(errs, fmt.Errorf("%s: log action requires 'message'", path))
			}
		case "store_text":
			if step.Name == "" {
				errs = append(errs, fmt.Errorf("%s: store_text action requires variable name 'name'", path))
			}
			if step.Selector == "" {
				errs = append(errs, fmt.Errorf("%s: store_text action requires 'selector'", path))
			}
		case "store_attribute":
			if step.Name == "" {
				errs = append(errs, fmt.Errorf("%s: store_attribute action requires variable name 'name'", path))
			}
			if step.Selector == "" && step.Attribute == "" {
				errs = append(errs, fmt.Errorf("%s: store_attribute action requires 'selector' and 'attribute'", path))
			}
		case "store_eval":
			if step.Name == "" {
				errs = append(errs, fmt.Errorf("%s: store_eval action requires variable name 'name'", path))
			}
			if step.Value == "" {
				errs = append(errs, fmt.Errorf("%s: store_eval action requires 'value'", path))
			}
		case "store_value":
			if step.Name == "" {
				errs = append(errs, fmt.Errorf("%s: store_value action requires variable name 'name'", path))
			}
			if step.As == "" {
				errs = append(errs, fmt.Errorf("%s: store_value action requires 'as'", path))
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
			if step.To == 0 {
				errs = append(errs, fmt.Errorf("%s: loop action requires 'to'", path))
			} else {
				if step.From > step.To {
					errs = append(errs, fmt.Errorf("%s: Loop range is invalid: from %d is greater than to %d", path, step.From, step.To))
				}
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
		case "screenshot":
			continue
		case "clear_cookies":
			continue
		case "add_header":
			if step.Name == "" || step.Value == "" {
				errs = append(errs, fmt.Errorf("%s: add_header action requires 'name' and 'value'", path))
			}
		case "set_header":
			if step.Name == "" || step.Value == "" {
				errs = append(errs, fmt.Errorf("%s: set_header action requires 'name' and 'value'", path))
			}
		case "remove_header":
			if step.Name == "" {
				errs = append(errs, fmt.Errorf("%s: remove_header action requires 'name'", path))
			}
		case "wait":
			if step.Duration == 0 {
				errs = append(errs, fmt.Errorf("%s: wait action requires 'duration'", path))
			}
		case "wait_visible":
			if step.Selector == "" {
				errs = append(errs, fmt.Errorf("%s: wait_visible action requires 'selector'", path))
			}
		case "wait_hidden":
			if step.Selector == "" {
				errs = append(errs, fmt.Errorf("%s: wait_hidden action requires 'selector'", path))
			}
		case "scrape":
			if step.Scrape == nil || len(step.Scrape.Fields) == 0 {
				errs = append(errs, fmt.Errorf("%s: scrape action requires 'scrape.fields'", path))
			}
		default:
			errs = append(errs, fmt.Errorf("%s: Action: %s (Not implemented)", path, step.Action))
		}
	}
	return errs
}
