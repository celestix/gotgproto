package handlers

import (
	"strings"

	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
)

// Command handler is executed when the update consists of tg.Message provided it is a command and satisfies all the conditions.
type Command struct {
	Prefix        []rune
	Name          string
	Callback      CallbackResponse
	Outgoing      bool
	UpdateFilters filters.UpdateFilter
}

// DefaultPrefix is the global variable consisting all the prefixes which will trigger the command.
var DefaultPrefix = []rune{'!', '/'}

// NewCommand creates a new Command handler with default fields, bound to call its response
func NewCommand(name string, response CallbackResponse) Command {
	return Command{
		Name:     name,
		Callback: response,
		Prefix:   DefaultPrefix,
		Outgoing: true,
	}
}

func (c Command) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	m := u.EffectiveMessage
	if m == nil || m.Text == "" {
		return nil
	}
	if !c.Outgoing && m.Out {
		return nil
	}
	if c.UpdateFilters != nil && !c.UpdateFilters(u) {
		return nil
	}
	arg := strings.ToLower(
		strings.Fields(m.Text)[0],
	)
	for _, prefix := range c.Prefix {
		if arg[0] == byte(prefix) {
			if arg[1:] == c.Name {
				return c.Callback(ctx, u)
			} else if split := strings.Split(arg[1:], "@"); split[0] == c.Name {
				if split[1] == strings.ToLower(ctx.Self.Username) {
					return c.Callback(ctx, u)
				}
			}
		}
	}
	return nil
}
