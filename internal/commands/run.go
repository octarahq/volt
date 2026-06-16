package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	var logs []string
	addToLogs := func(line string) {
		if script.Config.Output != "" {
			logs = append(logs, line)
		}
		fmt.Println(line)
	}

	var scriptTimeout time.Duration
	if strings.TrimSpace(script.Config.Timeout) != "" {
		scriptTimeout, err = time.ParseDuration(strings.TrimSpace(script.Config.Timeout))
		if err != nil {
			addToLogs(fmt.Sprintf("✘ Invalid Timout duration %q: %v", script.Config.Timeout, err))
			return
		}
	}

	browsers := script.Config.Browsers
	if len(browsers) == 0 {
		browsers = append(browsers, "chromium")
	}

	for _, b := range browsers {

		addToLogs(fmt.Sprintf("Volt : %s (%s)", script.Name, b))
		addToLogs("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		addToLogs(fmt.Sprintf("Config : Headless=%t", script.Config.Headless))
		addToLogs("")

		err = playwright.Install()
		if err != nil {
			addToLogs(fmt.Sprintf("✘ Playwright install error: %v", err))
			return
		}

		pw, err := playwright.Run()
		if err != nil {
			addToLogs(fmt.Sprintf("✘ Playwright start error: %v", err))
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
			addToLogs(fmt.Sprintf("✘ Browser launch error: %v", err))
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
			addToLogs("Warning: If you are not browsing your own site, WebKit is very easy to detect as a bot")
			newPageOptions.UserAgent = playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15")
		case "firefox":
			addToLogs("Warning: If you are not browsing your own site, Firefox is very easy to detect as a bot")
		}

		page, err := browser.NewPage(newPageOptions)
		if err != nil {
			addToLogs(fmt.Sprintf("✘ Page creation error: %v", err))
			return
		}

		page.AddInitScript(playwright.Script{Content: playwright.String(stealthScript)})

		state := engine.NewEngineState(script.Vars)
		start := time.Now()

		nbSteps := len(script.Steps)
		for i, s := range script.Steps {
			if scriptTimeout > 0 && time.Since(start) > scriptTimeout {
				addToLogs(fmt.Sprintf("  [%d/%d] ✘ Script timeout exceeded after %s", i+1, nbSteps, scriptTimeout))
				return
			}

			state.InterpolateStep(&s)
			line, err := utils.ProcessStep(page, state, s, script.Config.Humanize, script)
			if err != nil {
				addToLogs(fmt.Sprintf("  [%d/%d] ✘ %s", i+1, nbSteps, err))
				return
			}
			addToLogs(fmt.Sprintf("  [%d/%d] ✔ %s", i+1, nbSteps, line))

			if script.Config.ErrorIfCaptcha {
				if utils.CheckForCaptcha(page) {
					addToLogs(fmt.Sprintf("  [%d/%d] ✘ Captcha detected! Stopping script as requested by ErrorIfCaptcha.", i+1, nbSteps))
					return
				}
			}

			slowmo := script.Config.SlowMo
			if slowmo != "" {
				d, err := time.ParseDuration(slowmo)
				if err != nil {
					addToLogs(fmt.Sprintf("  [%d/%d] ✘ Invalid SlowMo duration %q: %v", i+1, nbSteps, slowmo, err))

					return
				}
				time.Sleep(d)
			}

			if scriptTimeout > 0 && time.Since(start) > scriptTimeout {
				addToLogs(fmt.Sprintf("  [%d/%d] ✘ Script timeout exceeded after %s", i+1, nbSteps, scriptTimeout))
				return
			}
		}
	}

	outPath := strings.TrimSpace(script.Config.Output)
	if outPath != "" {
		fmt.Printf("Save output to %s/output.log\n", outPath)
		if info, err := os.Stat(outPath); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(outPath, 0o755); err != nil {
					fmt.Printf("✘ Unable to create output directory %s: %v\n", outPath, err)
					return
				}
			} else {
				fmt.Printf("✘ Cannot access output directory %s: %v\n", outPath, err)
				return
			}
		} else if !info.IsDir() {
			fmt.Printf("✘ Output path %s is not a directory\n", outPath)
			return
		}

		logFilePath := filepath.Join(outPath, "output.log")
		var logContent strings.Builder
		for _, log := range logs {
			logContent.WriteString(fmt.Sprintf("%s\n", log))
		}
		if err := os.WriteFile(logFilePath, []byte(logContent.String()), 0o644); err != nil {
			fmt.Printf("✘ Unable to write log file %s: %v\n", logFilePath, err)
			return
		}
	}
}
