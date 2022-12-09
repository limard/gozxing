package detector

import (
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

// MultiDetector Encapsulates logic that can detect one or more QR Codes in an image,
// even if the QR Code is rotated or skewed, or partially obscured.
type MultiDetector struct {
	*Detector
}

func NewMultiDetector(image *gozxing.BitMatrix) *MultiDetector {
	return &MultiDetector{
		NewDetector(image),
	}
}

func (this *MultiDetector) DetectMulti(hints map[gozxing.DecodeHintType]interface{}) ([]*common.DetectorResult, error) {
	image := this.GetImage()
	resultPointCallback, _ := hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK].(gozxing.ResultPointCallback)

	finder := NewMultiFinderPatternFinder(image, resultPointCallback)
	infos, e := finder.FindMulti(hints)
	if e != nil || len(infos) == 0 {
		return nil, gozxing.WrapNotFoundException(e)
	}

	result := make([]*common.DetectorResult, 0)
	for _, info := range infos {
		r, e := this.ProcessFinderPatternInfo(info)
		if e != nil {
			// ignore
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

func (this *MultiDetector) DetectArea(hints map[gozxing.DecodeHintType]interface{}) (areas []gozxing.Area, e error) {
	min := func(a, b float64) int {
		if a > b {
			return int(b)
		}
		return int(a)
	}
	max := func(a, b float64) int {
		if a < b {
			return int(b)
		}
		return int(a)
	}

	if hints != nil {
		if cb, ok := hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK]; ok {
			this.resultPointCallback, _ = cb.(gozxing.ResultPointCallback)
		}
	}

	finder := NewMultiFinderPatternFinder(this.image, this.resultPointCallback)
	infos, e1 := finder.FindMulti(hints)
	if e1 != nil {
		e = e1
		return
	}

	for _, info := range infos {
		topLeft := info.GetTopLeft()
		topRight := info.GetTopRight()
		bottomLeft := info.GetBottomLeft()

		//log.Printf("topLeft:%v\n", topLeft)
		//log.Printf("topRight:%v\n", topRight)
		//log.Printf("bottomLeft:%v\n", bottomLeft)

		//每个方格占的像素大小
		moduleSize := this.calculateModuleSize(topLeft, topRight, bottomLeft)
		if moduleSize < 1.0 {
			e = gozxing.NewNotFoundException("moduleSize = %v", moduleSize)
			return
		}
		//log.Printf("moduleSize:%v\n", moduleSize)

		// 0  : x topLeft=bottomLeft(min) topLeft<topRight; y topLeft=topRight(min) topLeft<bottomLeft
		// 90 : x topLeft=topRight(max) topLeft>bottomLeft; y:topLeft=bottomLeft(min) topLeft<topRight
		// 180: x topLeft=bottomLeft(max) topLeft>topRight; y topLeft=topRight(max)  topLeft>bottomLeft
		// 270: x topLeft=topRight(min) topLeft<bottomLeft; y topLeft=bottomLeft(max) topLeft>topRight
		borderSize := int(3.5*moduleSize + 0.5) //向上取证
		area := gozxing.Area{}
		switch {
		case topLeft.GetX() < topRight.GetX() && topLeft.GetY() < bottomLeft.GetY():
			area.Rotate = 0
			area.ImageRectangle.Min.X = min(topLeft.GetX(), bottomLeft.GetX()) - borderSize
			area.ImageRectangle.Max.X = int(topRight.GetX()) + borderSize
			area.ImageRectangle.Min.Y = min(topLeft.GetY(), topRight.GetY()) - borderSize
			area.ImageRectangle.Max.Y = int(bottomLeft.GetY()) + borderSize
			area.PaperRectangle = area.ImageRectangle
		case topLeft.GetX() > bottomLeft.GetX() && topLeft.GetY() < topRight.GetY():
			area.Rotate = 90
			area.ImageRectangle.Min.X = int(bottomLeft.GetX()) - borderSize
			area.ImageRectangle.Max.X = max(topLeft.GetX(), topRight.GetX()) + borderSize
			area.ImageRectangle.Min.Y = min(topLeft.GetY(), bottomLeft.GetY()) - borderSize
			area.ImageRectangle.Max.Y = int(topRight.GetY()) + borderSize
			area.PaperRectangle.Min.X, area.PaperRectangle.Min.Y =
				rotatePaper(270, this.image.GetWidth(), this.image.GetHeight(), int(topLeft.GetX())+borderSize, int(topLeft.GetY())-borderSize)
			area.PaperRectangle.Max.X, area.PaperRectangle.Max.Y =
				rotatePaper(270, this.image.GetWidth(), this.image.GetHeight(), int(bottomLeft.GetX())-borderSize, int(topRight.GetY())+borderSize)
		case topLeft.GetX() > topRight.GetX() && topLeft.GetY() > bottomLeft.GetY():
			area.Rotate = 180
			area.ImageRectangle.Min.X = int(topRight.GetX()) - borderSize
			area.ImageRectangle.Max.X = max(topLeft.GetX(), bottomLeft.GetX()) + borderSize
			area.ImageRectangle.Min.Y = int(bottomLeft.GetY()) - borderSize
			area.ImageRectangle.Max.Y = max(topLeft.GetY(), topRight.GetY()) + borderSize
			area.PaperRectangle.Min.X, area.PaperRectangle.Min.Y =
				rotatePaper(180, this.image.GetWidth(), this.image.GetHeight(), int(topLeft.GetX())+borderSize, int(topLeft.GetY())+borderSize)
			area.PaperRectangle.Max.X, area.PaperRectangle.Max.Y =
				rotatePaper(180, this.image.GetWidth(), this.image.GetHeight(), int(topRight.GetX())-borderSize, int(bottomLeft.GetY())-borderSize)
		case topLeft.GetX() < bottomLeft.GetX() && topLeft.GetY() > topRight.GetY():
			area.Rotate = 270
			area.ImageRectangle.Min.X = min(topLeft.GetX(), topRight.GetX()) - borderSize
			area.ImageRectangle.Max.X = int(bottomLeft.GetX()) + borderSize
			area.ImageRectangle.Min.Y = int(topRight.GetY()) - borderSize
			area.ImageRectangle.Max.Y = max(topLeft.GetY(), bottomLeft.GetY()) + borderSize
			area.PaperRectangle.Min.X, area.PaperRectangle.Min.Y =
				rotatePaper(90, this.image.GetWidth(), this.image.GetHeight(), int(topLeft.GetX())-borderSize, int(topLeft.GetY())+borderSize)
			area.PaperRectangle.Max.X, area.PaperRectangle.Max.Y =
				rotatePaper(90, this.image.GetWidth(), this.image.GetHeight(), int(bottomLeft.GetX())+borderSize, int(topRight.GetY())-borderSize)
		}
		areas = append(areas, area)
	}

	return
}

// rotatePaper 计算纸张调转方向后，(x, y)在新纸张中的位置
// rotateAngle 纸张顺时针旋转的角度
func rotatePaper(rotateAngle, paperWidth, paperHeight, x, y int) (nX, nY int) {
	switch rotateAngle {
	case 0:
		nX = x
		nY = y
	case 90:
		nX = paperHeight - y
		nY = x
	case 180:
		nX = paperWidth - x
		nY = paperHeight - y
	case 270:
		nX = y
		nY = paperWidth - x
	}
	return
}
