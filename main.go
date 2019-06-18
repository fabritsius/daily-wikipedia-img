package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// URI – Wikipedia Daily Image API Link
const URI string = "https://en.wikipedia.org/w/api.php?action=featuredfeed&format=json&feed=potd"

var wg sync.WaitGroup

// main function sets server handlers and starts the server
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/styles.css", stylesHandler)
	http.HandleFunc("/sw.js", serviceWorkerHandler)
	http.HandleFunc("/pull-reload.js", reloadScriptHandler)
	http.HandleFunc("/manifest.json", manifestHandler)
	http.HandleFunc("/icons/", iconsHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	appengine.Main()
}

// indexHandler function handles path "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	wikiData := GetWikiData(URI, r)
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, wikiData)
	fmt.Println(r.Method, r.URL)
}

// stylesHandler function handles path "/styles.css"
func stylesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "css/main.css")
	fmt.Println(r.Method, r.URL)
}

// serviceWorkerHandler function handles path "/sw.js"
func serviceWorkerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "js/sw.js")
	fmt.Println(r.Method, r.URL)
}

func reloadScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "js/pull-reload.js")
}

// manifestHandler function handles path "/manifest.js"
func manifestHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "manifest.json")
	fmt.Println(r.Method, r.URL)
}

// iconsHandler function handles all "/icons/" paths
func iconsHandler(w http.ResponseWriter, r *http.Request) {
	iconPath := r.URL.Path[1:]
	http.ServeFile(w, r, iconPath)
	fmt.Println(r.Method, r.URL)
}

// faviconHandler function handles path "/favicon.ico"
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "icons/favicon.ico")
}

// WikiData – representation of Wikipedia's Daily Posts Data
type WikiData struct {
	Title       string      `xml:"channel>title"`
	Description string      `xml:"channel>description"`
	Items       []DailyItem `xml:"channel>item"`
}

// GetWikiData function fetches and returns "picture of the day" data
func GetWikiData(link string, r *http.Request) WikiData {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, _ := client.Get(URI)
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var data WikiData
	xml.Unmarshal(bytes, &data)

	// Reverse order of elements in DailyItems slice
	for i := len(data.Items)/2 - 1; i >= 0; i-- {
		j := len(data.Items) - 1 - i
		data.Items[i], data.Items[j] = data.Items[j], data.Items[i]
	}

	for i := range data.Items {
		wg.Add(1)
		go data.Items[i].FillWithValues()
	}
	wg.Wait()
	return data
}

// DailyItem – representation of Wikipedia Post Data
type DailyItem struct {
	Title       template.HTML
	Day         string 			`xml:"title"`
	Description template.HTML
	ImgSrc      string
	Link        string 			`xml:"link"`
	HTML        string 			`xml:"description"`
}

// FillWithValues function fills DailyItem structure with values from HTML
func (d *DailyItem) FillWithValues() {
	defer wg.Done()
	tokenizer := html.NewTokenizer(strings.NewReader(d.HTML))
	recordingDescription := false
	// Leave only the date for each post
	dayTitleWords := strings.Split(d.Day, " ")
	dtLen := len(dayTitleWords)
	d.Day = fmt.Sprintf("%s %s", dayTitleWords[dtLen-2], dayTitleWords[dtLen-1])
	// Look through HTML to find relevant pieces of information
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				fmt.Println(err)
			}
			return

		case html.StartTagToken:
			tag, hasAttr := tokenizer.TagName()
			sTag := string(tag)

			if sTag == "p" {
				recordingDescription = true
				d.Description += template.HTML(tokenizer.Raw())
			
			} else if recordingDescription {
				if sTag == "a" && hasAttr {
					attrs := getAttrVals(tokenizer)
					d.Description += buildWikipediaLink(attrs)
				} else {
					d.Description += template.HTML(tokenizer.Raw())
				}

			} else if hasAttr {
				attrs := getAttrVals(tokenizer)

				if isImage(sTag, attrs) {
					d.Title = template.HTML(attrs["title"])
					tokenizer.Next()
					imgAttrs := getAttrVals(tokenizer)
					d.ImgSrc = https(imgAttrs["src"])
					if srcset, ok := imgAttrs["srcset"]; ok {
						srcSets := strings.Fields(srcset)
						for i, item := range srcSets {
							if item == "2x" {
								if len(srcSets[i-1]) > 54 {
									d.ImgSrc = https(srcSets[i-1])
								}
							}
						}
					}

				} else if isVideo(sTag, attrs) {
					tokenizer.Next()
					imgAttrs := getAttrVals(tokenizer)
					d.ImgSrc = https(imgAttrs["src"])
					tokenizer.Next()
					videoAttrs := getAttrVals(tokenizer)
					d.Title = template.HTML(
						fmt.Sprintf("<a class='video-uri color-hover' target='_blank'" +
									"href='%s'>Play media</a>", videoAttrs["href"]))
								
				} else if isHeader(sTag, attrs) {
					tokenizer.Next()
					d.Title = template.HTML(string(tokenizer.Text()))
				}
			}
		
		case html.EndTagToken:
			if recordingDescription {
				tag, _ := tokenizer.TagName()
				if string(tag) == "p" {
					recordingDescription = false
				}
				d.Description += template.HTML(tokenizer.Raw())
			}

		case html.TextToken:
			if recordingDescription {
				d.Description += template.HTML(tokenizer.Raw())
			}
		}
	}
}

// getAttrVals function returns map with tag's attributes
func getAttrVals(t *html.Tokenizer) map[string]string {
	var result = make(map[string]string)
	for {
		attrName, attrVal, moreAttr := t.TagAttr()
		result[string(attrName)] = string(attrVal)
		if !moreAttr {
			return result
		}
	}
}

// isImage return true if tag is an image in wiki-xml
func isImage(tag string, attrs map[string]string) bool {
	return tag == "a" && attrs["class"] == "image"
}

// isVideo returns true if tag is a video in wiki-xml
func isVideo(tag string, attrs map[string]string) bool {
	return tag == "div" && attrs["class"] == "PopUpMediaTransform"
}

// isHeader returns true if tag is a header in wiki-xml
func isHeader(tag string, attrs map[string]string) bool {
	return tag == "h1" && attrs["id"] == "firstHeading"
}

// https returns a secure version of URI
func https(uri string) string {
	return "https:" + uri
}

// buildWikipediaLink converts relative URI to Wikipedia URI
func buildWikipediaLink(attrs map[string]string) template.HTML {
	wikipediaRoot := "https://en.wikipedia.org"
	linkTemplate := "<a href='%s' title='%s' target='_blank'>"
	return template.HTML(
		fmt.Sprintf(linkTemplate, wikipediaRoot + attrs["href"], attrs["title"]))
}