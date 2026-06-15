package utils

import (
	"math/rand"
	"time"

	"github.com/playwright-community/playwright-go"
)

const injectMousePointerScript = `
	() => {
		if (window._voltMouseInjected) return;
		window._voltMouseInjected = true;
		window._voltMouseX = window.innerWidth / 2;
		window._voltMouseY = window.innerHeight / 2;

		const box = document.createElement('div');
		box.id = 'volt-human-cursor';
		box.style.width = '20px';
		box.style.height = '20px';
		box.style.position = 'absolute';
		box.style.top = '0px';
		box.style.left = '0px';
		box.style.zIndex = '9999999999';
		box.style.pointerEvents = 'none';
		box.style.backgroundImage = 'url("data:image/svg+xml;utf8,<svg xmlns=\\"http://www.w3.org/2000/svg\\" width=\\"24\\" height=\\"24\\" viewBox=\\"0 0 24 24\\" fill=\\"none\\" stroke=\\"black\\" stroke-width=\\"2\\" stroke-linecap=\\"round\\" stroke-linejoin=\\"round\\"><path d=\\"m3 3 7.07 16.97 2.51-7.39 7.39-2.51L3 3z\\"/><path d=\\"m13 13 6 6\\"/></svg>")';
		box.style.backgroundSize = 'contain';
		box.style.transition = 'top 0.05s linear, left 0.05s linear';
		document.body.appendChild(box);

		document.addEventListener('mousemove', event => {
			window._voltMouseX = event.pageX;
			window._voltMouseY = event.pageY;
			box.style.left = event.pageX + 'px';
			box.style.top = event.pageY + 'px';
		});

		document.addEventListener('mousedown', event => {
			box.style.transform = 'scale(0.8)';
		});

		document.addEventListener('mouseup', event => {
			box.style.transform = 'scale(1)';
		});
	}
`

func InitMousePointer(page playwright.Page) error {
	_, err := page.Evaluate(injectMousePointerScript)
	return err
}

func HumanizeMouse(page playwright.Page, selector string) error {
	box, err := page.Locator(selector).BoundingBox()
	if err != nil {
		return err
	}

	var startX, startY float64
	xRaw, errX := page.Evaluate("window._voltMouseX || window.innerWidth / 2")
	yRaw, errY := page.Evaluate("window._voltMouseY || window.innerHeight / 2")

	if errX == nil && errY == nil {
		if xf, ok := xRaw.(float64); ok {
			startX = xf
		} else if xi, ok := xRaw.(int); ok {
			startX = float64(xi)
		}

		if yf, ok := yRaw.(float64); ok {
			startY = yf
		} else if yi, ok := yRaw.(int); ok {
			startY = float64(yi)
		}
	}

	targetX := float64(box.X) + float64(box.Width)/2.0 + (rand.Float64()*10 - 5)
	targetY := float64(box.Y) + float64(box.Height)/2.0 + (rand.Float64()*10 - 5)

	midX := (startX + targetX) / 2
	midY := (startY + targetY) / 2

	dx := targetX - startX
	dy := targetY - startY

	offsetFactor := (rand.Float64() * 0.6) - 0.3

	ctrlX := midX - (dy * offsetFactor)
	ctrlY := midY + (dx * offsetFactor)

	segments := 10

	for i := 1; i <= segments; i++ {
		t := float64(i) / float64(segments)

		x := (1-t)*(1-t)*startX + 2*(1-t)*t*ctrlX + t*t*targetX
		y := (1-t)*(1-t)*startY + 2*(1-t)*t*ctrlY + t*t*targetY

		steps := 3
		if i == segments {
			steps = 6
		}

		if err := page.Mouse().Move(x, y, playwright.MouseMoveOptions{Steps: playwright.Int(steps)}); err != nil {
			return err
		}

		time.Sleep(time.Duration(rand.Intn(10)+10) * time.Millisecond)
	}

	time.Sleep(time.Duration(rand.Intn(150)+50) * time.Millisecond)

	return nil
}
