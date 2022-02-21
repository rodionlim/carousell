// Package carousell provides primitives for querying Carousell (Singapore)
// and parsing the listings programatically.
//
// It also provides a simple caching
// mechanism for users to store state on the listings, for example, only
// caching after post-processing of a listing is completed successfully.
package carousell

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rodionlim/carousell/library/log"
	"golang.org/x/net/html"
)

const ENDPOINT = "https://www.carousell.sg"

// A Req is a structure that consists of relevant
// parameters to encapsulate a carousell GET request
type Req struct {
	endpoint   string
	searchTerm string
	queryParam map[string]string
}

func (r *Req) GetSearchTerm() string {
	return r.searchTerm
}

// Creates a new carousell blank Request.
// Modifiers, e.g. WithSearch, are used to add search criteria, before sending out the Request via Get.
func NewReq(opts ...func(*Req)) *Req {
	r := &Req{
		endpoint:   ENDPOINT,
		queryParam: make(map[string]string),
	}
	for _, f := range opts {
		f(r)
	}
	return r
}

// WithSearch adds a search term to a carousell Request.
func WithSearch(searchTerm string) func(r *Req) {
	return func(r *Req) {
		r.searchTerm += fmt.Sprintf("%q", searchTerm)
	}
}

// WithPriceFloor adds a filter for price floor.
func WithPriceFloor(px int) func(r *Req) {
	return func(r *Req) {
		r.queryParam["price_start"] = fmt.Sprint(px)
	}
}

func WithPriceCeil(px int) func(r *Req) {
	return func(r *Req) {
		r.queryParam["price_end"] = fmt.Sprint(px)
	}
}

// WithRecent ensures that only latest listings are being queried.
func WithRecent(r *Req) {
	r.queryParam["addRecent"] = "true"
	r.queryParam["sort_by"] = "3"
}

// Get gets and parse carousell listing based on user parameters.
// BUG: Get occassionally (very rare) fetches listing that does not take search terms into consideration, this should be handled if polling the function.
func (r *Req) Get() ([]Listing, error) {
	logger := log.Ctx(context.Background())

	if err := r.validate(); err != nil {
		return nil, err
	}
	url, err := url.Parse(fmt.Sprintf("%s/search/%s", r.endpoint, r.searchTerm))
	if err != nil {
		return nil, err
	}

	q := url.Query()
	for k, v := range r.queryParam {
		q.Set(k, v)
	}
	url.RawQuery = q.Encode()

	logger.Infof("Send req [%s]", url.String())
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	logger.Infof("Recv resp status_code[%v] headers[%v]", resp.StatusCode, resp.Header)

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	listings := deref(traverseDivs(node, make([]*Listing, 0)))
	return listings, nil
}

// A Cache is a structure that includes utilities,
// to store state on carousell listings after post-processing.
type Cache struct {
	Alerts map[string]bool
}

// Creates a new Cache of carousell alerts.
func NewCache() *Cache {
	return &Cache{
		Alerts: make(map[string]bool),
	}
}

// Store is a helper to cache listings by their IDs,
// to update the state that a listing has been processed.
func (c *Cache) Store(listings []Listing) {
	for _, listing := range listings {
		c.Alerts[listing.ID] = true
	}
}

// Process accepts a callback that processes a new listing,
// before storing it in the cache.
//
// checkListings is a flag that is used to ensure that listings are valid. When all listings are not cached,
// there is high probability that listings are invalid due to upstream error, since new posts should be infrequent.
// Set it to false to disable the checks and process all listings.
func (c *Cache) ProcessAndStore(listings []Listing, cb func(listing Listing) error, checkListings bool) {
	var toBeAlerted []Listing
	for _, listing := range listings {
		_, exists := c.Alerts[listing.ID]
		if !exists {
			toBeAlerted = append(toBeAlerted, listing)
			c.Alerts[listing.ID] = true
		}
	}
	if len(toBeAlerted) != len(listings) || !checkListings {
		for _, listing := range toBeAlerted {
			cb(listing)
		}
	}
}

// A Listing is a single carousell post.
// It contains all the information extracted off a carousell listing.
type Listing struct {
	Title       string
	Description string
	Price       float64
	Condition   string
	Url         string
	User        string
	Time        string
	ID          string
}

// Print is a summarized output of a listing printed out to console.
func (l *Listing) Print() {
	fmt.Print(l.Sprint())
}

// Sprint returns a summarized output of a listing.
func (l *Listing) Sprint() string {
	return fmt.Sprintf("%s - S$%.0f - %s\n%s\n", l.Title, l.Price, l.Condition, l.Url)
}

// ShortenListings return a summarized output of a list of listings.
func ShortenListings(l []Listing) []string {
	var res []string
	for _, listing := range l {
		res = append(res, (&listing).Sprint())
	}
	return res
}

func deref(l []*Listing) []Listing {
	// Convenience wrapper to convert pointers to slice of structs
	var listings []Listing
	for _, listing := range l {
		listings = append(listings, *listing)
	}
	return listings
}

func traverseText(n *html.Node, texts []string, links []string) ([]string, []string) {
	// Store all text data to texts slice
	if n.Type == html.TextNode {
		if n.Data != "Protection" {
			texts = append(texts, strings.TrimSpace(strings.Trim(n.Data, "\n")))
		}
	}
	// Store all hyper links to links slice
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		texts, links = traverseText(c, texts, links)
	}
	return texts, links
}

func parsePrice(px string) float64 {
	px = strings.Replace(strings.Replace(px, "S$", "", 1), ",", "", -1)
	pxFloat, err := strconv.ParseFloat(px, 64)
	if err != nil {
		return 0
	}
	return pxFloat
}

func nodeToListing(n *html.Node, id string) (*Listing, error) {
	var details []string
	var links []string
	var listing Listing

	itemNode := n.FirstChild
	details, links = traverseText(itemNode, details, links)
	if len(details) >= 6 {
		listing = Listing{
			ID:          id,
			User:        details[0],
			Time:        details[1],
			Title:       details[2],
			Price:       parsePrice(details[3]),
			Description: details[4],
			Condition:   details[5],
		}
	}
	if len(links) >= 2 {
		listing.Url = fmt.Sprintf("%s%s", ENDPOINT, links[1])
	}
	return &listing, nil
}

func traverseDivs(n *html.Node, divs []*Listing) []*Listing {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "data-testid" {
				div, err := nodeToListing(n, attr.Val)
				if err == nil {
					divs = append(divs, div)
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		divs = traverseDivs(c, divs)
	}
	return divs
}

func (r *Req) validate() error {
	if r.searchTerm == "" {
		return errors.New("no search term provided")
	}
	return nil
}
