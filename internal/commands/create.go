package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"volt/internal/types"

	"github.com/goccy/go-yaml"
	"github.com/playwright-community/playwright-go"
)

func Create() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Script name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Timeout (e.g. 30s): ")
	timeout, _ := reader.ReadString('\n')
	timeout = strings.TrimSpace(timeout)
	if timeout == "" {
		timeout = "30s"
	}

	fmt.Print("Headless (true/false, default false for creation): ")
	headlessStr, _ := reader.ReadString('\n')
	headlessStr = strings.TrimSpace(headlessStr)
	headless := headlessStr == "true"

	script := types.VoltScript{
		Name: name,
		Config: types.GlobalConfig{
			Headless: headless,
			Timeout:  timeout,
		},
		Steps: []types.Step{},
	}

	err := playwright.Install()
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
		Headless: playwright.Bool(false),
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

	err = page.ExposeFunction("recordAction", func(args ...interface{}) interface{} {
		if len(args) < 3 {
			return nil
		}
		action, _ := args[0].(string)
		selector, _ := args[1].(string)
		value, _ := args[2].(string)

		step := types.Step{
			Action: action,
		}
		switch action {
		case "click":
			step.Selector = selector
		case "type":
			step.Selector = selector
			step.Value = value
		case "press_key":
			step.Selector = selector
			step.Value = value
		}
		script.Steps = append(script.Steps, step)
		fmt.Printf("✔ Recorded: %s %s %s\n", action, selector, value)
		return nil
	})

	initScript := `
		let inputTimeout = null;
		let lastInputTarget = null;

		function getSelector(el) {
			if (el.id) return '#' + el.id;
			if (el.name) return '[name="' + el.name + '"]';
			if (el.placeholder) return '[placeholder="' + el.placeholder + '"]';
			let path = [];
			while (el && el.nodeType === Node.ELEMENT_NODE) {
				let selector = el.nodeName.toLowerCase();
				if (el.id) {
					selector += '#' + el.id;
					path.unshift(selector);
					break;
				}
				let sibling = el;
				let nth = 1;
				while (sibling = sibling.previousElementSibling) {
					if (sibling.nodeName.toLowerCase() == selector) nth++;
				}
				if (nth != 1) selector += ":nth-of-type("+nth+")";
				path.unshift(selector);
				el = el.parentNode;
			}
			return path.join(' > ');
		}

		document.addEventListener('click', e => {
			if (!e.isTrusted) return;
			let sel = getSelector(e.target);
			window.recordAction("click", sel, "");
		}, true);

		document.addEventListener('input', e => {
			if (!e.isTrusted) return;
			let sel = getSelector(e.target);
			lastInputTarget = e.target;
			if (inputTimeout) clearTimeout(inputTimeout);
			inputTimeout = setTimeout(() => {
				window.recordAction("type", sel, e.target.value);
				inputTimeout = null;
			}, 500);
		}, true);

		document.addEventListener('keydown', e => {
			if (!e.isTrusted) return;
			if (e.key === 'Enter') {
				if (inputTimeout) {
					clearTimeout(inputTimeout);
					let sel = getSelector(lastInputTarget);
					window.recordAction("type", sel, lastInputTarget.value);
					inputTimeout = null;
				}
				let sel = getSelector(e.target);
				window.recordAction("press_key", sel, "Enter");
			}
		}, true);
	`
	page.AddInitScript(playwright.Script{Content: playwright.String(initScript)})

	var lastUrl string
	page.On("framenavigated", func(frame playwright.Frame) {
		if frame == page.MainFrame() {
			url := frame.URL()
			if url != "about:blank" && url != lastUrl {
				lastUrl = url
				script.Steps = append(script.Steps, types.Step{
					Action: "navigate",
					URL:    url,
				})
				fmt.Printf("✔ Recorded: navigate %s\n", url)
			}
		}
	})

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Browser opened! Perform your actions in the browser.")
	fmt.Println("Press Enter in this terminal when you are done to save the script.")

	reader.ReadString('\n')

	filename := "generated.yaml"
	if name != "" {
		filename = strings.ReplaceAll(strings.ToLower(name), " ", "_") + ".yaml"
	}
	data, _ := yaml.Marshal(script)
	os.WriteFile(filename, data, 0644)
	fmt.Printf("✔ Script saved to %s\n", filename)
}
