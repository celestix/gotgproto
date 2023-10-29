package parsemode

import (
	"github.com/gotd/td/telegram/message/styling"
)

// ** bold
// ` mono
// __ italic
// ~~ strike
var mdMap = map[rune]string{
	'*': "bold",
	'`': "mono",
	'_': "italic",
	'~': "strike",
	'|': "spoiler",
}

var stylingMap = map[string]func(string) styling.StyledTextOption{
	"bold":    styling.Bold,
	"mono":    styling.Code,
	"italic":  styling.Italic,
	"strike":  styling.Strike,
	"spoiler": styling.Spoiler,
	"plain":   styling.Plain,
}

func StylizeText(s string) []styling.StyledTextOption {
	var a []styling.StyledTextOption
	var trigger string
	var triggered bool
	var tString string
	for _, c := range s {
		t, ok := mdMap[c]
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
