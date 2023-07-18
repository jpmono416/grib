package griblib

import (
	"fmt"
	"github.com/jpmono416/grib/griblib/jpeg2000"
	"image"
	"io"
)

// Data40 is a Grid point data - JPEG-2000 packing
//
// | Octet Number | Content
// ---------------------------------------------------------------------------------------
// | 12-15        | Reference value (R) (IEEE 32-bit floating-point value)
// | 16-17        | Binary scale factor (E)
// | 18-19        | Decimal scale factor (D)
// | 20           | Number of bits used for each packed value
// | 21           | Type of original field values
// | 22           | Group splitting method used
// | 23           | Missing value management used
// | 24-27        | Primary missing value substitute
// | 28-31        | Secondary missing value substitute
type Data40 struct {
	Data0
	OriginalFieldType          uint8  `json:"originalFieldType"`          // 21
	GroupSplittingMethod       uint8  `json:"groupSplittingMethod"`       // 22
	MissingValueManagement     uint8  `json:"missingValueManagement"`     // 23
	PrimaryMissingSubstitute   uint32 `json:"primaryMissingSubstitute"`   // 24-27
	SecondaryMissingSubstitute uint32 `json:"secondaryMissingSubstitute"` // 28-31
}

func ParseData40(dataReader io.Reader) ([]float64, error) {
	imageData, err := decodeImageData(dataReader)
	if err != nil {
		return nil, err
	}
	// Extract pixel data from imageData and return as []float64
	pixelData := extractPixelData(imageData)
	return pixelData, nil
}

func decodeImageData(dataReader io.Reader) (image.Image, error) {
	// Read the JPEG-2000 compressed data
	// TODO NOT WORKING, find a suitable way to decode JP2 format
	var byteData []byte
	n, err := io.ReadFull(dataReader, byteData)

	if err != nil || n == 0 {
		// Default 1px square empty image
		return image.NewRGBA(image.Rect(0, 0, 1, 1)),
			fmt.Errorf("Error reading compressed data: %s", err.Error())
	}

	imageData, err := jpeg2000.Parse(byteData)
	if err != nil {
		// Default 1px square empty image
		return image.NewRGBA(image.Rect(0, 0, 1, 1)),
			fmt.Errorf("Error parsing data: %s", err.Error())
	}
	return imageData, nil
}

func extractPixelData(imageData image.Image) []float64 {
	bounds := imageData.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	pixelData := make([]float64, width*height)
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := imageData.At(x, y).RGBA()
			pixelData[idx] = float64(r)
			idx++
		}
	}

	return pixelData
}
