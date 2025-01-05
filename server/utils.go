package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func GeneratePoints(size int, maxX, maxY float64) [][2]float64 {
	rand.Seed(time.Now().UnixNano())
	points := make([][2]float64, size)

	for i := 0; i < size; i++ {
		points[i][0] = rand.Float64() * maxX
		points[i][1] = rand.Float64() * maxY
	}

	return points;
}


func RandomColor() string {
	rand.Seed(time.Now().UnixNano())
	red := 128 + rand.Intn(128)
	green := 128 + rand.Intn(128)
	blue := 128 + rand.Intn(128)

	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Function to check if a point (px, py) is inside a circle with center (cx, cy) and radius r
func IsPointInCircle(cx, cy, r, px, py float64) bool {
	// Calculate the distance from the point (px, py) to the center (cx, cy)
	distance := math.Sqrt(math.Pow(px-cx, 2) + math.Pow(py-cy, 2))
	
	// Check if the distance is less than or equal to the radius
	return distance <= r
}

// Function to filter points that are inside a circle
func FilterPointsInCircle(cx, cy, r float64, points [][2]float64) [][2]float64 {
	var result [][2]float64
	
	// Iterate over the array of points
	for _, point := range points {
		px, py := point[0], point[1]
		// If the point is inside the circle, add it to the result
		if IsPointInCircle(cx, cy, r, px, py) {
			result = append(result, point)
		}
	}
	
	return result
}


// Function to subtract array b from array a
func SubtractArrays(a, b [][2]float64) [][2]float64 {
	// Create a map to store the count of elements in array b
	bMap := make(map[[2]float64]int)
	
	// Populate the map with elements from array b
	for _, point := range b {
		bMap[point]++
	}
	
	// Create a result array to store elements from a that are not in b
	var result [][2]float64
	
	// Iterate over array a and add elements to result that are not in b
	for _, point := range a {
		if bMap[point] > 0 {
			// Decrease the count of this element in bMap, indicating it's removed
			bMap[point]--
		} else {
			// If the element is not in b, add it to the result
			result = append(result, point)
		}
	}
	
	return result
}
