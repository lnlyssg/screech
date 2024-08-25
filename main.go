package main

import (
	_ "embed"
	"hello/config"
	"hello/input"
	"hello/screens"
	"hello/uilib"
	"log"
	"os"
	"runtime/debug"

	"github.com/veandco/go-sdl2/sdl"
)

//go:embed assets/NotoSans_Condensed-SemiBold.ttf
var NotoSans []byte

const (
	fontPath        = "./test.ttf"
	fontSize        = 36
	WinWidth  int32 = 1280
	WinHeight int32 = 720
	CenterX   int32 = WinWidth / 2
	CenterY   int32 = WinHeight / 2
)

type Text struct {
	Content string
	X, Y    int32
}

func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Unhandled error: %v\n", r)
			log.Println("Stack trace:")
			debug.PrintStack()
			os.Exit(-1)
		}
	}()

	config.InitVars()

	if err := uilib.InitSDL(); err != nil {
		panic(err)
	}

	if err := uilib.InitTTF(); err != nil {
		panic(err)
	}

	if err := uilib.InitFont(NotoSans, &config.BodyFont, 24); err != nil {
		panic(err)
	}

	if err := uilib.InitFont(NotoSans, &config.BodyBigFont, 58); err != nil {
		panic(err)
	}

	if err := uilib.InitFont(NotoSans, &config.LongTextFont, 24); err != nil {
		panic(err)
	}

	if err := uilib.InitFont(NotoSans, &config.HeaderFont, 28); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Systems List", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, config.ScreenWidth, config.ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	mainScreen, err := screens.NewMainScreen(renderer)
	if err != nil {
		panic(err)
	}

	scrapingScreen, err := screens.NewScrapingScreen(renderer)
	if err != nil {
		panic(err)
	}

	screensMap := map[string]func(){
		"main_screen":     mainScreen.Draw,
		"scraping_screen": scrapingScreen.Draw,
	}

	inputHandlers := map[string]func(input.InputEvent){
		"main_screen":     mainScreen.HandleInput,
		"scraping_screen": scrapingScreen.HandleInput,
	}

	input.StartListening()

	running := true
	for running {

		for {
			event := sdl.PollEvent()
			if event == nil {
				break
			}

			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		select {
		case inputEvent := <-input.InputChannel:
			if handler, ok := inputHandlers[config.CurrentScreen]; ok {
				handler(inputEvent)
			}
		default:
			// No event received
		}

		if drawFunc, ok := screensMap[config.CurrentScreen]; ok {
			drawFunc()
		}

		sdl.Delay(16)
	}
}
