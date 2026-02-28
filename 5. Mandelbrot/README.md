# Mandelbrot Set Generator

A concurrent Mandelbrot set renderer written in Go that outputs a 1980×1980 PNG image with a custom color palette.

## What is the Mandelbrot Set?

The Mandelbrot set is a famous fractal defined in the complex number plane. For each complex number **c**, the iteration **z = z² + c** is repeated starting from **z = 0**. If the value of **z** stays bounded (doesn't escape to infinity), **c** belongs to the Mandelbrot set. The boundary of this set produces infinitely complex, self-similar patterns that have become one of the most recognizable images in mathematics.

## What This Project Does

- Maps each pixel of a 1980×1980 image to a point on the complex plane (range −2 to +2 on both axes).
- Iterates the Mandelbrot formula up to 255 times per pixel to determine if the point escapes.
- Colors escaping points using smooth interpolation across a custom 4-color palette.
- Uses **8 goroutines** via a worker pool to compute pixels concurrently.
- Saves the result as `output.png`.

## Usage

```sh
go run main.go
```

The rendered image will be saved as `output.png` in the project directory.
