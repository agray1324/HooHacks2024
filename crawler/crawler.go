package crawler

import (
  "fmt"
  "os"
  "bufio"
  "math/rand"
  "strings"
  "time"
  "regexp"
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

// func randomAgent(userAgents []string) string {
//   return userAgents[rand.Intn(len(userAgents))]
// }

// func setupCollector() *colly.Collector {
//
// }

func Temp() {
  userAgents := loadAgents()

  c := colly.NewCollector(
    colly.MaxDepth(5),
  )


  c.OnRequest(func (r *colly.Request) {
    r.Headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
  })

  c.Limit(&colly.LimitRule{
    DomainGlob: "*",
    Parallelism: 2,
    RandomDelay: 50 * time.Millisecond,
  })

  c.OnHTML("p", func(e *colly.HTMLElement){
    fmt.Println(strings.TrimSpace(e.Text))
  })
  
  pattern := "*example.com*"

  c.OnHTML("a", func(e *colly.HTMLElement){
    nextPage := e.Request.AbsoluteURL(e.Attr("href"))
    b, _ := regexp.MatchString(pattern, nextPage)
    if b {
      c.Visit(nextPage)
    }
  }) 

  c.Visit("https://www.example.com")

}

