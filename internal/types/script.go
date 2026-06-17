package types

type VoltScript struct {
	Name   string            `yaml:"name" json:"name" toml:"name"`
	Config GlobalConfig      `yaml:"config" json:"config" toml:"config"`
	Vars   map[string]string `yaml:"vars,omitempty" json:"vars,omitempty" toml:"vars,omitempty"`
	Steps  []Step            `yaml:"steps" json:"steps" toml:"steps"`
}

type GlobalConfig struct {
	Headless       bool     `yaml:"headless" json:"headless" toml:"headless"`
	SlowMo         string   `yaml:"slow_mo,omitempty" json:"slow_mo,omitempty" toml:"slow_mo,omitempty"`
	Timeout        string   `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Output         string   `yaml:"output,omitempty" json:"output,omitempty" toml:"output,omitempty"`
	Browsers       []string `yaml:"browsers,omitempty" json:"browsers,omitempty" toml:"browsers,omitempty"`
	ErrorIfCaptcha bool     `yaml:"errorIfCaptcha,omitempty" json:"errorIfCaptcha,omitempty" toml:"errorIfCaptcha,omitempty"`
	Humanize       bool     `yaml:"humanize,omitempty" json:"humanize,omitempty" toml:"humanize,omitempty"`
}

type Step struct {
	Action    string            `yaml:"action" json:"action" toml:"action"`
	Name      string            `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	URL       string            `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	Selector  string            `yaml:"selector,omitempty" json:"selector,omitempty" toml:"selector,omitempty"`
	Value     string            `yaml:"value,omitempty" json:"value,omitempty" toml:"value,omitempty"`
	Duration  int               `yaml:"duration,omitempty" json:"duration,omitempty" toml:"duration,omitempty"`
	Key       string            `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	File      string            `yaml:"file,omitempty" json:"file,omitempty" toml:"file,omitempty"`
	As        string            `yaml:"as,omitempty" json:"as,omitempty" toml:"as,omitempty"`
	Attribute string            `yaml:"attribute,omitempty" json:"attribute,omitempty" toml:"attribute,omitempty"`
	Condition string            `yaml:"condition,omitempty" json:"condition,omitempty" toml:"condition,omitempty"`
	Then      []Step            `yaml:"then,omitempty" json:"then,omitempty" toml:"then,omitempty"`
	Else      []Step            `yaml:"else,omitempty" json:"else,omitempty" toml:"else,omitempty"`
	ForEach   []string          `yaml:"for_each,omitempty" json:"for_each,omitempty" toml:"for_each,omitempty"`
	Iterator  string            `yaml:"iterator,omitempty" json:"iterator,omitempty" toml:"iterator,omitempty"`
	Message   string            `yaml:"message,omitempty" json:"message,omitempty" toml:"message,omitempty"`
	Scrape    *ScrapeConfig     `yaml:"scrape,omitempty" json:"scrape,omitempty" toml:"scrape,omitempty"`
	Assert    *AssertConfig     `yaml:"assert,omitempty" json:"assert,omitempty" toml:"assert,omitempty"`
	Headers   map[string]string `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Pathname  string            `yaml:"pathname,omitempty" json:"pathname,omitempty" toml:"pathname,omitempty"`
	Position  *Position         `yaml:"position,omitempty" json:"position,omitempty" toml:"position,omitempty"`

	From int    `yaml:"from,omitempty" json:"from,omitempty" toml:"from,omitempty"`
	To   int    `yaml:"to,omitempty" json:"to,omitempty" toml:"to,omitempty"`
	Do   []Step `yaml:"do,omitempty" json:"do,omitempty" toml:"do,omitempty"`
}

type ScrapeConfig struct {
	Parent string            `yaml:"parent,omitempty" json:"parent,omitempty" toml:"parent,omitempty"`
	Fields map[string]string `yaml:"fields" json:"fields" toml:"fields"`
}

type Position struct {
	X        int  `yaml:"x,omitempty" json:"x,omitempty" toml:"x,omitempty"`
	Y        int  `yaml:"y,omitempty" json:"y,omitempty" toml:"y,omitempty"`
	Width    int  `yaml:"width,omitempty" json:"width,omitempty" toml:"width,omitempty"`
	Height   int  `yaml:"height,omitempty" json:"height,omitempty" toml:"height,omitempty"`
	FullPage bool `yaml:"full_page,omitempty" json:"full_page,omitempty" toml:"full_page,omitempty"`
}

type AssertConfig struct {
	Selector string `yaml:"selector,omitempty" json:"selector,omitempty" toml:"selector,omitempty"`
	Equals   string `yaml:"equals,omitempty" json:"equals,omitempty" toml:"equals,omitempty"`
	Contains string `yaml:"contains,omitempty" json:"contains,omitempty" toml:"contains,omitempty"`
	Eval     string `yaml:"eval,omitempty" json:"eval,omitempty" toml:"eval,omitempty"`
}
