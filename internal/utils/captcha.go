package utils

import (
	"strings"

	"github.com/playwright-community/playwright-go"
)

func CheckForCaptcha(page playwright.Page) bool {
	url := page.URL()
	if strings.Contains(url, "sorry/index") {
		return true
	}

	selectors := []string{
		"iframe[src*='recaptcha']",
		"iframe[src*='hcaptcha']",
		"iframe[src*='challenges.cloudflare.com']",
		"iframe[src*='datadome.co']",
		"#px-captcha",
	}

	for _, sel := range selectors {
		ct, err := page.Locator(sel).Count()
		if err == nil && ct > 0 {
			return true
		}
	}

	return false
}
