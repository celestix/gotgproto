package entityhelper

import "github.com/gotd/td/tg"

// EntityRoot is used to create message entities from the input string through its various methods.
type EntityRoot struct {
	String   string
	Offset   int
	Entities tg.MessageEntityClassArray
}

// StartParsing function creates an empty EntityRoot.
func StartParsing() *EntityRoot {
	return &EntityRoot{
		String:   "",
		Offset:   0,
		Entities: tg.MessageEntityClassArray{},
	}
}

// Bold appends the provided string as bold to the entity root.
func (root *EntityRoot) Bold(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityBold{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Italic appends the provided string as italic to the entity root.
func (root *EntityRoot) Italic(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityItalic{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Underline appends the provided string as underline to the entity root.
func (root *EntityRoot) Underline(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityUnderline{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Code appends the provided string as code/mono to the entity root.
func (root *EntityRoot) Code(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityCode{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Strike appends the provided string as strike to the entity root.
func (root *EntityRoot) Strike(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityStrike{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Spoiler appends the provided string as spoiler to the entity root.
func (root *EntityRoot) Spoiler(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntitySpoiler{Offset: root.Offset, Length: len(s)})
	root.Offset = len(s)
	root.String += s
	return root
}

// Link appends the provided link to the entity root.
func (root *EntityRoot) Link(text, url string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityTextURL{Offset: root.Offset, Length: len(text), URL: url})
	root.Offset = len(text)
	root.String += text
	return root
}

func (root *EntityRoot) GetEntities() []tg.MessageEntityClass {
	return root.Entities
}

func (root *EntityRoot) GetString() string {
	return root.String
}
