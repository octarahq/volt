package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"volt/internal/engine"
	"volt/internal/types"

	"github.com/expr-lang/expr"
	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
)

func ProcessStep(page playwright.Page, state *engine.EngineState, step types.Step, humanize bool, script types.VoltScript, logFunc func(string)) (string, error) {
	switch step.Action {
	case "navigate":
		if _, err := page.Goto(step.URL); err != nil {
			return "", fmt.Errorf("Navigate to %s (Error: %v)", step.URL, err)
		}
		if humanize {
			InitMousePointer(page)
		}
		return fmt.Sprintf("Navigate to %s", step.URL), nil
	case "click":
		if humanize {
			HumanizeMouse(page, step.Selector)
		}
		if err := page.Click(step.Selector); err != nil {
			return "", fmt.Errorf("Click on %s (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Click on %s", step.Selector), nil
	case "type":
		if humanize {
			HumanizeMouse(page, step.Selector)
		}
		if err := page.Locator(step.Selector).Fill(step.Value); err != nil {
			return "", fmt.Errorf("Type '%s' into %s (Error: %v)", step.Value, step.Selector, err)
		}
		return fmt.Sprintf("Type '%s' into %s", step.Value, step.Selector), nil
	case "press_key":
		if err := page.Locator(step.Selector).Press(step.Value); err != nil {
			return "", fmt.Errorf("Press '%s' on %s (Error: %v)", step.Value, step.Selector, err)
		}
		return fmt.Sprintf("Press '%s' on %s", step.Value, step.Selector), nil
	case "hover":
		if humanize {
			HumanizeMouse(page, step.Selector)
		}
		if err := page.Locator(step.Selector).Hover(); err != nil {
			return "", fmt.Errorf("Hover on %s (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Hover on %s", step.Selector), nil
	case "check":
		if err := page.Locator(step.Selector).Check(); err != nil {
			return "", fmt.Errorf("Check %s (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Check %s", step.Selector), nil
	case "uncheck":
		if err := page.Locator(step.Selector).Uncheck(); err != nil {
			return "", fmt.Errorf("Uncheck %s (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Uncheck %s", step.Selector), nil
	case "select":
		if _, err := page.Locator(step.Selector).SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(step.Value)}); err != nil {
			return "", fmt.Errorf("Select '%s' in %s (Error: %v)", step.Value, step.Selector, err)
		}
		return fmt.Sprintf("Select '%s' in %s", step.Value, step.Selector), nil
	case "upload":
		if err := page.Locator(step.Selector).SetInputFiles(step.File); err != nil {
			return "", fmt.Errorf("Upload file '%s' to %s (Error: %v)", step.File, step.Selector, err)
		}
		return fmt.Sprintf("Upload file '%s' to %s", step.File, step.Selector), nil
	case "log":
		return fmt.Sprintf("Console : %s", step.Message), nil
	case "store_value":
		name := step.Name
		value := step.As

		state.SetVar(name, value)
		return fmt.Sprintf("Store value in %s as '%s'", name, value), nil
	case "store_text":
		name := step.Name

		value, err := page.Locator(step.Selector).TextContent()
		if err != nil {
			return "", fmt.Errorf("Store text from %s (Error: %v)", step.Selector, err)
		}

		state.SetVar(name, value)
		return fmt.Sprintf("Store text in %s as '%s'", name, value), nil
	case "store_attribute":
		name := step.Name

		value, err := page.Locator(step.Selector).GetAttribute(step.Attribute)
		if err != nil {
			return "", fmt.Errorf("Store attribute '%s' from %s (Error: %v)", step.Value, step.Selector, err)
		}

		state.SetVar(name, value)
		return fmt.Sprintf("Store text in %s as '%s'", name, value), nil
	case "store_eval":
		name := step.Name

		value, err := page.Evaluate(step.Value)
		if err != nil {
			return "", fmt.Errorf("Store eval of %s (Error: %v)", step.Value, err)
		}

		stringValue := fmt.Sprint(value)
		state.SetVar(name, stringValue)
		return fmt.Sprintf("Store eval in %s as '%s'", name, stringValue), nil
	case "scroll":
		if step.Selector != "" {
			if err := page.Locator(step.Selector).ScrollIntoViewIfNeeded(); err != nil {
				return "", fmt.Errorf("Scroll to %s (Error: %v)", step.Selector, err)
			}
			return fmt.Sprintf("Scroll to %s", step.Selector), nil
		} else if step.Value == "bottom" {
			if _, err := page.Evaluate("window.scrollTo(0, document.body.scrollHeight)"); err != nil {
				return "", fmt.Errorf("Scroll to bottom (Error: %v)", err)
			}
			return "Scroll to bottom", nil
		} else if step.Value == "top" {
			if _, err := page.Evaluate("window.scrollTo(0, 0)"); err != nil {
				return "", fmt.Errorf("Scroll to top (Error: %v)", err)
			}
			return "Scroll to top", nil
		}
		return "", fmt.Errorf("Scroll action requires a selector, or value 'top' or 'bottom'")

	case "screenshot":
		path := step.Pathname
		if path == "" {
			var absPath string
			var err error
			filename := fmt.Sprintf("%s.png", uuid.NewString())
			if script.Config.Output != "" {
				absPath, err = filepath.Abs(fmt.Sprintf("%s/screenshots/%s", script.Config.Output, filename))
			} else {
				absPath, err = filepath.Abs(filename)
			}
			if err != nil {
				return "", fmt.Errorf("Get absolute screenshot path (Error: %v)", err)
			}
			path = absPath
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", fmt.Errorf("Create screenshot directory (Error: %v)", err)
		}

		opts := playwright.PageScreenshotOptions{
			Path: playwright.String(path),
		}

		if step.Position != nil {
			pos := step.Position
			if pos.FullPage {
				opts.FullPage = playwright.Bool(true)
			} else if pos.Width > 0 || pos.Height > 0 || pos.X != 0 || pos.Y != 0 {
				opts.Clip = &playwright.Rect{
					X:      float64(pos.X),
					Y:      float64(pos.Y),
					Width:  float64(pos.Width),
					Height: float64(pos.Height),
				}
			}
		}

		if _, err := page.Screenshot(opts); err != nil {
			return "", fmt.Errorf("Take screenshot and store to %s (Error: %v)", path, err)
		}

		return fmt.Sprintf("Take screenshot and store to %s", path), nil
	case "clear_cookies":
		if err := page.Context().ClearCookies(); err != nil {
			return "", fmt.Errorf("Clear cookies (Error: %v)", err)
		}
		return "Clear cookies", nil
	case "add_header":
		name := step.Name
		value := step.Value

		if err := page.Context().SetExtraHTTPHeaders(map[string]string{name: value}); err != nil {
			return "", fmt.Errorf("Add header '%s: %s' (Error: %v)", name, value, err)
		}
		return fmt.Sprintf("Add header '%s: %s'", name, value), nil
	case "set_header":
		name := step.Name
		value := step.Value

		if err := page.Context().SetExtraHTTPHeaders(map[string]string{name: value}); err != nil {
			return "", fmt.Errorf("Set header '%s: %s' (Error: %v)", name, value, err)
		}
		return fmt.Sprintf("Set header '%s: %s'", name, value), nil
	case "remove_header":
		name := step.Name

		if err := page.Context().SetExtraHTTPHeaders(map[string]string{name: ""}); err != nil {
			return "", fmt.Errorf("Remove header '%s' (Error: %v)", name, err)
		}
		return fmt.Sprintf("Remove header '%s'", name), nil
	case "wait":
		duration := step.Duration
		time.Sleep(time.Duration(duration) * time.Second)
		return fmt.Sprintf("Wait for %d seconds", duration), nil
	case "wait_visible":
		if err := page.Locator(step.Selector).WaitFor(); err != nil {
			return "", fmt.Errorf("Wait for %s to be visible (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Wait for %s to be visible", step.Selector), nil
	case "wait_hidden":
		if err := page.Locator(step.Selector).WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateHidden}); err != nil {
			return "", fmt.Errorf("Wait for %s to be hidden (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Wait for %s to be hidden", step.Selector), nil

	case "scrape":
		if step.Scrape == nil {
			return "", fmt.Errorf("Scrape action requires 'scrape' configuration")
		}

		var results interface{}
		if step.Scrape.Parent != "" {
			elements, err := page.Locator(step.Scrape.Parent).All()
			if err != nil {
				return "", fmt.Errorf("Finding parent elements %s (Error: %v)", step.Scrape.Parent, err)
			}

			items := make([]map[string]string, 0)
			for _, el := range elements {
				item := make(map[string]string)
				for field, fieldSelector := range step.Scrape.Fields {
					selector := fieldSelector
					attribute := ""
					if strings.Contains(selector, "@") {
						parts := strings.SplitN(selector, "@", 2)
						selector = strings.TrimSpace(parts[0])
						attribute = strings.TrimSpace(parts[1])
					}

					var val string
					var err error
					var locator playwright.Locator
					if selector != "" {
						locator = el.Locator(selector).First()
					} else {
						locator = el
					}

					if attribute != "" {
						val, err = locator.GetAttribute(attribute)
					} else {
						val, err = locator.TextContent()
					}

					if err == nil {
						item[field] = strings.TrimSpace(val)
					} else {
						item[field] = ""
					}
				}
				items = append(items, item)
			}
			results = items
		} else {
			item := make(map[string]string)
			for field, fieldSelector := range step.Scrape.Fields {
				selector := fieldSelector
				attribute := ""
				if strings.Contains(selector, "@") {
					parts := strings.SplitN(selector, "@", 2)
					selector = strings.TrimSpace(parts[0])
					attribute = strings.TrimSpace(parts[1])
				}

				var val string
				var err error
				if selector != "" {
					locator := page.Locator(selector).First()
					if attribute != "" {
						val, err = locator.GetAttribute(attribute)
					} else {
						val, err = locator.TextContent()
					}
				} else {
					val = ""
					err = fmt.Errorf("selector is required for single item scrape field")
				}

				if err == nil {
					item[field] = strings.TrimSpace(val)
				} else {
					item[field] = ""
				}
			}
			results = item
		}

		jsonBytes, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return "", fmt.Errorf("Serializing scrape results (Error: %v)", err)
		}

		if step.Name != "" {
			state.SetVar(step.Name, string(jsonBytes))
		}

		if script.Config.Output != "" {
			filename := "scraped_data.json"
			if step.Name != "" {
				filename = fmt.Sprintf("scraped_%s.json", step.Name)
			}
			path := filepath.Join(script.Config.Output, filename)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return "", fmt.Errorf("Create output directory (Error: %v)", err)
			}
			if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
				return "", fmt.Errorf("Write scrape results to %s (Error: %v)", path, err)
			}
			return fmt.Sprintf("Scraped data saved to %s", path), nil
		}

		return "Scraped data (stored in memory)", nil
	case "loop":
		from := step.From
		to := step.To

		iterations := to - from + 1
		for i := from; i <= to; i++ {
			state.SetVar("loop.index", strconv.Itoa(i))
			outerIdx := i - from + 1
			for j, s := range step.Do {
				state.InterpolateStep(&s)
				result, err := ProcessStep(page, state, s, humanize, script, logFunc)
				if err != nil {
					return "", err
				}
				if result != "" {
					logFunc(fmt.Sprintf("  [%d/%d] [%d/%d] ✔ %s", outerIdx, iterations, j+1, len(step.Do), result))
				}
			}
		}
		return fmt.Sprintf("Loop executed from %d to %d: ran %d iterations with %d step(s) each", from, to, iterations, len(step.Do)), nil
	case "if":
		condition := step.Condition

		env := map[string]interface{}{
			"false": false,
			"true":  true,
		}

		for n, v := range state.Vars {
			env[n] = v
		}

		program, err := expr.Compile(step.Condition, expr.Env(env))
		if err != nil {
			return "", fmt.Errorf("Invalid condition '%s' : %v", condition, err)
		}

		output, err := expr.Run(program, env)
		if err != nil {
			return "", fmt.Errorf("Execution error: %v", err)
		}

		var result bool
		switch v := output.(type) {
		case bool:
			result = v
		case string:
			if v == "true" {
				result = true
			} else if v == "false" {
				result = false
			} else {
				return "", fmt.Errorf("Condition '%s' returned a string \"%s\" which is not 'true' or 'false'. Make sure the condition evaluates to a boolean.", condition, v)
			}
		default:
			return "", fmt.Errorf("Condition '%s' did not evaluate to a boolean (got %T: %v)", condition, output, output)
		}

		if result {
			for j, s := range step.Then {
				state.InterpolateStep(&s)
				result, err := ProcessStep(page, state, s, humanize, script, logFunc)
				if err != nil {
					return "", err
				}
				if result != "" {
					logFunc(fmt.Sprintf("      [%d/%d] ✔ %s", j+1, len(step.Then), result))
				}
			}
		} else {
			for j, s := range step.Else {
				state.InterpolateStep(&s)
				result, err := ProcessStep(page, state, s, humanize, script, logFunc)
				if err != nil {
					return "", err
				}
				if result != "" {
					logFunc(fmt.Sprintf("      [%d/%d] ✔ %s", j+1, len(step.Else), result))
				}
			}
		}
		return "If executed", nil
	default:
		return "", fmt.Errorf("Action: %s (Not implemented)", step.Action)
	}
}
