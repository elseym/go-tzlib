package importers

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/elseym/go-tzlib"
)

const (
	modeFromURL  = iota
	modeFromFile = iota

	weightLimit float64 = 16.0
)

// TimeIs parses html
type TimeIs struct {
	mode     int
	resource string
}

// NewTimeIsFromFile returns a new TimeIs importer which will read a
// html file from the filesystem
func NewTimeIsFromFile(htmlfile string) (t TimeIs) {
	return TimeIs{modeFromFile, htmlfile}
}

// NewTimeIsFromURL returns a new TimeIs importer which will download
// html from the provided url
func NewTimeIsFromURL(url string) (t TimeIs) {
	return TimeIs{modeFromURL, url}
}

// Import returns a new Tzlib according to the current TimeIs importer config
func (ti TimeIs) Import() (t *tzlib.Tzlib, err error) {
	var doc *goquery.Document

	doc, err = ti.createDoc()
	if err == nil {
		t = parse(doc.Selection)
	}

	return
}

// parse
func parse(s *goquery.Selection) *tzlib.Tzlib {
	get := func(s *goquery.Selection, e string, f func(*goquery.Selection)) {
		s.Find(e).Each(func(i int, s *goquery.Selection) { f(s) })
	}

	t := tzlib.NewTzlib(time.Now(), parseUpdateTimestamp(s))

	get(s, "div.section", func(s *goquery.Selection) {
		zElem := firstChild(s, "h1")
		zUTC := zElem.Text()

		z := t.AddNewTimezone(zUTC)

		get(s, "div.scloud > ul > li[id^=c]", func(s *goquery.Selection) {
			cElem := firstChild(s, "a")
			cName := cElem.Text()
			cWeight := parseWeight(cElem)

			c := z.AddNewCountry(cName, cWeight)

			get(s, "ul > li", func(s *goquery.Selection) {
				yElem := firstChild(s, "a")
				yName := yElem.Text()
				yWeight := parseWeight(yElem)

				c.AddNewCity(yName, yWeight)
			})
		})
	})

	return t
}

// firstChild
func firstChild(s *goquery.Selection, e string) *goquery.Selection {
	return s.ChildrenFiltered(e).First()
}

// createDoc
func (ti TimeIs) createDoc() (doc *goquery.Document, err error) {
	switch ti.mode {
	case modeFromFile:
		var r io.Reader
		r, err = os.Open(ti.resource)
		if err == nil {
			doc, err = goquery.NewDocumentFromReader(r)
		}
	case modeFromURL:
		doc, err = goquery.NewDocument(ti.resource)
	}

	return
}

// parseWeight
func parseWeight(s *goquery.Selection) (weight float64) {
	if w, err := parseNumAttr(s, "class", "s"); err == nil {
		ww := float64(w)

		if s.HasClass("bold") {
			ww += 0.5
		}

		weight = (weightLimit - math.Min(weightLimit, math.Max(ww, 0))) / weightLimit
	}

	return weight
}

// parseNumAttr
func parseNumAttr(s *goquery.Selection, attr string, prefix string) (n int, err error) {
	if val := s.AttrOr(attr, ""); val != "" {
		for _, v := range strings.Split(val, " ") {
			vt := strings.TrimPrefix(v, prefix)
			if vt != v && strings.HasSuffix(v, vt) {
				return strconv.Atoi(vt)
			}
		}
	}
	return 0, fmt.Errorf("attribute ('%s') value does not contain parsable data", attr)
}

// parseUpdateTimestamp
func parseUpdateTimestamp(s *goquery.Selection) (t time.Time) {
	defer func() {
		if r := recover(); r != nil {
			t = time.Now().Add(24 * time.Hour).UTC()
		}
	}()

	text := s.Find("#cvwr > div > div > p").Last().Text()
	prefix := "The next update is scheduled for "
	layout := "2006-01-02 15:04 MST"
	start := strings.Index(text, prefix) + len(prefix)
	fragment := text[start : start+len(layout)]
	t, _ = time.Parse(layout, fragment)

	return t.UTC()
}
