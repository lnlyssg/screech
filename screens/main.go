package screens

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anibaldeboni/screech/components"
	"github.com/anibaldeboni/screech/config"
	"github.com/anibaldeboni/screech/input"
	"github.com/anibaldeboni/screech/uilib"

	"github.com/veandco/go-sdl2/sdl"
)

type MainScreen struct {
	renderer      *sdl.Renderer
	listComponent *components.List
	initialized   bool
}

func NewMainScreen(renderer *sdl.Renderer) (*MainScreen, error) {
	return &MainScreen{
		renderer: renderer,
		listComponent: components.NewList(renderer, 18, func(index int, item components.Item) string {
			return fmt.Sprintf("%d. %s", index+1, item.Text)
		}),
	}, nil
}

func (m *MainScreen) InitMain() {
	if m.initialized {
		return
	}
	systems := romDirsToList(listEmulatorDirs())
	m.listComponent.SetItems(systems)
	config.CurrentPlatform = systems[0].ID
	m.initialized = true
}

func romDirsToList(romDirs []RomDir) []components.Item {
	var items []components.Item
	for _, romDir := range romDirs {
		items = append(items, components.Item{
			Text:  romDir.Name,
			ID:    romDir.Name,
			Value: romDir.Path,
		})
	}
	return items
}

func (m *MainScreen) HandleInput(event input.InputEvent) {
	switch event.KeyCode {
	case "DOWN":
		m.listComponent.ScrollDown()
	case "UP":
		m.listComponent.ScrollUp()
	case "B":
		os.Exit(0)
	case "A":
		if len(m.listComponent.GetItems()) == 0 {
			return
		}
		config.CurrentScreen = "scraping_screen"
		config.CurrentSystem = m.SelectedSystem().ID
		SetScraping()
		m.initialized = false
	}
	m.updateLogo()
}

func (m *MainScreen) updateLogo() {
	selectedSystem := m.SelectedSystem()
	config.CurrentSystem = selectedSystem.ID
	uilib.RenderImage(m.renderer, "assets/logos/"+selectedSystem.ID+".png")
}

func (m *MainScreen) SelectedSystem() components.Item {
	return m.listComponent.GetItems()[m.listComponent.GetSelectedIndex()]
}

func (m *MainScreen) Draw() {
	m.InitMain()

	m.renderer.SetDrawColor(0, 0, 0, 255) // Background color
	m.renderer.Clear()

	uilib.RenderTexture(m.renderer, config.UiBackground, "Q2", "Q4")

	// Draw the title
	uilib.DrawText(m.renderer, "Systems", sdl.Point{X: 25, Y: 25}, config.Colors.PRIMARY, config.HeaderFont)

	m.listComponent.Draw(config.Colors.WHITE, config.Colors.SECONDARY)

	uilib.RenderTexture(m.renderer, config.UiControls, "Q3", "Q4")
	m.updateLogo()

	m.renderer.Present()
}

type RomDir struct {
	Name string
	Path string
}

func listEmulatorDirs() []RomDir {
	romsDir := config.Roms
	dirEntries, err := os.ReadDir(romsDir)
	if err != nil {
		panic(fmt.Errorf("error reading dir %s: %w", romsDir, err))
	}

	var dirs []RomDir
	for _, entry := range dirEntries {
		if entry.IsDir() {
			dirPath := filepath.Join(romsDir, entry.Name())
			dirFiles, err := os.ReadDir(dirPath)
			if err != nil {
				panic(fmt.Errorf("error reading dir %s: %w", dirPath, err))
			}
			if len(dirFiles) > 0 {
				dirs = append(dirs, RomDir{
					Name: entry.Name(),
					Path: dirPath,
				})
			}
		}
	}

	return dirs
}
