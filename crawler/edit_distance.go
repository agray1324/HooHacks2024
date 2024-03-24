package crawler

import (
  // "fmt"
  "regexp"
  "strings"
  "sync"
  "sort"
  "github.com/lithammer/fuzzysearch/fuzzy"
)

// implementation for this algorithm was inspired by this article
// https://www.geeksforgeeks.org/introduction-to-levenshtein-distance/
func LevDistance(s1 *string, s2 *string) int {
  m := len(*s1)
  n := len(*s2)
  
  memo := make([][]int, m + 1)
  for i := range memo {
    memo[i] = make([]int, n + 1)
  }

  for i := range m {
    memo[i][0] = i
  }

  for i := range n {
    memo[0][i] = i
  }

  for i := 1; i <= m; i++ {
    for j := 1; j <= n; j++ {
      if (*s1)[i - 1] == (*s2)[j-1] {
        memo[i][j] = memo[i-1][j-1]
      } else {
        memo[i][j] = 1 + min(memo[i][j-1], min(memo[i-1][j], memo[i-1][j-1]))
      }
    }
  }

  return memo[m][n]
}

func Tokenize(search string) []string {
  // select alphanumeric
  spaces := regexp.MustCompile(`\s+`)
  re := regexp.MustCompile("[^a-zA-Z0-9 -]")
  tokenized := re.ReplaceAllString(search, "")

  // clean white space
  tokenized = spaces.ReplaceAllString(tokenized, " ")
  tokenized = strings.ReplaceAll(tokenized, "\"", "'")
  search = strings.TrimSpace(search)
  tokenized = strings.ToLower(tokenized)
  token_sli := strings.Split(tokenized, " ")

  var revised []string
  for i := 0; i < len(token_sli); i++ {
    if len(token_sli[i]) >= 3 {
      revised = append(revised, token_sli[i])
    }
  }

  return revised
}

func Relevances(search string, page []string) []int {
  relevance := make([]int, len(page))
  search_tokens := Tokenize(search)
  for i := 0; i < len(relevance); i++ {
      for j := 0; j < len(search_tokens); j++ {
        scoreSub := LevDistance(&search_tokens[j], &page[i])
        if scoreSub < 3 {
          relevance[i] += 5
          relevance[i] *= 2
        }

        // if scoreSub > 3 {
        //   scoreSub = max(len(search_tokens[j]), len(page[i]))
        // }
      }
  }
  return relevance
}

func PageScore(search string, page []string) (int, int) {
  r := Relevances(search, page)
  rank := 0
  maxRank := 0


  for i, rel := range r {
    rank += rel
    if rel > r[maxRank] {
      maxRank = i
    }
  }

  return rank, maxRank
}

func PageRankings(search string, pages []string, urls []string) ([]string, [][]string) {
  // var wg sync.WaitGroup
  scores := make([]int, len(pages))
  rel := make([]int, len(pages))
  tokens := make([][]string, len(pages))

  for i := 0; i < len(pages); i++ {
    tokens[i] = Tokenize(pages[i])
  }

  for i := 0; i < len(tokens); i++ {
    scores[i], rel[i] = PageScore(search, tokens[i])
  }

  asyncSortRankDouble(tokens, scores)

  asyncSortRank(urls, scores)

  asyncSortRankInt(rel, scores)

  sort.Slice(scores, func(i, j int) bool {
    return scores[i] > scores[j]
  })

  keywords := make([][]string, 10)
  for i := 0; i < len(keywords); i++ {
    keywords[i] = getKeywords(rel[i], tokens[i])
  }

  return urls[:10], keywords
}


// parallelized sort by rank similarity
func FuzzyRank(search string, content []string) []int {
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

func asyncSortRank(sli []string, ranks []int) {
  sort.Slice(sli, func(i, j int) bool {
    return ranks[i] > ranks[j]
  })
}

func asyncSortRankDouble(sli [][]string, ranks []int) {
  sort.Slice(sli, func(i, j int) bool {
    return ranks[i] > ranks[j]
  })
}

func asyncSortRankInt(sli []int, ranks []int) {
  sort.Slice(sli, func(i, j int) bool {
    return ranks[i] > ranks[j]
  })
}

func getKeywords(rel int, tokens []string) []string{
  return tokens[max(0, rel-5):min(len(tokens), rel+5)]
}
