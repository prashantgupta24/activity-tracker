package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("hello")

	x := 1
	y := 500
	robotgo.MoveMouse(x, y)
	color := robotgo.GetPixelColor(x, y)
	fmt.Println("color : ", color)

	// abool := robotgo.ShowAlert("test", "robotgo")
	// if abool == 0 {
	// 	fmt.Println("ok@@@ ", "ok")
	// }

	bitmap1 := robotgo.CaptureScreen(10, 20, 300, 400)
	bitmap2 := robotgo.CaptureScreen(20, 30, 300, 400)
	fmt.Println("...", bitmap1)
	fmt.Println("...", bitmap2)
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	// defer robotgo.FreeBitmap(bitmap1)
	// defer robotgo.FreeBitmap(bitmap2)
	fx1, fy1 := robotgo.FindBitmap(bitmap1)
	fx2, fy2 := robotgo.FindBitmap(bitmap2)

	fmt.Println("...", fx1, fx2, fy1, fy2)

	bitmap := robotgo.CaptureScreen(10, 20, 300, 400)
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	defer robotgo.FreeBitmap(bitmap)

	fmt.Println("...", bitmap)

	fx, fy := robotgo.FindBitmap(bitmap)
	fmt.Println("FindBitmap------ ", fx, fy)
	robotgo.SaveBitmap(bitmap, "test.png")

	title := robotgo.GetTitle()
	fmt.Println("title@@@ ", title)

	//robotgo.Start()
	// go func() {
	// 	for {
	// 		ok := robotgo.AddEvent("mleft")
	// 		if ok {
	// 			fmt.Println("you press...\n", "k")
	// 		}
	// 	}

	// }()

	// go func() {
	// 	for {
	// 		mleft := robotgo.AddMouse("left")
	// 		if mleft {
	// 			fmt.Println("you press...", "mouse left button")
	// 		}
	// 	}

	// }()
	//	robotgo.Start()

	// go func() {
	// 	for {
	// 		mleft := robotgo.AddEvent("mright")
	// 		if mleft {
	// 			fmt.Println("you press...", "mouse right button")
	// 		}
	// 	}

	// }()

	//	time.Sleep(time.Second * 200)
	// mright := robotgo.AddEvent("mright")
	// if mright {
	// 	log.Printf("mright clicked \n")
	// }
}
