package main

import (
  "fmt"
  "os"
  "bufio"
  "github.com/gocolly/colly/v2"
  "math/rand"
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

  for i := 0; i < 5; i++ {
    fmt.Println(userAgents[rand.Intn(len(userAgents))])
  }

  c := colly.NewCollector(
    colly.AllowedDomains("toscrape.com"),
    colly.MaxDepth(5),
    // colly.UserAgent(randomAgent()),
  )
  c.Visit("www.example.com")
}
