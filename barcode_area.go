package gozxing

import (
	"fmt"
	"image"
)

// Area describes the position and rotation of barcode
type Area struct {
	ImageRectangle image.Rectangle // postion in image (no rotation)
	PaperRectangle image.Rectangle // Rotate the paper according to the direction of Barcode. postion on image (rotatsion)
	Rotate         int             // rotation angle (clockwise)
}

func (t Area) String() string {
	return fmt.Sprintf("Image:%s Paper:%s Rotate:%d", t.ImageRectangle.String(), t.PaperRectangle.String(), t.Rotate)
}
