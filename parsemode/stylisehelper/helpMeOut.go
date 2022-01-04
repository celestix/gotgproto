package stylisehelper

import "github.com/gotd/td/telegram/message/styling"

type StyliseHelper struct {
	StoArray []styling.StyledTextOption
}

func Start(style styling.StyledTextOption) *StyliseHelper {
	return &StyliseHelper{StoArray: []styling.StyledTextOption{style}}
}

func (sh *StyliseHelper) Bold(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Bold(s))
	return sh
}

func (sh *StyliseHelper) Code(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Code(s))
	return sh
}

func (sh *StyliseHelper) Strike(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Strike(s))
	return sh
}

func (sh *StyliseHelper) Underline(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Underline(s))
	return sh
}

func (sh *StyliseHelper) Italic(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Italic(s))
	return sh
}

func (sh *StyliseHelper) Plain(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Plain(s))
	return sh
}

func (sh *StyliseHelper) Link(text, url string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.TextURL(text, url))
	return sh
}

func (sh *StyliseHelper) Spoiler(s string) *StyliseHelper {
	sh.StoArray = append(sh.StoArray, styling.Spoiler(s))
	return sh
}
