package spiders

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nozzle/throttler"

	"github.com/walkmiao/fake-useragent/downloader"
	"github.com/walkmiao/fake-useragent/scheduler"
	"github.com/walkmiao/fake-useragent/setting"
	"github.com/walkmiao/fake-useragent/useragent"
)

type Spider struct {
	Attribute
	FullUrl string
}

type Attribute struct {
	Tag      string
	Category string
	Page     int
}

var urlAttributeResults = make(map[string]Attribute)

func NewBrowserSpider() *Spider {
	return &Spider{}
}

func (a *Attribute) GetSpider() *Spider {
	return &Spider{
		Attribute: Attribute{
			Tag:      a.Tag,
			Category: a.Category,
			Page:     a.Page,
		},
		FullUrl: fmt.Sprintf(setting.BROWSER_URL, a.Tag, a.Category, a.Page),
	}
}

func (s *Spider) AppendBrowser(maxPage int) {
	for tag, categories := range setting.BrowserUserAgentMaps {
		for _, category := range categories {
			for page := 1; page <= maxPage; page++ {
				attribute := Attribute{Tag: tag, Category: category, Page: page}
				urlAttributeResults[attribute.GetSpider().FullUrl] = attribute
				scheduler.AppendUrl(attribute.GetSpider().FullUrl)
			}
		}
	}
}

func (s *Spider) StartBrowser(delay time.Duration, timeout time.Duration) {
	count := scheduler.CountUrl()
	th := throttler.New(5, count)
	for i := 0; i <= count; i++ {
		go func() {
			var (
				resp *http.Response
				doc  *goquery.Document
				err  error
			)
			defer th.Done(err)

			if url := scheduler.PopUrl(); url != "" {
				vDownloader := downloader.Download{Delay: delay, Timeout: timeout}
				resp, err = vDownloader.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				doc, err = goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					return
				}

				doc.Find("td.useragent a").Each(func(i int, selection *goquery.Selection) {
					if value := selection.Text(); value != "" {
						useragent.UA.Set(urlAttributeResults[url].Category, value)
					}
				})
			}
		}()

		th.Throttle()
	}
}
