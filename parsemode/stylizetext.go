package parsemode

import (
	"strings"

	"github.com/gotd/td/telegram/message/styling"
)

var mdMap = map[rune]string{
	'*': "bold",
	'`': "mono",
	'_': "italic",
	'~': "strike",
	'|': "spoiler",
}

var htmlMap = map[string]string{
	"<b>":                         "bold",
	"<i>":                         "italic",
	"<code>":                      "mono",
	"<s>":                         "strike",
	"<span class=\"tg-spoiler\">": "spoiler",
}

var stylingMap = map[string]func(string) styling.StyledTextOption{
	"bold":    styling.Bold,
	"mono":    styling.Code,
	"italic":  styling.Italic,
	"strike":  styling.Strike,
	"spoiler": styling.Spoiler,
	"plain":   styling.Plain,
}

// StylizeText converts a formatted string into a slice of styled text options.
//
// Parameters:
// - s: The input string containing the text to be styled.
// - mode: (Optional) Defines the parsing mode. Can be "markdown" or "html".
//   - "markdown": Uses markdown syntax (e.g., *bold*, _italic_).
//   - "html": Uses HTML tags (e.g., <b>bold</b>, <i>italic</i>).
//   - If the mode is not specified or is different from "html", the default is "markdown".
//
// Returns:
// - []styling.StyledTextOption: A slice of styled text options that can be used with the styling package.
func StylizeText(s string, mode ...string) []styling.StyledTextOption {
	var a []styling.StyledTextOption
	var trigger string
	var triggered bool
	var tString string

	// Determine mode: "markdown" (default) or "html"
	selectedMode := "markdown"
	if len(mode) > 0 && mode[0] == "html" {
		selectedMode = "html"
	}

	for i := 0; i < len(s); i++ {
		c := rune(s[i])
		var t string
		var ok bool

		if selectedMode == "markdown" {
			t, ok = mdMap[c]
		} else if selectedMode == "html" {
			for tag, style := range htmlMap {
				if strings.HasPrefix(s[i:], tag) {
					t = style
					i += len(tag) - 1 // Move index to end of tag
					ok = true
					break
				}
			}
		}

		if !ok && !triggered {
			trigger = "plain"
			tString += string(c)
		}
		if triggered {
			if ok {
				if t == trigger {
					a = append(a, stylingMap[trigger](tString))
					tString = ""
					triggered = false
				}
				continue
			}
			tString += string(c)
			continue
		}
		if ok {
			if trigger == "plain" {
				a = append(a, stylingMap[trigger](tString))
			}
			tString = ""
			trigger = t
			triggered = true
		}
	}

	if tString != "" { // without this check, the last style is not added
		a = append(a, stylingMap[trigger](tString))
	}
	return a
}
