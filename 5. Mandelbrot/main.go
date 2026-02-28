package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"sync"
)

const WIDTH = 1980
const HEIGHT = 1980

// func mandelbrot(c complex128) color.Color {
//     const maxIter = 100
//     z := complex(0, 0)
//     for i := range maxIter {
//         z = z*z + c
//         if cmplx.Abs(z) > 2 {
//             shade := uint8(255 * i / maxIter)
//             return color.RGBA{shade, shade, shade, 255}
//         }
//     }
//     return color.Black
// }

func mandelbrot(c complex128) color.Color {
	const maxIter = 255
	z := complex(0, 0)

	colors := []color.RGBA{
		{43, 46, 74, 255},
		{232, 69, 69, 255},
		{144, 55, 73, 255},
		{83, 53, 74, 255},
	}

	for i := range maxIter {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			t := math.Sqrt(float64(i) / float64(maxIter))

			scaledT := t * float64(len(colors)-1)
			idx := int(scaledT)
			nextIdx := (idx + 1) % len(colors)
			frac := scaledT - float64(idx)

			return color.RGBA{
				uint8(float64(colors[idx].R)*(1-frac) + float64(colors[nextIdx].R)*frac),
				uint8(float64(colors[idx].G)*(1-frac) + float64(colors[nextIdx].G)*frac),
				uint8(float64(colors[idx].B)*(1-frac) + float64(colors[nextIdx].B)*frac),
				255,
			}
		}
	}

	return colors[3]
}

func convertPixelToComplex(x, y int) complex128 {
	normX := float64(x) / float64(WIDTH)
	normY := float64(y) / float64(HEIGHT)

	real := float64(normX*4 - 2)
	imag := float64(normY*4 - 2)

	return complex(real, imag)
}

func worker(jobs chan [2]int, results chan result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		x, y := job[0], job[1]
		c := convertPixelToComplex(x, y)
		results <- result{
			x:     x,
			y:     y,
			color: mandelbrot(c),
		}
	}

}

type result struct {
	x, y  int
	color color.Color
}

func main() {

	jobs := make(chan [2]int, 100)
	results := make(chan result, 100)

	var wg sync.WaitGroup
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for range 8 {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for y := range HEIGHT {
			for x := range WIDTH {
				jobs <- [2]int{x, y}
			}
		}
		close(jobs)
	}()

	for result := range results {
		img.Set(result.x, result.y, result.color)
	}

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
