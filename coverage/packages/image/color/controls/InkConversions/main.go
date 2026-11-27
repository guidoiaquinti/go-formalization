package main

import (
	"fmt"
	"image/color"
)

func main() {
	c, m, y, k := color.RGBToCMYK(0, 0, 0)
	if (color.CMYK{C: c, M: m, Y: y, K: k}) != (color.CMYK{K: 255}) {
		panic("RGB black did not become full black ink")
	}
	c, m, y, k = color.RGBToCMYK(7, 5, 1)
	if (color.CMYK{C: c, M: m, Y: y, K: k}) != (color.CMYK{M: 72, Y: 218, K: 248}) {
		panic("RGB conversion did not truncate the scaled deficits")
	}

	r, g, b := color.CMYKToRGB(9, 0, 255, 13)
	if [3]uint8{r, g, b} != [3]uint8{234, 242, 0} {
		panic("CMYK conversion did not use widened scaling")
	}
	r, g, b = color.CMYKToRGB(0, 1, 2, 255)
	if r != 0 || g != 0 || b != 0 {
		panic("full black ink did not produce black RGB")
	}

	roundedDeficit := (255*(7-5) + 7/2) / 7
	narrowProduct := (255 - 9) * (255 - 13) / 255
	if roundedDeficit != 73 || narrowProduct != 233 {
		panic("unexpected mutant results")
	}
	fmt.Println("Ink conversions: Go passed four fixtures; rounding and direct byte scaling rejected")
}
