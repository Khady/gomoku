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

var game Gomoku
var stop, iamode bool
var menuitem *gtk.MenuItem
var gdkwin *gdk.Window
var pixmap *gdk.Pixmap
var gc *gdk.GC
var player, countTake int
var statusbar *gtk.Statusbar
var drawingarea *gtk.DrawingArea
var hint *gtk.Label

func event_play(x, y int) bool {
	vic, stones, err := game.Play(x, y)
	if err != nil {
		return false
	}
	context_id := statusbar.GetContextId("go-gtk")
	statusbar.Push(context_id, fmt.Sprintf("[Player 1/2 : %d/%d stone before death] Last move is Player %d : %d/%d",
		game.countTake[1], game.countTake[0], player, x+1, y+1))
	for _, stone := range stones {
		countTake++
		if countTake > 19 {
			break
		}
		draw_square(gc, pixmap, stone[0], stone[1])
		if player == 1 {
			gc.SetRgbFgColor(gdk.NewColor("white"))
		} else {
			gc.SetRgbFgColor(gdk.NewColor("black"))
		}
		tmpx := 800
		tmpy := countTake * INTER
		tmpx = ((tmpx-INTER/2)/INTER)*INTER + INTER
		tmpy = ((tmpy-INTER/2)/INTER)*INTER + INTER
		pixmap.GetDrawable().DrawArc(gc, true, tmpx-(STONE/2)+10,
			tmpy-(STONE/2), STONE, STONE, 0, 64*360)
	}
	if player == 1 {
		gc.SetRgbFgColor(gdk.NewColor("black"))
		player = 2
	} else {
		gc.SetRgbFgColor(gdk.NewColor("white"))
		player = 1
	}
	x = x*INTER + INTER
	y = y*INTER + INTER
	pixmap.GetDrawable().DrawArc(gc, true, x-(STONE/2), y-(STONE/2), STONE, STONE, 0, 64*360)
	if vic != 0 {
		WINNER = fmt.Sprintf("And the winner is \"Player %d\"", vic)
		context_id := statusbar.GetContextId("go-gtk")
		statusbar.Push(context_id, WINNER)
		stop = true
	}
	drawingarea.GetWindow().Invalidate(nil, false)
	if vic != 0 {
		messagedialog := gtk.NewMessageDialog(
			statusbar.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK,
			WINNER)
		messagedialog.Response(func() {
			messagedialog.Destroy()
		})
		messagedialog.Run()
	}
	return true
}

func clean_side(gc *gdk.GC, pixmap *gdk.Pixmap, x1, y1, x2, y2 int) {
	gc.SetRgbFgColor(gdk.NewColor("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true,
		x1,
		y1,
		x2,
		y2)
}

func draw_square(gc *gdk.GC, pixmap *gdk.Pixmap, x, y int) {
	gc.SetRgbFgColor(gdk.NewColor("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true,
		x*INTER+DEC,
		y*INTER+DEC,
		40,
		40)
	gc.SetRgbFgColor(gdk.NewColor("black"))
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

func display_init_grid(gc *gdk.GC, pixmap *gdk.Pixmap) {
	gc.SetRgbFgColor(gdk.NewColor("grey"))
	pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	for x := 0; x < 19; x++ {
		for y := 0; y < 19; y++ {
			gc.SetRgbFgColor(gdk.NewColor("grey"))
			draw_square(gc, pixmap, x, y)
		}
	}
	gc.SetRgbFgColor(gdk.NewColor("black"))
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

func status_bar(vbox *gtk.VBox) {
	statusbar = gtk.NewStatusbar()
	context_id := statusbar.GetContextId("go-gtk")
	statusbar.Push(context_id, "(not so) Proudly propulsed by the inglorious Gomoku Project, with love, and Golang!")
	vbox.PackStart(statusbar, false, false, 0)
}

func menu_bar(vbox *gtk.VBox) {
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	buttons := gtk.NewAlignment(0, 0, 0, 0)
	checkbox := gtk.NewAlignment(1, 0, 0, 0)
	newPlayerGameButton := gtk.NewButtonWithLabel("Player vs Player")
	newIaGameButton := gtk.NewButtonWithLabel("Player vs AI")
	hint = gtk.NewLabel("Hint: Not yet")
	threeCheckBox := gtk.NewCheckButtonWithLabel("Three and three")
	endCheckBox := gtk.NewCheckButtonWithLabel("Unbreakable end")
	hbox := gtk.NewHBox(false, 1)
	hbox0 := gtk.NewHBox(false, 1)
	hbox1 := gtk.NewHBox(false, 1)
	hbox0.Add(newPlayerGameButton)
	hbox0.Add(newIaGameButton)
	hbox1.Add(threeCheckBox)
	hbox1.Add(endCheckBox)
	buttons.Add(hbox0)
	checkbox.Add(hbox1)
	hbox.Add(buttons)
	hbox.Add(hint)
	hbox.Add(checkbox)
	vbox.PackStart(hbox, false, true, 0)

	cascademenu := gtk.NewMenuItemWithMnemonic("_Game")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	playermenuitem := gtk.NewMenuItemWithMnemonic("_Player Vs Player")
	playermenuitem.Connect("activate", func() {
		gc.SetRgbFgColor(gdk.NewColor("grey"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		game = Gomoku{make([]int, 361), true, game.endgameTake, game.doubleThree, 1, [2]int{10, 10}, 0}
		player = 1
		countTake = 0
		iamode = false
		display_init_grid(gc, pixmap)
		drawingarea.Hide()
		drawingarea.Show()
		stop = false
		context_id := statusbar.GetContextId("go-gtk")
		statusbar.Push(context_id, "(not so) Proudly propulsed by the inglorious Gomoku Project, with love, and Golang!")
	})
	submenu.Append(playermenuitem)
	newPlayerGameButton.Clicked(func() {
		playermenuitem.Activate()
	})
	iamenuitem := gtk.NewMenuItemWithMnemonic("_Player Vs AI")
	iamenuitem.Connect("activate", func() {
		gc.SetRgbFgColor(gdk.NewColor("grey"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		game = Gomoku{make([]int, 361), true, game.endgameTake, game.doubleThree, 1, [2]int{10, 10}, 0}
		player = 1
		countTake = 0
		iamode = true
		display_init_grid(gc, pixmap)
		drawingarea.Hide()
		drawingarea.Show()
		stop = false
		context_id := statusbar.GetContextId("go-gtk")
		statusbar.Push(context_id, "(not so) Proudly propulsed by the inglorious Gomoku Project, with love, and Golang!")
	})
	submenu.Append(iamenuitem)
	newIaGameButton.Clicked(func() {
		iamenuitem.Activate()
	})
	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Rules")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	threemenuitem := gtk.NewCheckMenuItemWithMnemonic("_Three and three")
	threemenuitem.Connect("activate", func() {
		if game.doubleThree == false {
			game.doubleThree = true
		} else {
			game.doubleThree = false
		}
	})
	submenu.Append(threemenuitem)
	threeCheckBox.Connect("toggled", func() {
		threemenuitem.Activate()
	})

	endmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Unbreakable end")
	endmenuitem.Connect("activate", func() {
		if game.endgameTake == false {
			game.endgameTake = true
		} else {
			game.endgameTake = false
		}
	})
	submenu.Append(endmenuitem)
	endCheckBox.Connect("toggled", func() {
		endmenuitem.Activate()
	})

}

func configure_board(vbox *gtk.VBox) {
	drawingarea = gtk.NewDrawingArea()
	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.Allocation
		drawingarea.GetAllocation(&allocation)
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		display_init_grid(gc, pixmap)
	})

	drawingarea.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		// Check if the game is running and if player click in the goban
		if stop == true {
			return
		}
		if gdkwin == nil {
			gdkwin = drawingarea.GetWindow()
		}
		arg := ctx.Args(0)
		mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
		var mt gdk.ModifierType
		var x, y int
		if mev.IsHint != 0 {
			gdkwin.GetPointer(&x, &y, &mt)
		} else {
			x, y = int(mev.X), int(mev.Y)
		}
		x = ((x-INTER/2)/INTER)
		y = ((y-INTER/2)/INTER)
		if x < 0 || x >= 19 || y < 0 || y >= 19 {
			return
		}
		// end check
		if event_play(x, y) && iamode && stop != true {
			fmt.Println("ai turn")
			x, y = IATurn(&game)
			event_play(x, y)
		}
		x, y = IATurn(&game)
		hint.SetLabel(fmt.Sprintf("Hint: %d/%d", x, y))
	})

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	drawingarea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK))

	vbox.Add(drawingarea)
}

func board_display() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Gomoku")
	window.SetResizable(false)
	window.Connect("destroy", gtk.MainQuit)

	player = 1
	game = Gomoku{make([]int, 361), true, false, false, 1, [2]int{10, 10}, 0}

	vbox := gtk.NewVBox(false, 1)

	menu_bar(vbox)
	configure_board(vbox)
	status_bar(vbox)

	window.Add(vbox)
	window.SetSizeRequest(WIDTH+40, HEIGHT+50)
	window.ShowAll()
	gtk.Main()
}
