# Volt

Volt is a universal web automation, scripting, testing, and Robotic Process Automation (RPA) tool written in Go. It provides the power of Playwright through simple, lightweight, and portable YAML-TOML-JSON configuration files, allowing for a 100% No-Code/Low-Code experience.

## Features

- **File-Based Automation**: Define complex browser automation workflows using readable YAML-TOML-JSON syntax without writing code.
- **Powered by Playwright**: Reliable and fast cross-browser automation engine.
- **Headless and Headful Modes**: Run scripts invisibly in the background for CI/CD pipelines, or visibly with an optional slow-motion mode for local debugging.
- **Standalone Binary**: Compiled in Go, requiring no Node.js environment or `node_modules`.
- **Static Validation**: Built-in syntax and semantic checker to validate scripts before execution.

## Installation

Ensure you have Go installed (version 1.22 or higher).

Clone the repository and build the project:

```bash
git clone https://github.com/octarahq/volt.git
cd volt
go build -o volt main.go
```

To install the required Playwright browsers, run the command below on your first setup:

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
```

## Usage

Volt provides a command-line interface to interact with your scripts.

### Run a Script

Execute an automation script:

```bash
volt run path/to/script.yaml
```

### Validate a Script

Check the syntax and structure of your YAML file without running the browser:

```bash
volt check path/to/script.yaml
```

### Schema Validation in Editors

You don't need to configure your editor settings to get autocompletion and validation. Simply add the schema link at the very top of your script file.

**For YAML files** (`script.yaml`):

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/octarahq/volt/main/volt-schema.json
# ...
```

**For JSON files** (`script.json`):

```json
{
  "$schema": "https://raw.githubusercontent.com/octarahq/volt/main/volt-schema.json"
  // ...
}
```

**For TOML files** (`script.toml`):

```toml
#:schema https://raw.githubusercontent.com/octarahq/volt/main/scripts/volt-schema.json
# ...
```

## Script Structure

A Volt script is composed of global configurations, variables, and sequential steps.

### Supported Configurations

The `config` block allows you to define global settings for your automation script. The following options are supported:

- `headless` (boolean): Run the browser in headless mode (invisible).
- `slow_mo` (string): Delay between Playwright operations. Useful for local debugging (e.g., `"500ms"`).
- `timeout` (string): Maximum execution time for the automation script.
- `output` (string): Define an output directory path for saved data (like screenshots or scraped data).
- `browsers` (list of strings): A list of browsers to execute the automation on (e.g., `["chromium", "firefox", "webkit"]`). The script will be run for each browser specified.

Example `script.yaml`:

```yaml
name: "Example Automation"
config:
  headless: false
  slow_mo: "500ms"
  browsers:
    - chromium
    - firefox
steps:
  - action: "navigate"
    url: "https://example.com"

  - action: "type"
    selector: "input[name='q']"
    value: "Volt automation"

  - action: "press_key"
    selector: "input[name='q']"
    value: "Enter"

  - action: "click"
    selector: "a.result-link"
```

## Supported Actions

Here is the implementation status of the actions planned for Volt:

### Interactions

- [x] `navigate`: Go to a specific URL (`url`).
- [x] `click`: Click on an element targeting a CSS selector (`selector`).
- [x] `hover`: Hover over an element (`selector`).
- [x] `type`: Type text into an input field (`selector`, `value`).
- [x] `press_key`: Press a specific keyboard key on an element (`selector`, `value`).
- [x] `check` / `uncheck`: Toggle checkboxes.
- [x] `select`: Select options from dropdowns.
- [x] `upload`: Upload files.
- [x] `scroll`: Scroll the page (top, bottom, or to a specific selector).

### Variables & RPA

- [x] `store_value`: Store constant value into a variable.
- [x] `store_text`: Store text content of an element into a variable.
- [x] `store_attribute`: Store an attribute of an element into a variable.
- [x] `store_eval`: Store the result of a JS evaluation into a variable.

### Scraping

- [x] `scrape`: Extract structured list data into JSON/CSV.

### Logic and Flows

- [x] `if`: Conditional execution block.
- [x] `loop`: Indexed numerical loop.
- [ ] `for_each`: Iterate over a list.
- [x] `log`: Console logging.

### Waiting & Delays

- [x] `wait`: Fixed time delay.
- [x] `wait_visible` / `wait_hidden`: Wait for an element's visibility state.

### Assertions (Test Mode)

- [ ] `assert_visible` / `assert_not_visible`: Assert element visibility.
- [ ] `assert_text`: Assert text content matches.
- [ ] `assert_eval`: Assert the result of a JS evaluation.

### System

- [x] `screenshot`: Capture a screenshot.
- [x] `clear_cookies`: Clear browser cookies.
- [x] `add_header`: Add custom HTTP header.
- [x] `set_header`: Set custom HTTP header.
- [x] `remove_header`: Remove custom HTTP header.
