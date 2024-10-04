package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	var URL string = "https://www.gutenberg.org/files/1727/1727-0.txt"
	bookText, err := downloadBook(URL)
	if err != nil {
		panic(err)
	}

	numArrays := 50
	arrayLenBase := 1000
	dataSets := make([][]string, 0)
	for i := range 10 {
		newSet, err := buildDataSets(bookText, numArrays, arrayLenBase*(1+i))
		if err != nil {
			panic(err)
		}
		dataSets = append(dataSets, newSet)
	}
	//fmt.Println(dataSets[0][0])

}

func findChar(s string, c byte) int {
	for i := range s {
		if s[i] == c {
			return i
		}
	}
	return len(s)
}

func buildDataSets(text string, num int, size int) ([]string, error) {
	dataSet := make([]string, 0)
	l := len(text)
	for i := 0; i < num; i++ {
		start := rand.Intn(l - size - 1)
		dataSet = append(dataSet, text[start:start+size])
	}
	//fmt.Println(len(dataSet))
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

	start := "*** START OF THE PROJECT GUTENBERG EBOOK THE ODYSSEY ***"
	end := "*** END OF THE PROJECT GUTENBERG EBOOK THE ODYSSEY ***"
	startIndex := strings.Index(fullText, start)
	endIndex := strings.Index(fullText, end)

	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return "", fmt.Errorf("could not find start or end markers in the text")
	}

	bookText := fullText[startIndex+len(start) : endIndex]

	re := regexp.MustCompile("[[:^ascii:]]")
	bookText = re.ReplaceAllLiteralString(bookText, "")

	return bookText, nil
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

	   	file.Write([]byte(fullText))
	   	return file, nil */
}
