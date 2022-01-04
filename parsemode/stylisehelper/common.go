package stylisehelper

import "github.com/gotd/td/telegram/message/styling"

// StyledTextRoot is used to create an array of styling.StyledTextOption from the input string through its various methods.
type StyledTextRoot struct {
	StoArray []styling.StyledTextOption
}

// Start function creates an StyledTextRoot with the provided styling.StyledTextOption.
func Start(style styling.StyledTextOption) *StyledTextRoot {
	return &StyledTextRoot{StoArray: []styling.StyledTextOption{style}}
}

// Bold appends the provided string as bold to the styled text root.
func (sh *StyledTextRoot) Bold(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Bold(s))
	return sh
}

// Code appends the provided string as code/mono to the styled text root.
func (sh *StyledTextRoot) Code(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Code(s))
	return sh
}

// Strike appends the provided string as strike to the styled text root.
func (sh *StyledTextRoot) Strike(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Strike(s))
	return sh
}

// Underline appends the provided string as underline to the styled text root.
func (sh *StyledTextRoot) Underline(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Underline(s))
	return sh
}

// Italic appends the provided string as italic to the styled text root.
func (sh *StyledTextRoot) Italic(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Italic(s))
	return sh
}

// Plain appends he provided string as plain text to the styled text root.
func (sh *StyledTextRoot) Plain(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Plain(s))
	return sh
}

// Link appends the provided link to the styled text root.
func (sh *StyledTextRoot) Link(text, url string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.TextURL(text, url))
	return sh
}

// Spoiler appends the provided string as spoiler to the styled text root.
func (sh *StyledTextRoot) Spoiler(s string) *StyledTextRoot {
	sh.StoArray = append(sh.StoArray, styling.Spoiler(s))
	return sh
}
