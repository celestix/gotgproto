package handlers

import (
	"strings"

	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
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
	if u.EffectiveMessage == nil || u.EffectiveMessage.Message == "" {
		return nil
	}
	if !c.Outgoing && u.EffectiveMessage.Out {
		return nil
	}
	if c.UpdateFilters != nil && !c.UpdateFilters(u) {
		return nil
	}
	args := strings.Fields(strings.ToLower(u.EffectiveMessage.Message))
	for _, prefix := range c.Prefix {
		if args[0][0] == byte(prefix) {
			if args[0][1:] == c.Name {
				return c.Callback(ctx, u)
			} else if split := strings.Split(args[0][1:], "@"); split[0] == c.Name {
				if split[1] == strings.ToLower(ctx.Self.Username) {
					return c.Callback(ctx, u)
				}
			}
		}
	}
	return nil
}
