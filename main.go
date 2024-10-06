package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Brute force search Algo
// finds char c in a string
// requires string be only ascii characters
// verification in a main_test.go
func findChar(s string, c byte) int {
	for i := range s {
		if s[i] == c {
			return i
		}
	}
	return len(s)
}

func main() {

	var URL string = "https://www.gutenberg.org/files/1727/1727-0.txt"
	bookText, err := downloadBook(URL)
	if err != nil {
		panic(err)
	}

	numArrays := 50
	arrayLenBase := 1000
	numLengths := 10
	dataSets, _ := buildDataSets(&bookText, numArrays, arrayLenBase, numLengths)

	testChars := []byte{'e', 'm', 'Q', '%'}
	var sets []statsSet
	for _, c := range testChars {
		sets = append(sets, MinMaxAvg(dataSets, c))
	}

	plotStat(sets, "worst")
	plotStat(sets, "best")
	plotStat(sets, "avg")

}

func plotStat(plotData []statsSet, stat string) {

	p := plot.New()

	plotArgs := make([]interface{}, 0)
	for i := range plotData {
		pts := make(plotter.XYs, len(plotData[i].stats))
		for j := range plotData[i].stats {
			pts[j].X = float64(plotData[i].stats[j].l)
			switch stat {
			case "worst":
				pts[j].Y = float64(plotData[i].stats[j].max)
			case "best":
				pts[j].Y = float64(plotData[i].stats[j].min)
			case "avg":
				pts[j].Y = float64(plotData[i].stats[j].avg)
			}
			//pts[j].Y /= pts[j].X
		}
		plotArgs = append(plotArgs, string(plotData[i].char))
		plotArgs = append(plotArgs, pts)
	}

	p.Title.Text = fmt.Sprintf("Plot of %s", stat)
	p.X.Label.Text = "Array Size"
	p.Y.Label.Text = "Normalized Run Time"

	err := plotutil.AddLinePoints(p, plotArgs...)
	if err != nil {
		panic(err)
	}

	p.Legend.YOffs = vg.Length(.5 * vg.Inch)
	//p.X.Padding.Right = vg.Length(1.5 * vg.Inch)
	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, fmt.Sprintf("%s.png", stat)); err != nil {
		panic(err)
	}
}

type stats struct {
	min int
	max int
	avg int
	l   int
}

type statsSet struct {
	char  byte
	stats []stats
}

// finds the worst,best and average search case for char c in the data set
func MinMaxAvg(s [][]string, c byte) statsSet {
	var results statsSet
	results.char = c
	for i := range len(s) { //number of array sizes
		max := 0
		min := len(s[i][0])
		total := 0
		numArrays := len(s[i])
		for j := range numArrays { //number of arrays of one size
			res := findChar(s[i][j], c)
			if res > max {
				max = res
			}
			if res < min {
				min = res
			}
			total += res
		}
		avg := total / numArrays
		results.stats = append(results.stats, stats{min: min, max: max, avg: avg, l: len(s[i][0])})
	}
	return results
}

// for each Array Size, the book is split into
func buildDataSets(text *string, numArrays int, baseSize int, numSizes int) ([][]string, error) {

	//new random source w/ fixed seed
	r := rand.New(rand.NewSource(int64('o') + int64('d') + int64('y') + int64('s') + int64('s') + int64('e') + int64('y')))

	dataSet := make([][]string, 0)

	l := len(*text)
	segSize := l / numArrays

	for j := range numSizes {
		arraySize := baseSize * (j + 1)
		dataSet = append(dataSet, make([]string, 0))
		for i := range numArrays {
			segStart := i * segSize
			start := segStart + r.Intn(segSize-arraySize)
			subString := (*text)[start : start+arraySize]
			dataSet[j] = append(dataSet[j], subString)
		}
	}
	//fmt.Println(dataSet[0][0])
	return dataSet, nil
}

func downloadBook(bookUrl string) (string, error) {
	resp, err := http.Get(bookUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	var body *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var extractText func(*html.Node) string
	extractText = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		var result string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result += extractText(c)
		}
		return result
	}

	fullText := extractText(body)

	//if we need to save the text.
	/* 	s := "The Project Gutenberg eBook of"
	   	titleStart := strings.Index(fullText, s)
	   	titleEnd := strings.Index(fullText, "\n")
	   	title := fullText[titleStart+len(s) : titleEnd]
	   	title = strings.ReplaceAll(title, " ", "")
	   	title = strings.ReplaceAll(title, ",", "")

	   	fileName := title + ".txt"
	   	file, err := os.Create(fileName)
	   	if err != nil {
	   		return nil, err
	   	}
	   	defer file.Close()

	   	file.Write([]byte(fullText))*/

	start := "*** START OF THE PROJECT GUTENBERG"
	end := "*** END OF THE PROJECT GUTENBERG"
	startIndex := strings.Index(fullText, start)
	endIndex := strings.Index(fullText, end)

	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return "", fmt.Errorf("could not find start or end markers in the text")
	}

	bookText := fullText[startIndex+len(start) : endIndex]

	re := regexp.MustCompile("[[:^ascii:]]")
	bookText = re.ReplaceAllLiteralString(bookText, "")

	return bookText, nil
}
