package main

import ("fmt"
        "net/http"
        "io/ioutil"
        "io"
        "encoding/xml"
        "golang.org/x/net/html"
        "html/template"
        "strings"
        "sync")

const PORT string = ":5000"
const URI string = "https://en.wikipedia.org/w/api.php?action=featuredfeed&format=json&feed=potd"

var wg sync.WaitGroup

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/styles.css", stylesHandler)
  http.HandleFunc("/favicon.ico", favIconHandler)
  http.HandleFunc("/icons/github.png", gitIconHandler)
  fmt.Println("...Serving on port", PORT)
  http.ListenAndServe(PORT, nil)
}

// Handles path "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
  wikiData := GetWikiData(URI)
  t, _ := template.ParseFiles("index.html")
  t.Execute(w, wikiData)
  fmt.Println(r.Method, r.URL)
}

// Handles path "/styles.css"
func stylesHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "styles.css")
  fmt.Println(r.Method, r.URL)
}

// Handles path "/favicon.ico"
func favIconHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "icons/wikipedia.ico")
  fmt.Println(r.Method, r.URL)
}

// Handles path "/github.png"
func gitIconHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "icons/github.png")
  fmt.Println(r.Method, r.URL)
}

type WikiData struct {
  Title       string      `xml:"channel>title"`
  Description string      `xml:"channel>description"`
  Items       []DailyItem `xml:"channel>item"`
}

// Fetch and return "picture of the day" data
func GetWikiData(link string) WikiData {
  resp, _ := http.Get(URI)
  bytes, _ := ioutil.ReadAll(resp.Body)
  resp.Body.Close()

  var data WikiData
  xml.Unmarshal(bytes, &data)

  // Reverse order of elements in DailyItems slice
  for i := len(data.Items)/2-1; i >= 0;i-- {
    j := len(data.Items)-1-i
    data.Items[i], data.Items[j] = data.Items[j], data.Items[i]
  }

  for i, _ := range data.Items {
    wg.Add(1)
    go data.Items[i].FillWithValues()
  }
  wg.Wait()
  return data
}

type DailyItem struct {
  Title       string
  Day         string `xml:"title"`
  Description string
  ImgSrc      string
  Link        string `xml:"link"`
  HTML        string `xml:"description"`
}

// Fill DailyItem structure with values from HTML
func (d *DailyItem) FillWithValues() {
  defer wg.Done()
  tokenizer := html.NewTokenizer(strings.NewReader(d.HTML))
  isFirstParagraph := true
  recordingDescription := false
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
        s_tag := string(tag)
        if hasAttr {
          attrs := getAttrVals(tokenizer)
          if s_tag == "a" {
            if attrs["class"] == "image" {
              d.Title = attrs["title"]
              tokenizer.Next()
              imgAttrs := getAttrVals(tokenizer)
              d.ImgSrc = "https:" + imgAttrs["src"]
              if srcset, ok := imgAttrs["srcset"]; ok {
                srcSets := strings.Fields(srcset)
                for i, item := range srcSets {
                  if item == "2x" {
                    if len(srcSets[i-1]) > 54 {
                      d.ImgSrc = "https:" + srcSets[i-1]
                    }
                  }
                }
              }
            }
          } else if s_tag == "video" {
            // [WARNING] unfinished code
            // since videos are rare among daily images
            // this part is hard to improve at the moment
            d.ImgSrc = "https:" + attrs["poster"]
          } else if s_tag == "h1" && attrs["id"] == "firstHeading" {
            tokenizer.Next()
            d.Title = string(tokenizer.Text())
          } 
        } else if s_tag == "p" && !recordingDescription {
          recordingDescription = true
        }
      case html.EndTagToken:
        tag, _ := tokenizer.TagName()
        if string(tag) == "p" && isFirstParagraph {
          isFirstParagraph = false
        }
      case html.TextToken:
        if recordingDescription && isFirstParagraph {
          d.Description += string(tokenizer.Text())
        }
    }
  }
}

// Return map with tag's attributes
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

// Pretty print for Daily Items (unused at the moment)
func (d DailyItem) String() string {
  return fmt.Sprintf(`
  Title         %s
  Day           %s
  Description   %s
  ImgSrc        %s
  Link          %s
  `, d.Title, d.Day, d.Description, d.ImgSrc, d.Link)
}