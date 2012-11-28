package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"unsafe"
)

var (
	HEIGHT = 800
	WIDTH  = HEIGHT
	INTER  = HEIGHT / 20
	DEC    = INTER / 2
	CIRCLE = DEC / 2
	STONE  = CIRCLE + DEC
)

func clean_side(gc *gdk.GdkGC, pixmap *gdk.GdkPixmap, x1, y1, x2, y2 int) {
	gc.SetRgbFgColor(gdk.Color("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true,
		x1,
		y1,
		x2,
		y2)
}

func draw_square(gc *gdk.GdkGC, pixmap *gdk.GdkPixmap, x, y int) {
	gc.SetRgbFgColor(gdk.Color("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true,
		x*INTER+DEC,
		y*INTER+DEC,
		x*INTER+INTER+DEC,
		y*INTER+INTER+DEC)
	gc.SetRgbFgColor(gdk.Color("black"))
	pixmap.GetDrawable().DrawLine(gc,
		x*INTER+INTER/2+DEC,
		y*INTER+DEC,
		x*INTER+INTER/2+DEC,
		y*INTER+INTER+DEC)
	pixmap.GetDrawable().DrawLine(gc,
		x*INTER+DEC,
		y*INTER+INTER/2+DEC,
		x*INTER+INTER+DEC,
		y*INTER+INTER/2+DEC)
	if x == 0 {
		clean_side(gc, pixmap, 0, 0, INTER, -1) // LEFT
	} else if x == 18 {
		clean_side(gc, pixmap, HEIGHT-INTER+1, 0, HEIGHT, -1) // RIGHT
	}
	if y == 0 {
		clean_side(gc, pixmap, 0, 0, -1, INTER) // TOP
	} else if y == 18 {
		clean_side(gc, pixmap, 0, HEIGHT-INTER+1, -1, -1) // BOT
	}
}

func display_init_grid(gc *gdk.GdkGC, pixmap *gdk.GdkPixmap) {
	for x := 0; x < 19; x++ {
		for y := 0; y < 19; y++ {
			draw_square(gc, pixmap, x, y)
		}
	}
	gc.SetRgbFgColor(gdk.Color("black"))
	pixmap.GetDrawable().DrawArc(gc, true, (4*INTER)-(CIRCLE/2), (4*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (10*INTER)-(CIRCLE/2), (4*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (16*INTER)-(CIRCLE/2), (4*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (4*INTER)-(CIRCLE/2), (10*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (10*INTER)-(CIRCLE/2), (10*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (16*INTER)-(CIRCLE/2), (10*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (4*INTER)-(CIRCLE/2), (16*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (10*INTER)-(CIRCLE/2), (16*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
	pixmap.GetDrawable().DrawArc(gc, true, (16*INTER)-(CIRCLE/2), (16*INTER)-(CIRCLE/2),
		CIRCLE, CIRCLE, 0, 64*360)
}

func board_display(game Gomoku) {
	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("Gomoku")
	window.Connect("destroy", func() {
		println("got destroy!")
		gtk.MainQuit()
	})

	var gdkwin *gdk.GdkWindow
	var pixmap *gdk.GdkPixmap
	var gc *gdk.GdkGC
	var player int
	player = 1

	vbox := gtk.VBox(true, 0)
	drawingarea := gtk.DrawingArea()

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.GtkAllocation
		drawingarea.GetAllocation(&allocation)
		pixmap = gdk.Pixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.GC(pixmap.GetDrawable())
		display_init_grid(gc, pixmap)
	})

	drawingarea.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		if gdkwin == nil {
			gdkwin = drawingarea.GetWindow()
		}
		arg := ctx.Args(0)
		mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
		var mt gdk.GdkModifierType
		var x, y int
		if mev.IsHint != 0 {
			gdkwin.GetPointer(&x, &y, &mt)
		} else {
			x, y = int(mev.X), int(mev.Y)
		}
		vic, err := game.Play(((x - INTER/2) / INTER), ((y - INTER/2) / INTER))
		if err != nil {
			return
		}
		if player == 1 {
			gc.SetRgbFgColor(gdk.Color("black"))
			player = 2
		} else {
			gc.SetRgbFgColor(gdk.Color("white"))
			player = 1
		}
		x = ((x-INTER/2)/INTER)*INTER + INTER
		y = ((y-INTER/2)/INTER)*INTER + INTER
		pixmap.GetDrawable().DrawArc(gc, true, x-(STONE/2), y-(STONE/2), STONE, STONE, 0, 64*360)
		if vic != 0 {
			fmt.Println("Player", vic, "win")
		}
		drawingarea.GetWindow().Invalidate(nil, false)
	})

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	drawingarea.SetEvents(int(gdk.GDK_POINTER_MOTION_MASK | gdk.GDK_POINTER_MOTION_HINT_MASK | gdk.GDK_BUTTON_PRESS_MASK))
	vbox.Add(drawingarea)

	window.Add(vbox)
	window.SetSizeRequest(WIDTH, HEIGHT)
	window.ShowAll()
	gtk.Main()
}

func game_mode() int {
	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("Gomoku")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	var mode int

	vbox := gtk.VBox(true, 0)
	pvp := gtk.ButtonWithLabel("Player Vs Player")
	pvp.Clicked(func() {
		mode = 0
		window.Destroy()
	})
	pvai := gtk.Label("Player Vs Ai")
	vbox.Add(pvp)
	vbox.Add(pvai)
	pvp.Show()
	window.Add(vbox)
	window.ShowAll()
	gtk.Main()
	return mode
}