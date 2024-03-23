package main

import (
  "fmt"
  "os"
  "bufio"
  "math/rand"
  "strings"
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

func main() {
  userAgents := loadAgents()

  c := colly.NewCollector(
    // colly.AllowedDomains("toscrape.com"),
    colly.MaxDepth(5),
    colly.UserAgent(userAgents[rand.Intn(len(userAgents))]),
  )

  c.OnHTML("p", func(e *colly.HTMLElement){
    fmt.Println(strings.TrimSpace(e.Text))
  })

  c.Visit("https://www.example.com")
}
