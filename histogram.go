package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

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
	// Convert to string and get first non-zero, non-decimal digit
	str := fmt.Sprintf("%f", num)
	for _, c := range str {
		if c != '0' && c != '.' {
			digit, _ := strconv.Atoi(string(c))
			return digit
		}
	}
	return 0
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

	// Count frequencies of first digits
	digitCounts := make(map[int]int)
	totalCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := sizeRegex.FindStringSubmatch(line)
		if len(matches) >= 3 {
			size, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				continue
			}

			// Convert MiB to KiB if necessary
			if matches[2] == "MiB" {
				size *= 1024
			}

			digit := getFirstDigit(size)
			if digit > 0 {
				digitCounts[digit]++
				totalCount++
			}
		}
	}

	// Create the bar chart
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
	)

	// Prepare data for the chart
	digits := make([]int, 9)
	actualFreq := make([]float64, 9)
	expectedFreq := make([]float64, 9)

	for i := 0; i < 9; i++ {
		digits[i] = i + 1
		count := float64(digitCounts[i+1])
		actualFreq[i] = (count / float64(totalCount)) * 100
		expectedFreq[i] = benfordExpected[i+1]
	}

	// Convert digits to strings for x-axis
	xAxis := make([]string, 9)
	for i := range digits {
		xAxis[i] = fmt.Sprintf("%d", digits[i])
	}

	// Add the data series
	bar.SetXAxis(xAxis)
	bar.AddSeries("Actual", generateBarItems(actualFreq))
	bar.AddSeries("Expected (Benford's Law)", generateBarItems(expectedFreq))

	// Print numerical comparison
	fmt.Println("\nNumerical Comparison:")
	fmt.Println("Digit | Actual % | Expected %")
	// fmt.Println("-" * 30)
	for i := 0; i < 9; i++ {
		fmt.Printf("%5d | %8.1f | %8.1f\n", 
			digits[i], 
			actualFreq[i], 
			expectedFreq[i])
	}

	// Calculate chi-square statistic
	chiSquare := 0.0
	for i := 0; i < 9; i++ {
		expected := (benfordExpected[i+1] / 100) * float64(totalCount)
		observed := float64(digitCounts[i+1])
		chiSquare += math.Pow(observed-expected, 2) / expected
	}
	// Add chi-square interpretation
	fmt.Printf("\nChi-square statistic: %.2f\n", chiSquare)
	fmt.Println("Interpretation:")
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
	// Save the chart to an HTML file
	f, _ := os.Create("benford.html")
	bar.Render(f)
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i := 0; i < len(values); i++ {
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}
