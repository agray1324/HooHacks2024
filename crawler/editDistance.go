package crawler

import (
  // "fmt"
  "regexp"
  "strings"
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

  return strings.Split(strings.ToLower(tokenized), " ")
}
