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

	browsers := script.Config.Browsers
	if len(browsers) == 0 {
		browsers = append(browsers, "chromium")
	}

	for _, b := range browsers {

		fmt.Printf("Volt : %s (%s)\n", script.Name, b)
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

		var browser playwright.Browser
		switch b {
		case "firefox":
			browser, err = pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(script.Config.Headless),
			})
		case "webkit":
			browser, err = pw.WebKit.Launch(playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(script.Config.Headless),
			})
		default:
			browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(script.Config.Headless),
				Args: []string{
					"--disable-blink-features=AutomationControlled",
				},
			})
		}
		if err != nil {
			fmt.Printf("✘ Browser launch error: %v\n", err)
			return
		}
		defer browser.Close()

		var stealthScript string

		switch b {
		case "firefox":
			stealthScript = `
				Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
			`
		case "webkit":
			stealthScript = `
				Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
				Object.defineProperty(navigator, 'platform', {get: () => 'MacIntel'});
			`
		default:
			stealthScript = `
				Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
				window.navigator.chrome = { runtime: {} };
			`
		}

		var newPageOptions playwright.BrowserNewPageOptions
		switch b {
		case "chromium":
			newPageOptions.UserAgent = playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		case "webkit":
			fmt.Println("Warning: If you are not browsing your own site, WebKit is very easy to detect as a bot")
			newPageOptions.UserAgent = playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15")
		case "firefox":
			fmt.Println("Warning: If you are not browsing your own site, Firefox is very easy to detect as a bot")
		}

		page, err := browser.NewPage(newPageOptions)
		if err != nil {
			fmt.Printf("✘ Page creation error: %v\n", err)
			return
		}

		page.AddInitScript(playwright.Script{Content: playwright.String(stealthScript)})

		state := engine.NewEngineState(script.Vars)

		nbSteps := len(script.Steps)
		for i, s := range script.Steps {
			state.InterpolateStep(&s)
			line, err := utils.ProcessStep(page, state, s, script.Config.Humanize)
			if err != nil {
				fmt.Printf("  [%d/%d] ✘ %s\n", i+1, nbSteps, err)
				return
			}
			fmt.Printf("  [%d/%d] ✔ %s\n", i+1, nbSteps, line)

			if script.Config.ErrorIfCaptcha {
				if utils.CheckForCaptcha(page) {
					fmt.Printf("  [%d/%d] ✘ Captcha detected! Stopping script as requested by ErrorIfCaptcha.\n", i+1, nbSteps)
					return
				}
			}
		}
	}
}
