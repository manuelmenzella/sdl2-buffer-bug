package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func run() error {
	// Create a Surface of size 31x24.
	// The width and pixel format are chosen such that their product
	// is not a multiple of four, which causes SDL to add padding to
	// the underlying buffer.
	// https://github.com/libsdl-org/SDL/blob/c59d4dcd38c382a1e9b69b053756f1139a861574/src/video/SDL_surface.c#L51
	surface, err := sdl.CreateRGBSurfaceWithFormat(
		0, 31 /* width */, 24 /* height */, 24, sdl.PIXELFORMAT_RGB24)
	if err != nil {
		return err
	}
	defer surface.Free()

	// The Surface pitch (bytes per row) is 31*3, but gets rounded up to a
	// multiple of four -> 96. The underlying buffer should then be 24*96 = 2304 bytes.
	// However, Surface::Pixels miscalculates this as w*h*3, or 31*24*3 = 2232 bytes.
	pixels := surface.Pixels()

	// Setting the last byte to 0xff should set the last pixel
	// (bottom-right, row 23, column 30) to blue. However, with the current
	// implementation of the SDL wrapper, this actually sets the byte corresponding
	// to the blue channel for the pixel at row 23, column 7 -- ((2231 mod 96) - 2) / 3.
	pixels[len(pixels)-1] = 0xff

	// All code below this line is strictly boiler plate to blit the small surface into
	// a window large enough to see the effect.
	window, err := sdl.CreateWindow(
		"window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return err
	}
	defer window.Destroy()

	windowSurface, err := window.GetSurface()
	if err != nil {
		return err
	}

	if err = surface.BlitScaled(nil, windowSurface, nil); err != nil {
		return nil
	}

	err = window.UpdateSurface()
	if err != nil {
		return err
	}

	for {
		event := sdl.WaitEvent()
		switch event := event.(type) {
		case *sdl.WindowEvent:
			if event.Event == sdl.WINDOWEVENT_CLOSE {
				return nil
			}
		case *sdl.KeyboardEvent:
			if event.Keysym.Scancode == sdl.GetScancodeFromKey(sdl.K_ESCAPE) {
				return nil
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
