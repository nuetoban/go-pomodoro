package main

import (
	"fmt"
	"time"
)

type BarTemplate struct {
	time   time.Time
	name   string
	finish string
	emoji  string
}

func NewBarTemplate() *BarTemplate {
	t := &BarTemplate{}
	t.time = time.Now()
	t.finish = ` {{ bar . "[" "-" (cycle . "🕛" "🕐" "🕑" "🕒" "🕓" "🕔" "🕕" "🕖" "🕗" "🕘" "🕙" "🕚" ) "." "]"}} {{etime .}} {{percent .}}`
	return t
}

func (t *BarTemplate) SetEmoji(emoji string) *BarTemplate {
	t.emoji = emoji
	return t
}

func (t *BarTemplate) SetName(name string) *BarTemplate {
	t.name = name
	return t
}

func (t *BarTemplate) String() string {
	return fmt.Sprintf(`%s %30s `+
		`{{ red "%s" }} `+
		`{{ bar . "[" "-" (cycle . "🕛" "🕐" "🕑" "🕒" "🕓" "🕔" "🕕" "🕖" "🕗" "🕘" "🕙" "🕚" ) "." "]"}} `+
		`{{etime .}} `+
		`{{percent .}}`, t.time.Format(time.Stamp), t.name, t.emoji)
}
