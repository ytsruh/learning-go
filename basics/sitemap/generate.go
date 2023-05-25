package sitemap

import (
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func Generate(url string, depth int) (success bool, err error) {
	if url == "" || !strings.HasPrefix(url, "http") {
		return false, errors.New("error: url must start with http and be real url")
	}

	if depth > 10 {
		depth = 10
	}
	log.Println("Starting...")
	pages := bfs(url, depth)

	toXML := urlset{
		Xmlns: xmlns,
	}
	log.Println("Generating XML...")
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}

	log.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		log.Println(err)
		return false, errors.New("error: problem creating XML")
	}
	log.Println("Completed")

	return true, nil
}

func get(urlStr string) []string {
	log.Println("Getting url...")
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(reader io.Reader, base string) []string {
	links, _ := Parse(reader)

	var hrefs []string

	for _, l := range links {
		switch {
		// Find relative paths
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		// Find absolute paths
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		// Ignore everything else
		default:
			//log.Println("Ignoring: " + l.Href)
		}
	}
	return hrefs
}

// Filter takes in slice of strings but also a function that returns bool. Filter function returns slice of string
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	// Use empty struct here because it takes up less memory than a bool
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		// move nextQueue to queue and replace nextQueue with empty map
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			// If the url is in the seen map then ok is true as it means the page was visited
			if _, ok := seen[url]; ok {
				continue
			}
			// If url is not in map then need to visit and get new links
			// struct{}{} is defining the type of an empty struct and then instantiating it
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	// Create new variable with an empty slice and max length so that memory is allocated ahead of time
	seenUrls := make([]string, 0, len(seen))
	for url := range seen {
		seenUrls = append(seenUrls, url)
	}
	return seenUrls
}
