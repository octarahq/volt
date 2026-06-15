package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"volt/internal/engine"
	"volt/internal/types"

	"github.com/playwright-community/playwright-go"
)

func ProcessStep(page playwright.Page, state *engine.EngineState, step types.Step) (string, error) {
	switch step.Action {
	case "navigate":
		if _, err := page.Goto(step.URL); err != nil {
			return "", fmt.Errorf("Navigate to %s (Error: %v)", step.URL, err)
		}
		return fmt.Sprintf("Navigate to %s", step.URL), nil
	case "click":
		if err := page.Click(step.Selector); err != nil {
			return "", fmt.Errorf("Click on %s (Error: %v)", step.Selector, err)
		}
		return fmt.Sprintf("Click on %s", step.Selector), nil
	case "type":
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
			absPath, err := filepath.Abs("screenshot.png")
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

	default:
		return "", fmt.Errorf("Action: %s (Not implemented)", step.Action)
	}
}
