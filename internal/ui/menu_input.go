package ui

type MenuEvent int

const (
	MenuNone MenuEvent = iota
	MenuUp
	MenuDown
	MenuSelect
	MenuQuit
)

// MapMenuKey maps a key token produced by input.ReadKey to a menu event.
// Vim mode is irrelevant on the menu: j/k always work.
func MapMenuKey(key string) MenuEvent {
	switch key {
	case "\x1b[A", "k":
		return MenuUp
	case "\x1b[B", "j":
		return MenuDown
	case "\r", " ":
		return MenuSelect
	case "q", "\x03":
		return MenuQuit
	default:
		return MenuNone
	}
}
