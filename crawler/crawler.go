package crawler

import (
  "fmt"
  "os"
  "bufio"
  // "io/ioutil"
  "math/rand"
  "strings"
  // "time"
  "regexp"
  "net/http"
  "errors"
  "sync"
  "sort"
  // "io"
  // "bytes"
  // "golang.org/x/net/html"
  "github.com/gocolly/colly/v2"
  "github.com/lithammer/fuzzysearch/fuzzy"
  // usearch "github.com/unum-cloud/usearch/golang"
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
  URL []string
  Content []string
  Pattern string
  mtex sync.Mutex
  Count int
}

func (c *Crawler) Init() {
  c.Collector = setupCollector()
  c.Pattern = ""
  c.Count = 0
}

func NewCrawler() *Crawler {
  c := &Crawler{}
  c.Init()
  return c
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

  // random delay on accesses
  c.Limit(&colly.LimitRule{
    DomainGlob: "*",
    Parallelism: 100,
    // RandomDelay: 50 * time.Millisecond,
  })

  c.Async = true

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

func (c *Crawler) Visit(url string) {
  c.Collector.Visit(url)
}

func (c *Crawler) Wait() {
  c.Collector.Wait()
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

// needs a reference to a waitgroup to parallelize storage
func (c *Crawler) CleanBody(wg *sync.WaitGroup) {
  c.Collector.OnHTML("body", func(e *colly.HTMLElement){
    // create a new thread that stores data
    wg.Add(1)
    go func() {
      defer wg.Done()

      text := strings.TrimSpace(e.Text)
      s := regexp.MustCompile(`\s+`)
      text = s.ReplaceAllString(text, " ")
      text = strings.ReplaceAll(text, "\"", "'")
      // text = strings.ReplaceAll(text, "'", "")
      url := e.Request.URL.String()

      c.mtex.Lock()
      c.Content = append(c.Content, text)
      c.URL = append(c.URL, url)
      c.Count += 1
      c.mtex.Unlock()
    }()
  })
}

// parallelized sort by rank similarity
func Rank(search string, content []string) []int {
  var wg sync.WaitGroup
  ranks := make([]int, len(content))

  for idx, str := range content {
    wg.Add(1)
    go func() {
      defer wg.Done()
      ranks[idx] = fuzzy.RankMatchNormalizedFold(search, str)
    }()
  }

  wg.Wait()

  return ranks
}

func asyncSortRank(sli []string, ranks []int, wg *sync.WaitGroup) {
  defer wg.Done()

  sort.Slice(sli, func(i, j int) bool {
    return ranks[i] > ranks[j]
  })
}

// reorder the content and urls to reflect the best matches
func (c *Crawler) FuzzySearch(search string) {
  ranks := Rank(search, c.Content)

  var wg sync.WaitGroup

  wg.Add(2)
  go asyncSortRank(c.Content, ranks, &wg)
  go asyncSortRank(c.URL, ranks, &wg)
  wg.Wait()

  // sort ranks
  sort.Slice(ranks, func(i, j int) bool {
    return ranks[i] > ranks[j]
  })

  fmt.Println("Top 3 related links:")
  for i := 0; i < 3; i++ {
    fmt.Println("\t", i+1, ".", c.Content[i])
  }
}

func (c *Crawler) Search(search string) string {
  var urls []string

  for idx, con := range c.Content {
    if strings.Contains(con, search) {
      urls = append(urls, c.URL[idx])
    }
  }

  return strings.Join(urls, `\\,\\`)
  // // Create Index
  //   vector_size := 3
  //   vectors_count := 100
  //   conf := usearch.DefaultConfig(uint(vector_size))
  //   index,err := usearch.NewIndex(conf)
  //   if err != nil {
  //       panic("Failed to create Index")
  //   }
  //   defer index.Destroy()
  //
  //   // Add to Index
  //   err = index.Reserve(uint(vectors_count))
  //   for i := 0; i < vectors_count; i++ {
  //       err = index.Add(usearch.Key(i), []float32{float32(i), float32(i+1), float32(i+2)})
  //       if err != nil {
  //           panic("Failed to add")
  //       }
  //   }
  //
  //   // Search
  //   keys, distances, err := index.Search([]float32{0.0, 1.0, 2.0}, 3)
  //   if err != nil {
  //       panic("Failed to search")
  //   }
  //   fmt.Println(keys, distances)
//   url := "https://api.predictionguard.com/chat/completions"
//   method := "POST"
//
//   ind := 0
//
//   for i, item := range c.Content {
//     if strings.Contains(item, "Schulz") {
//       ind = i
//     }
//   }
//
//   fmt.Println(c.Content[ind])
//   fmt.Println(c.URL[ind])
//   question := `### Instruction:\n Read and respond to the following passage with only an integer value between 0 and 100 based on how related the passage is to the following query: ` + search + `. Do not include any other text aside from the integer value. The passage is as follows: \n` + c.Content[ind]
//
//   // question := `You are a ranking machine trying to determine the most relevant website for a given query. You are given the plain text on the website and must rate the entire text passage with one score. You can only respond with a single integer value and you cannot say any other words or characters. What portion of the following passage addresses the topic ` + search + `? Please rate this with an integer ranging from 1 to 100, where 1 is none and 100 is all of it: ` + c.Content[0]
//
//   payload := strings.NewReader(`{
//     "model": "Neural-Chat-7B",
//     "messages": [
//         {
//             "role": "user",
//             "content": "` + question + `"
//         }
//     ]
// }`)
//
//   client := &http.Client {
//   }
//
//   req, err := http.NewRequest(method, url, payload)
//
//   if err != nil {
//     fmt.Println(err)
//     return
//   }
//
//   file, _ := os.Open("prediction_guard.key")
//   s := bufio.NewScanner(file)
//   for s.Scan() {
//     key := s.Text()
//
//     req.Header.Add("Content-Type", "application/json")
//     req.Header.Add("x-api-key", key)
//   }
//
//
//   res, err := client.Do(req)
//   if err != nil {
//     fmt.Println(err)
//     return
//   }
//   defer res.Body.Close()
//
//   body, err := ioutil.ReadAll(res.Body)
//   if err != nil {
//     fmt.Println(err)
//     return
//   }
//   fmt.Println(string(body))
}

func Index(url string) (*Crawler, error) {
  var wg sync.WaitGroup
  c := NewCrawler()
  c.CleanBody(&wg)

  if (TestResponse(url)) {
    trimmedURL := trimURL(url)
    c.SetLinkPattern(trimmedURL)
    
    c.Visit(url)

    c.Wait()
    wg.Wait()

    fmt.Println("Searched", c.Count, "web addresses")

    return c, nil
  } else {
    return c, errors.New("URL is not valid")
  }
}

