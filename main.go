package main

import (
	//"fmt"
	"flag"
	"github.com/go-vgo/robotgo"
	"time"
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

func main() {
	var scale float64 // 设置鼠标增加的倍数, 仅在上屏放大
	flag.Float64Var(&scale, "scale", 0.5, "鼠标速度增加的放大倍数")
	flag.Parse()
	var lastMousePosX, lastMousePosY int
	var currentMousePosX, currentMousePosY int
	lastMousePosX, lastMousePosY = robotgo.Location()
	for {
		currentMousePosX, currentMousePosY = robotgo.Location()
		if currentMousePosY >= 0 {
			continue
		}
		xDiff := currentMousePosX - lastMousePosX
		yDiff := currentMousePosY - lastMousePosY

		if xDiff >= diffLimit || yDiff >= diffLimit {
			//fmt.Printf("Mouse moved, x: %v, y: %v\n", xDiff, yDiff)
			// 计算新的鼠标位置
			curScale := scale * min((max(float64(xDiff), float64(yDiff))/diffMaxScaleBoundary), 1.0)
			newX := float64(currentMousePosX) + float64(xDiff)*curScale
			newY := float64(currentMousePosY) + float64(yDiff)*curScale

			// 设置鼠标位置
			robotgo.Move(int(newX), int(newY))
			lastMousePosX, lastMousePosY = int(newX), int(newY)
		} else {
			lastMousePosX, lastMousePosY = currentMousePosX, currentMousePosY
		}

		// 每10毫秒读取一次鼠标的位置
		time.Sleep(10 * time.Millisecond)
	}
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
