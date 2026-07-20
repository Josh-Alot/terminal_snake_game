package ui

import "testing"

func TestNewStartMenu_HasFourItems(t *testing.T) {
	menu := NewStartMenu(false)

	items := menu.Items()
	if len(items) != 4 {
		t.Fatalf("expected 4 items, got %d", len(items))
	}

	want := []MenuItem{
		{ID: IDStart, Label: "Start Game", Kind: ItemAction},
		{ID: IDLeaderboard, Label: "Leaderboard", Kind: ItemAction},
		{ID: IDVimToggle, Label: "Vim Mode", Kind: ItemToggle},
		{ID: IDQuit, Label: "Quit", Kind: ItemAction},
	}

	for i, item := range want {
		if items[i] != item {
			t.Errorf("item %d: expected %+v, got %+v", i, item, items[i])
		}
	}

	if menu.SelectedIndex() != 0 {
		t.Errorf("expected initial selected index 0, got %d", menu.SelectedIndex())
	}
}

func TestMenu_MoveDownAdvances(t *testing.T) {
	menu := NewStartMenu(false)

	menu.MoveDown()

	if menu.SelectedIndex() != 1 {
		t.Errorf("expected selected index 1, got %d", menu.SelectedIndex())
	}
}

func TestMenu_MoveUpDecrements(t *testing.T) {
	menu := NewStartMenu(false)
	menu.MoveDown()
	menu.MoveDown()

	menu.MoveUp()

	if menu.SelectedIndex() != 1 {
		t.Errorf("expected selected index 1, got %d", menu.SelectedIndex())
	}
}

func TestMenu_NavigationClampsAtBounds(t *testing.T) {
	menu := NewStartMenu(false)

	menu.MoveUp()
	if menu.SelectedIndex() != 0 {
		t.Errorf("expected selected index to clamp at 0, got %d", menu.SelectedIndex())
	}

	last := len(menu.Items()) - 1
	for i := 0; i < len(menu.Items())+2; i++ {
		menu.MoveDown()
	}
	if menu.SelectedIndex() != last {
		t.Errorf("expected selected index to clamp at %d, got %d", last, menu.SelectedIndex())
	}
}

func TestMenu_ToggleVim(t *testing.T) {
	menu := NewStartMenu(false)

	if menu.IsVimMode() {
		t.Fatal("expected vim mode to start disabled")
	}

	menu.ToggleVim()
	if !menu.IsVimMode() {
		t.Error("expected vim mode enabled after one toggle")
	}

	menu.ToggleVim()
	if menu.IsVimMode() {
		t.Error("expected vim mode disabled after two toggles")
	}
}

func TestNewStartMenu_InitialVimMode(t *testing.T) {
	menu := NewStartMenu(true)

	if !menu.IsVimMode() {
		t.Error("expected vim mode to start enabled")
	}
}

func TestMenu_Selected(t *testing.T) {
	menu := NewStartMenu(false)

	if got := menu.Selected(); got.ID != IDStart {
		t.Errorf("expected selected ID %d, got %d", IDStart, got.ID)
	}

	menu.MoveDown()
	menu.MoveDown()

	if got := menu.Selected(); got.ID != IDVimToggle {
		t.Errorf("expected selected ID %d, got %d", IDVimToggle, got.ID)
	}
}
