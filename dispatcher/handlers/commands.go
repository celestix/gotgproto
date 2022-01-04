package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"strings"
)

type Command struct {
	Prefix        []rune
	Name          string
	Callback      CallbackResponse
	UpdateFilters filters.UpdateFilter
}

var DefaultPrefix = []rune{'!', '/'}

func NewCommand(name string, response CallbackResponse) Command {
	return Command{
		Name:          name,
		Callback:      response,
		Prefix:        DefaultPrefix,
		UpdateFilters: nil,
	}
}

func (c Command) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage == nil || u.EffectiveMessage.Out || len(u.EffectiveMessage.Message) == 0 {
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
