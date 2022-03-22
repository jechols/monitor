package oregonnews

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"om-gwtf/internal/config"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type statusInfo struct {
	Start          time.Time
	Error          string `json:",omitempty"`
	Success        bool
	URL            string
	ResponseBody   []byte `json:"-"`
	StatusCode     int    `json:",omitempty"`
	Words          string
	ResultCount    int64 `json:",omitempty"`
	DurationMillis int64 `json:",omitempty"`
}

func (s *statusInfo) Print(out io.Writer) {
	var data, err = json.Marshal(s)
	if err != nil {
		panic("error marshaling statusInfo: " + err.Error())
	}
	fmt.Fprintln(out, string(data))
}

func (s *statusInfo) get() {
	s.Start = time.Now()
	var c = &http.Client{Timeout: time.Second * 60}
	var resp, err = c.Get(s.URL)
	if err != nil {
		s.Error = fmt.Sprintf("failed: %s", err)
		return
	}

	s.DurationMillis = time.Since(s.Start).Milliseconds()
	s.StatusCode = resp.StatusCode
	defer resp.Body.Close()
	s.ResponseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Error = fmt.Sprintf("failed: %s", err)
	}
}

// Subset of the Dolch list nouns.  These are some of the most common English
// words.  We choose five of these randomly and perform an "or" search.
var words = []string{
	"apple", "baby", "back", "ball", "bear", "bed", "bell", "bird", "birthday", "boat",
	"box", "boy", "bread", "brother", "cake", "car", "cat", "chair", "chicken", "children",
	"coat", "corn", "cow", "day", "dog", "doll", "door", "duck", "egg", "eye",
	"farm", "farmer", "father", "feet", "fire", "fish", "floor", "flower", "game", "garden",
	"girl", "grass", "ground", "hand", "head", "hill", "home", "horse", "house",
	"kitty", "leg", "letter", "man", "men", "milk", "money", "morning", "mother", "name",
	"nest", "night", "paper", "party", "picture", "pig", "rabbit", "rain", "ring", "robin",
	"school", "seed", "sheep", "shoe", "sister", "snow", "song", "squirrel", "stick", "street",
	"sun", "table", "thing", "time", "top", "toy", "tree", "watch", "water", "way",
	"wind", "window", "wood",
}

// We always expect at least triple-digit results.  With five random, common
// words, it is basicaly impossible to not have a lot of results in production
// or staging.
var searchRE = regexp.MustCompile(`\s+(\d\d\d+) results\s+containing`)

// Run executes the oregonnews search test
func Run(conf *config.Config) bool {
	var info statusInfo

	info.Words = randomWords(5)
	info.URL = makeSearchURL(conf.URL, info.Words)
	info.Error = ""

	info.get()
	if info.Error != "" {
		info.StatusCode = -1
		info.Print(os.Stdout)
		return false
	}

	var matches = searchRE.FindSubmatch(info.ResponseBody)
	if len(matches) != 2 {
		info.StatusCode = -1
		info.Error = "no search results found by regexp"
		info.Print(os.Stdout)
		return false
	}

	info.Success = true
	if conf.PrintBody {
		fmt.Println(string(info.ResponseBody))
	}

	// Ignore the error from Atoi since the regex already forces digit-only input
	info.ResultCount, _ = strconv.ParseInt(string(matches[1]), 10, 64)
	info.Print(os.Stdout)
	return true
}

func randomWords(n int) string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	return strings.Join(words[:n], " ")
}

func makeSearchURL(base *url.URL, query string) string {
	// Yes it's stupid to re-parse the URL just to create a new one, but it's
	// effective and ensures a perfect clone including user/pass if they exist
	var newU, _ = url.Parse(base.String())
	newU.Path = "/search/pages/results"

	var vals = newU.Query()
	vals.Set("ortext", query)
	vals.Set("andtext", "")
	vals.Set("phrasetext", "")
	vals.Set("proxtext", "")
	vals.Set("proxdistance", "5")
	vals.Set("city", "")
	vals.Set("county", "")
	vals.Set("date1", "1846-01-01")
	vals.Set("date2", time.Now().Format("2006-01-02"))
	vals.Set("language", "")
	vals.Set("frequency", "")
	vals.Set("rows", "20")
	vals.Set("searchType", "advanced")

	newU.RawQuery = vals.Encode()
	return newU.String()
}
