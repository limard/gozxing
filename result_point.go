package gozxing

import (
	"github.com/makiuchi-d/gozxing/common/detector"
)

type ResultPoint interface {
	GetX() float64
	GetY() float64
}

type ResultPointBase struct {
	x float64
	y float64
}

func NewResultPoint(x, y float64) ResultPoint {
	return ResultPointBase{x, y}
}

func (rp ResultPointBase) GetX() float64 {
	return rp.x
}

func (rp ResultPointBase) GetY() float64 {
	return rp.y
}

// Orders an array of three ResultPoints in an order [A,B,C] such that AB is less than AC
// and BC is less than AC, and the angle between BC and BA is less than 180 degrees.
// @param patterns array of three {@code ResultPoint} to order
func ResultPoint_OrderBestPatterns(patterns []ResultPoint) {
	// Find distances between pattern centers
	zeroOneDistance := distance(patterns[0], patterns[1])
	oneTwoDistance := distance(patterns[1], patterns[2])
	zeroTwoDistance := distance(patterns[0], patterns[2])

	var pointA, pointB, pointC ResultPoint

	// Assume one closest to other two is B; A and C will just be guesses at first
	if oneTwoDistance >= zeroOneDistance && oneTwoDistance >= zeroTwoDistance {
		pointB = patterns[0]
		pointA = patterns[1]
		pointC = patterns[2]
	} else if zeroTwoDistance >= oneTwoDistance && zeroTwoDistance >= zeroOneDistance {
		pointB = patterns[1]
		pointA = patterns[0]
		pointC = patterns[2]
	} else {
		pointB = patterns[2]
		pointA = patterns[0]
		pointC = patterns[1]
	}

    // Use cross product to figure out whether A and C are correct or flipped.
    // This asks whether BC x BA has a positive z component, which is the arrangement
    // we want for A, B, C. If it's negative, then we've got it flipped around and
    // should swap A and C.
	if crossProductZ(pointA, pointB, pointC) < 0.0 {
		pointA, pointC = pointC, pointA
	}

	patterns[0] = pointA
	patterns[1] = pointB
	patterns[2] = pointC
}

func distance(pattern1, pattern2 ResultPoint) float64 {
	return detector.MathUtils_DistanceFloat(pattern1.GetX(), pattern1.GetY(), pattern2.GetX(), pattern2.GetY())
}

func crossProductZ(pointA, pointB, pointC ResultPoint) float64 {
	bX := pointB.GetX()
	bY := pointB.GetY()
	return ((pointC.GetX() - bX) * (pointA.GetY() - bY)) - ((pointC.GetY() - bY) * (pointA.GetX() - bX))
}