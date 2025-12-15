package gui

import (
	"image/color"
)

// Category color mappings for visual organization
var categoryColors = map[string]color.Color{
	"Python":       color.NRGBA{R: 55, G: 118, B: 171, A: 255}, // Blue
	"JavaScript":   color.NRGBA{R: 240, G: 219, B: 79, A: 255}, // Yellow
	"Go":           color.NRGBA{R: 0, G: 173, B: 216, A: 255},  // Cyan
	"C":            color.NRGBA{R: 85, G: 85, B: 85, A: 255},   // Dark Gray
	"SQL":          color.NRGBA{R: 0, G: 122, B: 204, A: 255},  // Azure
	"HTML":         color.NRGBA{R: 227, G: 79, B: 38, A: 255},  // Orange
	"CSS":          color.NRGBA{R: 41, G: 101, B: 241, A: 255}, // Royal Blue
	"Personal":     color.NRGBA{R: 156, G: 39, B: 176, A: 255}, // Purple
	"Professional": color.NRGBA{R: 76, G: 175, B: 80, A: 255},  // Green
	"Symbols":      color.NRGBA{R: 255, G: 152, B: 0, A: 255},  // Deep Orange
	"General":      color.NRGBA{R: 96, G: 125, B: 139, A: 255}, // Blue Gray
	"DateTime":     color.NRGBA{R: 63, G: 81, B: 181, A: 255},  // Indigo
}

// GetCategoryColor returns the color for a given category
func GetCategoryColor(category string) color.Color {
	if c, ok := categoryColors[category]; ok {
		return c
	}
	// Default color for unknown categories
	return color.NRGBA{R: 158, G: 158, B: 158, A: 255} // Gray
}

// GetAllCategories returns a list of all available categories with their colors
func GetAllCategories() map[string]color.Color {
	return categoryColors
}

// CountByCategory counts expansions in each category
func CountByCategory(expansions []interface{}) map[string]int {
	counts := make(map[string]int)

	for _, exp := range expansions {
		if expMap, ok := exp.(map[string]interface{}); ok {
			if category, ok := expMap["category"].(string); ok && category != "" {
				counts[category]++
			} else {
				counts["Uncategorized"]++
			}
		}
	}

	return counts
}
