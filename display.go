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
	WINNER string
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
		40,
		40)
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
		clean_side(gc, pixmap, 0, y*INTER+DEC, INTER, INTER) // LEFT
	} else if x == 18 {
		clean_side(gc, pixmap, HEIGHT-INTER+1, y*INTER+DEC, -1, INTER) // RIGHT
	}
	if y == 0 {
		clean_side(gc, pixmap, x*INTER+DEC, 0, x*INTER+INTER+DEC, INTER) // TOP
	} else if y == 18 {
		clean_side(gc, pixmap, x*INTER+DEC, HEIGHT-INTER+1, INTER, -1) // BOT
	}
}

func display_init_grid(gc *gdk.GdkGC, pixmap *gdk.GdkPixmap) {
	gc.SetRgbFgColor(gdk.Color("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	for x := 0; x < 19; x++ {
		for y := 0; y < 19; y++ {
			gc.SetRgbFgColor(gdk.Color("grey"))
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

func board_display() {
	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetPosition(gtk.GTK_WIN_POS_CENTER)
	window.SetTitle("Gomoku")
	window.SetResizable(false)
	window.Connect("destroy", func() {
		println("got destroy!")
		gtk.MainQuit()
	})

	var game Gomoku
	var endGame, doubleThree, stop bool
	var menuitem *gtk.GtkMenuItem
	var gdkwin *gdk.GdkWindow
	var pixmap *gdk.GdkPixmap
	var gc *gdk.GdkGC
	var player int
	player = 1
	game = Gomoku{make([]int, 361), true, endGame, doubleThree, 1, [2]int{10, 10}}

	vbox := gtk.VBox(false, 1)

	menubar := gtk.MenuBar()
	vbox.PackStart(menubar, false, false, 0)
	statusbar := gtk.Statusbar()
	context_id := statusbar.GetContextId("go-gtk")
	statusbar.Push(context_id, "(not so) Proudly propulsed by the inglorious Gomoku Project, with love, and Golang!")

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
		if stop == true {
			return
		}
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
		if ((x-INTER/2)/INTER) < 0 || ((x-INTER/2)/INTER) >= 19 ||
			((y-INTER/2)/INTER) < 0 || ((y-INTER/2)/INTER) >= 19 {
			return
		}
		vic, stones, err := game.Play(((x - INTER/2) / INTER), ((y - INTER/2) / INTER))
		if err != nil {
			return
		}
		statusbar.Push(context_id, fmt.Sprintf("Last move is Player %d : %d/%d", player, ((x-INTER/2)/INTER)+1, ((y-INTER/2)/INTER)+1))
		for _, stone := range stones {
			draw_square(gc, pixmap, stone[0], stone[1])
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
			WINNER = fmt.Sprintf("And the winner is \"Player %d\"", vic)
			context_id := statusbar.GetContextId("go-gtk")
			statusbar.Push(context_id, WINNER)
			stop = true
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

	cascademenu := gtk.MenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.Menu()
	cascademenu.SetSubmenu(submenu)
	menuitem = gtk.MenuItemWithMnemonic("_Player Vs Player")
	menuitem.Connect("activate", func() {
		gc.SetRgbFgColor(gdk.Color("grey"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		fmt.Println(endGame, doubleThree)
		game = Gomoku{make([]int, 361), true, endGame, doubleThree, 1, [2]int{10, 10}}
		player = 1
		display_init_grid(gc, pixmap)
		drawingarea.Hide()
		drawingarea.Show()
		stop = false
	})
	submenu.Append(menuitem)
	menuitem = gtk.MenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.MenuItemWithMnemonic("_Options")
	menubar.Append(cascademenu)
	submenu = gtk.Menu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.CheckMenuItemWithMnemonic("_Double trois")
	checkmenuitem.Connect("activate", func() {
		if doubleThree == false {
			doubleThree = true
		} else {
			doubleThree = false
		}
		game = Gomoku{make([]int, 361), true, endGame, doubleThree, 1, [2]int{10, 10}}
		player = 1
		display_init_grid(gc, pixmap)
		drawingarea.Hide()
		drawingarea.Show()
		stop = false
	})
	submenu.Append(checkmenuitem)

	checkmenuitem = gtk.CheckMenuItemWithMnemonic("P_rise de fin de partie")
	checkmenuitem.Connect("activate", func() {
		if endGame == false {
			endGame = true
		} else {
			endGame = false
		}
		game = Gomoku{make([]int, 361), true, endGame, doubleThree, 1, [2]int{10, 10}}
		player = 1
		display_init_grid(gc, pixmap)
		drawingarea.Hide()
		drawingarea.Show()
		stop = false
	})
	submenu.Append(checkmenuitem)

	vbox.PackStart(statusbar, false, false, 0)

	window.Add(vbox)
	window.SetSizeRequest(WIDTH, HEIGHT+20)
	window.ShowAll()
	gtk.Main()
}
