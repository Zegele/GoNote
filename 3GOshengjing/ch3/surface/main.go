// Surface computes an SVG rendering of a 3-D surface function.

package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320            // canvas size in pixels(画布尺寸像素)
	cells         = 100                 // number of grid(网格，方格) cells （网格单元数）
	xyrange       = 30.0                //axis ranges (-xyrange..+xyrange) (轴的范围，负到正的范围)
	xyscale       = width / 2 / xyrange // 10； pixels per x or y unit (每x或y单位的像素)
	zscale        = height * 0.4        // 128； pixels per z unit（每z单位的像素）
	angle         = math.Pi / 6         //angle of x, y axes (30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Printf("<polygon points = '%g, %g, %g, %g, %g, %g, %g, %g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)

		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64) {
	//Find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	//Compute surface height z.
	z := f(x, y)

	//project (x, y, z) isometrically onto 2-D SVG canvas (sx, sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from(0,0)
	return math.Sin(r) / r
}

// go build  生成一个.exe可执行文件。
// .surface.exe > s.svg 生成一个svg文件。
