package crawler

import (
  "fmt"
  "os"
  "bufio"
  "math/rand"
  "strings"
  "time"
  "regexp"
  "net/http"
  "github.com/gocolly/colly/v2"
)

func loadAgents() []string {
  file, _ := os.Open("user-agents.txt")
  s := bufio.NewScanner(file)
  userAgents := []string{}

  for s.Scan() {
    userAgents = append(userAgents, s.Text())
  }
  return userAgents
}

type Crawler struct {
  Collector *colly.Collector
  URLs []string
  Pattern string
}

func setupCollector() *colly.Collector {
  userAgents := loadAgents()

  c := colly.NewCollector(
    colly.MaxDepth(5),
  )

  c.OnRequest(func (r *colly.Request) {
    // on each request set a random user agent
    r.Headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
  })

  // c.OnHTML("p", func(e *colly.HTMLElement){
  //   fmt.Println(strings.TrimSpace(e.Text))
  // })

  // random delay on accesses
  c.Limit(&colly.LimitRule{
    DomainGlob: "*",
    Parallelism: 2,
    RandomDelay: 50 * time.Millisecond,
  })

  return c
}

func (c *Crawler) SetLinkPattern(pattern string) {
  c.Collector.OnHTML("a", func(e *colly.HTMLElement){
    nextPage := e.Request.AbsoluteURL(e.Attr("href"))
    b, _ := regexp.MatchString(pattern, nextPage)
    if b {
      c.Visit(nextPage)
    }
  }) 
}

func (c *Crawler) Init() {
  c.Collector = setupCollector()
  c.URLs = []string{}
  c.Pattern = ""
}

func NewCrawler() *Crawler {
  c := &Crawler{}
  c.Init()
  return c
}

func (c *Crawler) Visit(url string) {
  c.Collector.Visit(url)
}

func trimURL(url string) string {
  m := regexp.MustCompile(`\.([^.].*)`)
  url = m.FindStringSubmatch(url)[0]
  url, _ = strings.CutPrefix(url, ".")
  url = strings.Split(url, "/")[0]
  return url
}

func TestResponse(url string) bool {
  response, err := http.Get(url)
  if err != nil {
    return false
  }
  return response.StatusCode < 400
}

func Crawl(url string) {
  c := NewCrawler()

  fmt.Println(TestResponse(url))

  // TODO: throw error if link does not contain http
  trimmedURL:= trimURL(url)

  c.SetLinkPattern(trimmedURL)
  c.Visit(url)
}

