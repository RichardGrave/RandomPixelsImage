package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const argsImagePos = 1
const argsWidthPos = 2
const argsHeightPos = 3
const defaultPixels = 500
const fileExtension = ".png"

func main() {
	start := time.Now()

	//Set random Seed for generating Random Numbers
	rand.Seed(time.Now().UTC().UnixNano())
	imageName := "random_image"
	width := defaultPixels
	height := defaultPixels

	if len(os.Args) > 1 {
		imageName = os.Args[argsImagePos]
		if len(os.Args) > 2 {
			if tmpWidth, err := strconv.Atoi(os.Args[argsWidthPos]); err == nil {
				width = tmpWidth
			} else {
				fmt.Println("Width is not a number")
			}
			if len(os.Args) > 3 {
				if tmpHeight, err := strconv.Atoi(os.Args[argsHeightPos]); err == nil {
					height = tmpHeight
				} else {
					fmt.Println("Height is not a number")
				}
			} else {
				fmt.Println("Same length of width is used for height")
				height = width
			}
		} else {
			fmt.Println("Default width and height used")
		}
	} else {
		fmt.Println("Default imagename, width and height used")
	}
	createImage(imageName, width, height)

	elapsed := time.Since(start)

	fmt.Println("Start: ", start, " - Elapsed time: ", elapsed)
}

func createImage(imageName string, width int, height int) {
	imageFile := imageName + fileExtension

	if _, err := os.Stat(imageFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist. So its being created\n", imageName)
			newFile, err := os.Create(imageFile)
			if err != nil {
				log.Fatal(err)
			}
			defer newFile.Close()
			writePixelsToFile(newFile, width, height)
		}
	} else {
		if file, removeErr := os.OpenFile(imageFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); removeErr == nil {
			fmt.Printf("File %s already exist. So its being reused\n", imageName)
			file.Truncate(0)
			file.Seek(0, 0)
			defer file.Close()
			writePixelsToFile(file, width, height)

		} else {
			fmt.Println("Could not open file: ", removeErr)
		}
	}
}

//RandomColor :Save a colorPoint for the image
type RandomColor struct {
	//RandomColor
	posx int
	posy int

	colorR uint8
	colorG uint8
	colorB uint8
	colorA uint8
}

func writePixelsToFile(imageFile *os.File, width, height int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	maxItems := width * height

	randomColorChan := make(chan RandomColor, maxItems)

	go createRandomPixel(randomColorChan, 0, 0, width, height)

	iterColorLoop := 0
	for iterColorLoop < maxItems {

		randomColorType := <-randomColorChan
		img.Set(randomColorType.posx, randomColorType.posy,
			color.RGBA{randomColorType.colorR, randomColorType.colorG, randomColorType.colorB, randomColorType.colorA})

		iterColorLoop++
	}

	png.Encode(imageFile, img)

}

//TODO:RG this is not OK.
// We don't want recursion because it uses to much memory if the width and height are to big.
func createRandomPixel(randomColorChan chan RandomColor, posx, posy, width, height int) {
	var randomCol RandomColor

	randomCol.posx = posx
	randomCol.posy = posy

	randomCol.colorR = uint8(rand.Intn(255))
	randomCol.colorG = uint8(rand.Intn(255))
	randomCol.colorB = uint8(rand.Intn(255))
	randomCol.colorA = uint8(rand.Intn(255))

	fmt.Println(posy, " - ", posx)

	randomColorChan <- randomCol

	if posx == (width-1) && posy == (height-1) {
		return
	}

	if posx == (width - 1) {
		posx = 0
		posy++
	} else {
		posx++
	}

	go createRandomPixel(randomColorChan, posx, posy, width, height)
}
