package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
)

/*

|
|	上屏
|(0,0)
|----------------------->x
|
|	下屏
|
\/
y

*/

const (
	diffLimit            = 3
	diffMaxScaleBoundary = 10
)

type Scale struct {
	scale  float64 // 设置鼠标增加的倍数, 仅在上屏放大
	RWLock sync.RWMutex
}

var scaleObj Scale

func main() {
	flag.Float64Var(&scaleObj.scale, "scale", 1.0, "鼠标速度增加的放大倍数")
	flag.Parse()
	var lastMousePosX, lastMousePosY int
	var currentMousePosX, currentMousePosY int
	lastMousePosX, lastMousePosY = robotgo.Location()
	go func() {
		for {
			currentMousePosX, currentMousePosY = robotgo.Location()
			if currentMousePosY >= 0 {
				continue
			}
			xDiff := currentMousePosX - lastMousePosX
			yDiff := currentMousePosY - lastMousePosY

			if xDiff >= diffLimit || yDiff >= diffLimit {
				scaleObj.RWLock.RLock()
				curScale := scaleObj.scale * min((max(float64(xDiff), float64(yDiff))/diffMaxScaleBoundary), 1.0)
				scaleObj.RWLock.RUnlock()
				newX := float64(currentMousePosX) + float64(xDiff)*curScale
				newY := float64(currentMousePosY) + float64(yDiff)*curScale

				robotgo.Move(int(newX), int(newY))
				lastMousePosX, lastMousePosY = int(newX), int(newY)
			} else {
				lastMousePosX, lastMousePosY = currentMousePosX, currentMousePosY
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()
	systray.Run(onReady, onExit)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func onReady() {
	var strTitle string
	strTitle = fmt.Sprintf("%.1f", scaleObj.scale+1)
	systray.SetTitle("x" + strTitle)
	quit := systray.AddMenuItem("Quit", "Quit the app")
	addScale := systray.AddMenuItem("Add", "Set the scale of mouse speed")
	subScale := systray.AddMenuItem("Sub", "Set the scale of mouse speed")
	terminalInput := systray.AddMenuItem("Terminal", "Input the scale of mouse speed")
	go func() {
		for {
			select {
			case <-quit.ClickedCh:
				systray.Quit()
			case <-addScale.ClickedCh:
				scaleObj.RWLock.Lock()
				scaleObj.scale += 0.1
				strTitle = fmt.Sprintf("%.1f", scaleObj.scale+1)
				scaleObj.RWLock.Unlock()
			case <-subScale.ClickedCh:
				scaleObj.RWLock.Lock()
				scaleObj.scale -= 0.1
				if scaleObj.scale <= 0.01 {
					scaleObj.scale = 0
				}
				strTitle = fmt.Sprintf("%.1f", scaleObj.scale+1)
				scaleObj.RWLock.Unlock()
			case <-terminalInput.ClickedCh:
				fmt.Println("Please input the scale of mouse speed:")
				var scale float64
				fmt.Scanln(&scale)
				if scale <= 1 {
					scale = 1
				}
				scaleObj.RWLock.Lock()
				scaleObj.scale = scale - 1
				strTitle = fmt.Sprintf("%.1f", scaleObj.scale+1)
				scaleObj.RWLock.Unlock()
			}
			systray.SetTitle("x" + strTitle)
		}
	}()
}

func onExit() {
	os.Exit(0)
}
