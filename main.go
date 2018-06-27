package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const metubeLiveURL = "https://www.metube.id/live"

// PlaylistItem good
type PlaylistItem struct {
	Title string
	URL   string
	Logo  string
	Error bool
}

func main() {

	// Instantiate default collector
	c := colly.NewCollector()

	tvList := []string{}

	// On every a element which has href attribute call callback
	c.OnHTML(".livetv-thumbnail a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		tvList = append(tvList, link)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(metubeLiveURL)

	playList := []PlaylistItem{}

	for _, url := range tvList {
		item := getPlayListItem(url)
		playList = append(playList, item)
	}

	m3uString := generateM3UText(playList)

	fmt.Println(m3uString)
}

func getPlayListItem(url string) PlaylistItem {
	item := PlaylistItem{}
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		m3u := m3u8URLExtract(string(r.Body))
		if m3u == "" {
			item.Error = true
		} else {
			item.URL = m3u
		}

	})

	c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
		logo := e.Attr("content")
		item.Logo = logo
	})

	c.OnHTML("meta[property='og:title']", func(e *colly.HTMLElement) {
		title := e.Attr("content")
		item.Title = title
	})

	c.Visit(url)

	return item
}

func m3u8URLExtract(content string) string {
	r, _ := regexp.Compile(`(http.*?\.m3u8)`)
	result := r.FindString(content)
	return result
}

func generateM3UText(playlist []PlaylistItem) string {

	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	for _, item := range playlist {
		if item.Error == true {
			continue
		}

		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("#EXTINF:-1 tvg-logo=\"%s\",%s\n", item.Logo, item.Title))
		sb.WriteString(item.URL)
		sb.WriteString("\n")
	}

	return sb.String()
}
