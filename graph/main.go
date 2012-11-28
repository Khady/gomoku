package main

import (
	"os"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"unsafe"
)

var (
	HEIGHT int = 800
	WIDTH int = HEIGHT
	INTER int = HEIGHT / 20
	DEC int = INTER - 5
)

type point struct {
	x int
	y int
}

func init_grid(gc *gdk.GdkGC, pixmap *gdk.GdkPixmap) {
	gc.SetRgbFgColor(gdk.Color("gray"))
	pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	gc.SetRgbFgColor(gdk.Color("black"))
	gc.SetRgbBgColor(gdk.Color("white"))
	for i := 0; i < 19; i++ {
		pixmap.GetDrawable().DrawLine(gc, (i * INTER) + DEC, DEC, (i * INTER) + DEC, HEIGHT - (DEC + 10))
	}
	for i := 0; i < 19; i++ {
		pixmap.GetDrawable().DrawLine(gc, DEC, (i * INTER) + DEC, WIDTH - (DEC + 10), (i * INTER) + DEC)
	}
	pixmap.GetDrawable().DrawArc(gc, true, (3 * INTER) + DEC - 5, (3 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (9 * INTER) + DEC - 5, (3 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (15 * INTER) + DEC - 5, (3 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (3 * INTER) + DEC - 5, (9 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (9 * INTER) + DEC - 5, (9 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (15 * INTER) + DEC - 5, (9 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (3 * INTER) + DEC - 5, (15 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (9 * INTER) + DEC - 5, (15 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
	pixmap.GetDrawable().DrawArc(gc, true, (15 * INTER) + DEC - 5, (15 * INTER) + DEC - 5, 10, 10, 0, 64 * 360)
}

func main() {
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
	var p1, p2 point
	p1.x = -1
	p1.y = -1

	vbox := gtk.VBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea := gtk.DrawingArea()

	// pixbuf := gtk.ImageFromFile("goban.png")

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.GtkAllocation
		drawingarea.GetAllocation(&allocation)
		pixmap = gdk.Pixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.GC(pixmap.GetDrawable())
		init_grid(gc, pixmap)
	})

	drawingarea.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		if gdkwin == nil {
			gdkwin = drawingarea.GetWindow()
		}
		arg := ctx.Args(0)
		mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
		var mt gdk.GdkModifierType
		if mev.IsHint != 0 {
			gdkwin.GetPointer(&p2.x, &p2.y, &mt)
		} else {
			p2.x, p2.y = int(mev.X), int(mev.Y)
		}
		if p1.x != -1 && p2.x != -1 && (gdk.GdkEventMask(mt)&gdk.GDK_BUTTON_PRESS_MASK) != 0 {
			pixmap.GetDrawable().DrawLine(gc, p1.x, p1.y, p2.x, p2.y)
			drawingarea.GetWindow().Invalidate(nil, false)
		}
		p1 = p2
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