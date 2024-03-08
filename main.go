package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"syscall/js"

	"github.com/nfnt/resize"
)

func exports() {
	js.Global().Set("optimizeImage", js.FuncOf(optimizeImage))
}

func optimizeImage(this js.Value, p []js.Value) interface{} {
	inputImage := make([]byte, p[0].Get("length").Int())
	width := uint(p[1].Int())
	height := uint(p[2].Int())
	filter := p[3].String()
	js.CopyBytesToGo(inputImage, p[0])

	img, _, err := image.Decode(bytes.NewReader(inputImage))
	if err != nil {
		return err
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
	var outputBuffer bytes.Buffer
	if filter == "true" {
		bouds := resizedImg.Bounds()
		gray := image.NewGray(bouds)

		for y := bouds.Min.Y; y < bouds.Max.Y; y++ {
			for x := bouds.Min.X; x < bouds.Max.X; x++ {
				px := img.At(x, y)
				gr := color.GrayModel.Convert(px)
				gray.Set(x, y, gr)
			}
		}

		err = jpeg.Encode(&outputBuffer, gray, nil)
		if err != nil {
			return err
		}
	} else {
		err = jpeg.Encode(&outputBuffer, resizedImg, nil)
		if err != nil {
			return err
		}
	}

	jsOutput := js.Global().Get("Uint8Array").New(len(outputBuffer.Bytes()))
	js.CopyBytesToJS(jsOutput, outputBuffer.Bytes())
	fmt.Println("optimizeImage called")
	return jsOutput
}

func main() {
	fmt.Println("WASM Go Initialized")
	c := make(chan struct{})
	exports()
	<-c
}
