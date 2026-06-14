package utils

import (
	"fmt"
	"volt/internal/types"

	"github.com/playwright-community/playwright-go"
)

func ProcessStep(page playwright.Page, step types.Step) (string, error) {
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
	default:
		return "", fmt.Errorf("Action: %s (Not implemented)", step.Action)
	}
}
