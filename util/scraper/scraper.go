package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tmc/langchaingo/tools"
)

const ErrScrapingFailed = "scraper could not read URL, or scraping is not allowed for provided URL"

type Scraper struct {
	MaxDepth  int
	Parallels int
	Delay     int64
}

var _ tools.Tool = Scraper{}

func NewScraper(maxDepth ...int) (*Scraper, error) {
	depth := 1
	parallels := 2
	delay := 3

	if len(maxDepth) > 0 {
		depth = maxDepth[0]
	}

	return &Scraper{
		MaxDepth:  depth,
		Parallels: parallels,
		Delay:     int64(delay),
	}, nil
}

func (scraper Scraper) Name() string {
	return "Web Scraper"
}

func (scraper Scraper) Description() string {
	return `
		Web Scraper will scan a url and return the content of the web page.
		Input should be a working url.
	`
}

func (scraper Scraper) Call(ctx context.Context, input string) (string, error) {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	c := colly.NewCollector(
		colly.MaxDepth(scraper.MaxDepth),
		colly.Async(true),
	)

	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: scraper.Parallels,
		Delay:       time.Duration(scraper.Delay) * time.Second,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	var siteData strings.Builder
	homePageLinks := make(map[string]bool)
	scrapedLinks := make(map[string]bool)

	c.OnRequest(func(r *colly.Request) {
		if ctx.Err() != nil {
			r.Abort()
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		currentURL := e.Request.URL.String()

		// Only process the page if it hasn't been visited yet
		if !scrapedLinks[currentURL] {
			scrapedLinks[currentURL] = true

			siteData.WriteString("\n\nPage URL: " + currentURL)

			title := e.ChildText("title")
			if title != "" {
				siteData.WriteString("\nPage Title: " + title)
			}

			description := e.ChildAttr("meta[name=description]", "content")
			if description != "" {
				siteData.WriteString("\nPage Description: " + description)
			}

			siteData.WriteString("\nHeaders:")
			e.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
				siteData.WriteString("\n" + el.Text)
			})

			siteData.WriteString("\nContent:")
			e.ForEach("p", func(_ int, el *colly.HTMLElement) {
				siteData.WriteString("\n" + el.Text)
			})

			if currentURL == input {
				e.ForEach("a", func(_ int, el *colly.HTMLElement) {
					link := el.Attr("href")
					if link != "" && !homePageLinks[link] {
						homePageLinks[link] = true
						siteData.WriteString("\nLink: " + link)
					}
				})
			}
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)

		// Parse the link to get the hostname
		u, err := url.Parse(absoluteLink)
		if err != nil {
			// Handle the error appropriately
			return
		}

		// Check if the link's hostname matches the current request's hostname
		if u.Hostname() != e.Request.URL.Hostname() {
			return
		}

		// Check for redundant pages
		blacklist := []string{
			"login",
			"signup",
			"signin",
			"register",
			"logout",
			"download",
			"redirect",
		}
		for _, item := range blacklist {
			if strings.Contains(u.Path, item) {
				return
			}
		}

		// Normalize the path to treat '/' and '/index.html' as the same path
		if u.Path == "/index.html" || u.Path == "" {
			u.Path = "/"
		}

		// Only visit the page if it hasn't been visited yet
		if !scrapedLinks[u.String()] {
			err := c.Visit(u.String())
			if err != nil {
				siteData.WriteString(fmt.Sprintf("\nError following link %s: %v", link, err))
			}
		}
	})

	err = c.Visit(input)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		c.Wait()
	}

	// Append all scraped links
	siteData.WriteString("\n\nScraped Links:")
	for link := range scrapedLinks {
		siteData.WriteString("\n" + link)
	}

	return siteData.String(), nil
}
