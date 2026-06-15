package commands

import (
	"fmt"
	"os"
	"volt/internal/engine"
	"volt/internal/utils"

	"github.com/playwright-community/playwright-go"
)

func Run(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		utils.ErrorCannotReadFile(path)
		return
	}

	script, err := utils.ParseScriptFile(path, data)
	if err != nil {
		fmt.Printf("✘ Error parsing file: %v\n", err)
		return
	}

	utils.CheckScript(&script)

	fmt.Printf("Volt : %s\n", script.Name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Config : Headless=%t\n", script.Config.Headless)
	fmt.Println()

	err = playwright.Install()
	if err != nil {
		fmt.Printf("✘ Playwright install error: %v\n", err)
		return
	}

	pw, err := playwright.Run()
	if err != nil {
		fmt.Printf("✘ Playwright start error: %v\n", err)
		return
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(script.Config.Headless),
		Args: []string{
			"--disable-blink-features=AutomationControlled",
		},
	})
	if err != nil {
		fmt.Printf("✘ Browser launch error: %v\n", err)
		return
	}
	defer browser.Close()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	})
	if err != nil {
		fmt.Printf("✘ Page creation error: %v\n", err)
		return
	}

	stealthScript := `
		Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
		window.navigator.chrome = { runtime: {} };
		Object.defineProperty(navigator, 'plugins', {get: () => [1, 2, 3]});
		Object.defineProperty(navigator, 'languages', {get: () => ['en-US', 'en', 'fr']});
	`
	page.AddInitScript(playwright.Script{Content: playwright.String(stealthScript)})

	state := engine.NewEngineState(script.Vars)

	nbSteps := len(script.Steps)
	for i, s := range script.Steps {
		state.InterpolateStep(&s)
		line, err := utils.ProcessStep(page, state, s)
		if err != nil {
			fmt.Printf("  [%d/%d] ✘ %s\n", i+1, nbSteps, err)
			return
		}
		fmt.Printf("  [%d/%d] ✔ %s\n", i+1, nbSteps, line)
	}
}
