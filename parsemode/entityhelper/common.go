package entityhelper

import "github.com/gotd/td/tg"

type EntityRoot struct {
	String   string
	Offset   int
	Entities tg.MessageEntityClassArray
}

func StartParsing() *EntityRoot {
	return &EntityRoot{
		String:   "",
		Offset:   0,
		Entities: tg.MessageEntityClassArray{},
	}
}

func (root *EntityRoot) Bold(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityBold{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) Italic(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityItalic{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) Underline(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityUnderline{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) Code(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityCode{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) Strike(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityStrike{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) Spoiler(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntitySpoiler{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

func (root *EntityRoot) GetEntities() []tg.MessageEntityClass {
	return root.Entities
}

func (root *EntityRoot) GetString() string {
	return root.String
}
