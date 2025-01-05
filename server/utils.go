package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func GeneratePoints(size, maxX, maxY int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	points := make([][2]int, size)

	for i := 0; i < size; i++ {
		points[i][0] = rand.Intn(maxX)
		points[i][1] = rand.Intn(maxY)
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
func IsPointInCircle(cx, cy, r float64, px, py int) bool {
	// Calculate the distance from the point (px, py) to the center (cx, cy)
	distance := math.Sqrt(math.Pow(float64(px)-cx, 2) + math.Pow(float64(py)-cy, 2))
	
	// Check if the distance is less than or equal to the radius
	return distance <= r
}

// Function to check if Circle A (center cxA, cyA, radius rA) includes Circle B (center cxB, cyB, radius rB)
func IsCircleInCircle(cxA, cyA, rA, cxB, cyB, rB float64) bool {
	// Calculate the distance between the centers of Circle A and Circle B
	distance := math.Sqrt(math.Pow(cxA-cxB, 2) + math.Pow(cyA-cyB, 2))
	
	// Check if the distance + radius of Circle B is less than or equal to radius of Circle A
	return distance + rB <= rA
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
