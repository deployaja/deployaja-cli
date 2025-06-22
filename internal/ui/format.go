package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// FormatCurrency formats a float64 value as Indonesian Rupiah
func FormatCurrency(amount float64) string {
	// Convert to string with 0 decimal places since IDR typically doesn't use decimals
	amountStr := fmt.Sprintf("%.0f", amount)

	// Add thousand separators (dots in Indonesian format)
	// Start from the right and add dots every 3 digits
	runes := []rune(amountStr)
	var result []rune

	for i, r := range runes {
		if i > 0 && (len(runes)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, r)
	}

	return "Rp " + string(result)
}

func FormatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("2006-01-02 15:04:05")
}

// TableColumn represents a column in a table
type TableColumn struct {
	Header string
	Width  int
}

// stripANSI removes ANSI color codes from a string
func stripANSI(s string) string {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// getStringWidth returns the visual width of a string (ignoring ANSI codes)
func getStringWidth(s string) int {
	return len(stripANSI(s))
}

// FormatTable creates a formatted table similar to Kubernetes output
func FormatTable(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	// Calculate column widths based on visual width (ignoring ANSI codes)
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = getStringWidth(header)
	}

	// Find maximum width for each column
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				cellWidth := getStringWidth(cell)
				if cellWidth > colWidths[i] {
					colWidths[i] = cellWidth
				}
			}
		}
	}

	var result strings.Builder

	// Print header
	result.WriteString(formatRow(headers, colWidths))
	result.WriteString(formatSeparator(colWidths))

	// Print rows
	for _, row := range rows {
		result.WriteString(formatRow(row, colWidths))
	}

	return result.String()
}

func formatRow(cells []string, widths []int) string {
	var result strings.Builder
	for i, cell := range cells {
		if i < len(widths) {
			// Calculate padding needed
			cellWidth := getStringWidth(cell)
			padding := widths[i] - cellWidth

			// Add the cell content
			result.WriteString(cell)

			// Add padding
			if padding > 0 {
				result.WriteString(strings.Repeat(" ", padding))
			}

			// Add column separator with better spacing
			if i < len(cells)-1 {
				result.WriteString("   ")
			}
		}
	}
	result.WriteString("\n")
	return result.String()
}

func formatSeparator(widths []int) string {
	var result strings.Builder
	for i, width := range widths {
		result.WriteString(strings.Repeat("-", width))
		if i < len(widths)-1 {
			result.WriteString("   ")
		}
	}
	result.WriteString("\n")
	return result.String()
}
