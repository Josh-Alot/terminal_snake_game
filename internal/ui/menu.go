package ui

type ItemKind int

const (
	ItemAction ItemKind = iota // leaves the menu and triggers an action
	ItemToggle                 // flips a boolean and stays on screen
)

type MenuItem struct {
	ID    int
	Label string
	Kind  ItemKind
}

// Well-known item IDs for the start menu.
const (
	IDStart = iota
	IDLeaderboard
	IDVimToggle
	IDQuit
)

type Menu struct {
	items    []MenuItem
	selected int
	vimMode  bool
}

// NewStartMenu builds the 4-item start menu. vimMode is the initial state
// of the Vim Mode toggle.
func NewStartMenu(vimMode bool) *Menu {
	return &Menu{
		items: []MenuItem{
			{ID: IDStart, Label: "Start Game", Kind: ItemAction},
			{ID: IDLeaderboard, Label: "Leaderboard", Kind: ItemAction},
			{ID: IDVimToggle, Label: "Vim Mode", Kind: ItemToggle},
			{ID: IDQuit, Label: "Quit", Kind: ItemAction},
		},
		vimMode: vimMode,
	}
}

func (m *Menu) MoveUp() {
	if m.selected > 0 {
		m.selected--
	}
}

func (m *Menu) MoveDown() {
	if m.selected < len(m.items)-1 {
		m.selected++
	}
}

func (m *Menu) SelectedIndex() int {
	return m.selected
}

func (m *Menu) Selected() MenuItem {
	return m.items[m.selected]
}

func (m *Menu) Items() []MenuItem {
	return m.items
}

func (m *Menu) IsVimMode() bool {
	return m.vimMode
}

func (m *Menu) ToggleVim() {
	m.vimMode = !m.vimMode
}
