package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Expected frequencies according to Benford's Law
var benfordExpected = map[int]float64{
	1: 30.1,
	2: 17.6,
	3: 12.5,
	4: 9.7,
	5: 7.9,
	6: 6.7,
	7: 5.8,
	8: 5.1,
	9: 4.6,
}

func getFirstDigit(num float64) int {
	str := fmt.Sprintf("%f", num)
	for _, c := range str {
		if c != '0' && c != '.' {
			digit, _ := strconv.Atoi(string(c))
			return digit
		}
	}
	return 0
}

func analyzeSizeRange(sizes []float64, minSize, maxSize float64) {
	digitCounts := make(map[int]int)
	totalCount := 0

	// Count frequencies for files in the specified range
	for _, size := range sizes {
		if size >= minSize && size < maxSize {
			firstDigit := getFirstDigit(size)
			if firstDigit > 0 {
				digitCounts[firstDigit]++
				totalCount++
			}
		}
	}

	if totalCount == 0 {
		fmt.Printf("\nNo files found between %.0f KiB and %.0f KiB\n", minSize, maxSize)
		return
	}

	// Calculate and print distribution
	fmt.Printf("\nAnalysis for files between %.0f KiB and %.0f KiB (Total files: %d):\n", minSize, maxSize, totalCount)
	fmt.Println("Digit | Actual % | Expected % | Count")
	fmt.Println(strings.Repeat("-", 45))

	chiSquare := 0.0
	for digit := 1; digit <= 9; digit++ {
		actual := float64(digitCounts[digit]) / float64(totalCount) * 100
		expected := benfordExpected[digit]
		fmt.Printf("%5d | %8.1f | %8.1f | %5d\n", digit, actual, expected, digitCounts[digit])

		// Calculate chi-square contribution
		expectedCount := (expected / 100) * float64(totalCount)
		observed := float64(digitCounts[digit])
		chiSquare += math.Pow(observed-expectedCount, 2) / expectedCount
	}

	fmt.Printf("\nChi-square statistic: %.2f (for this range)\n", chiSquare)
	if chiSquare < 15.51 {
		fmt.Println("Result: Good fit to Benford's Law (p > 0.05)")
	} else if chiSquare < 20.09 {
		fmt.Println("Result: Marginal fit to Benford's Law (0.01 < p < 0.05)")
	} else {
		fmt.Println("Result: Poor fit. Does not follow Benford's Law (p < 0.01)")
	}
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	// Open and read the file
	file, err := os.Open("/home/jan/Desktop/benford/pacman.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Regular expression to match sizes
	sizeRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(MiB|KiB)`)
	speedRegex := regexp.MustCompile(`/s`)

	// Store all sizes for multiple analyses
	var sizes []float64
	digitCounts := make(map[int]int)
	totalCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := sizeRegex.FindStringSubmatch(line)
		if len(matches) >= 3 && !speedRegex.MatchString(matches[0]) {
			size, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				continue
			}

			// Convert MiB to KiB if necessary
			if matches[2] == "MiB" {
				size *= 1024
			}

			sizes = append(sizes, size)
			digit := getFirstDigit(size)
			if digit > 0 {
				digitCounts[digit]++
				totalCount++
			}
		}
	}

	// First, show overall analysis
	fmt.Println("=== Overall Analysis ===")
	fmt.Printf("Total number of packages: %d\n", totalCount)
	fmt.Println("\nDigit | Actual % | Expected % | Count")
	fmt.Println(strings.Repeat("-", 45))

	actualFreq := make([]float64, 9)
	expectedFreq := make([]float64, 9)
	chiSquare := 0.0

	for i := 1; i <= 9; i++ {
		count := float64(digitCounts[i])
		actual := (count / float64(totalCount)) * 100
		expected := benfordExpected[i]
		actualFreq[i-1] = actual
		expectedFreq[i-1] = expected

		fmt.Printf("%5d | %8.1f | %8.1f | %5d\n", i, actual, expected, digitCounts[i])

		expectedCount := (expected / 100) * float64(totalCount)
		chiSquare += math.Pow(count-expectedCount, 2) / expectedCount
	}

	fmt.Printf("\nOverall Chi-square statistic: %.2f\n", chiSquare)
	fmt.Println("Critical values (degrees of freedom = 8):")
	fmt.Println("  α = 0.05: 15.51")
	fmt.Println("  α = 0.01: 20.09")
	if chiSquare < 15.51 {
		fmt.Println("Result: Good fit! The data follows Benford's Law (p > 0.05)")
	} else if chiSquare < 20.09 {
		fmt.Println("Result: Marginal fit to Benford's Law (0.01 < p < 0.05)")
	} else {
		fmt.Println("Result: Poor fit. Data does not follow Benford's Law (p < 0.01)")
	}

	// Analyze different size ranges
	fmt.Println("\n=== Analysis by Size Ranges ===")
	analyzeSizeRange(sizes, 0, 100)       // 0-100 KiB
	analyzeSizeRange(sizes, 100, 1000)    // 100-1000 KiB
	analyzeSizeRange(sizes, 1000, 10000)  // 1-10 MiB
	analyzeSizeRange(sizes, 10000, math.MaxFloat64) // >10 MiB

	// Create visualization
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "First Digit Distribution vs Benford's Law",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Frequency (%)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "First Digit",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)

	// Convert digits to strings for x-axis
	xAxis := make([]string, 9)
	for i := range xAxis {
		xAxis[i] = fmt.Sprintf("%d", i+1)
	}

	// Add the data series
	bar.SetXAxis(xAxis)
	bar.AddSeries("Actual", generateBarItems(actualFreq)).
		AddSeries("Expected (Benford's Law)", generateBarItems(expectedFreq))

	// Save the chart to an HTML file
	outputFile := "benford.html"
	f, _ := os.Create(outputFile)
	bar.Render(f)
	f.Close()

	// Get the absolute path of the HTML file
	absPath, _ := os.Getwd()
	htmlPath := "file://" + absPath + "/" + outputFile

	fmt.Printf("\nOpening chart in browser: %s\n", htmlPath)
	err = openBrowser(htmlPath)
	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
		fmt.Printf("Please open %s manually in your web browser\n", outputFile)
	}
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i := 0; i < len(values); i++ {
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}
