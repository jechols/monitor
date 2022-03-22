package libweb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"om-gwtf/internal/config"
	"os"
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

// Run executes the libweb page request test
func Run(conf *config.Config) bool {
	var info statusInfo

	info.URL = makeTestURL(conf.URL)
	info.Error = ""

	info.get()
	if info.Error != "" {
		info.StatusCode = -1
		info.Print(os.Stdout)
		return false
	}

	if !strings.Contains(string(info.ResponseBody), "Knight Library") {
		info.StatusCode = -1
		info.Error = "Expected text not found in the response body"
		info.Print(os.Stdout)
		return false
	}

	info.Success = true
	if conf.PrintBody {
		fmt.Println(string(info.ResponseBody))
	}

	info.Print(os.Stdout)
	return true
}

func makeTestURL(base *url.URL) string {
	// Yes it's stupid to re-parse the URL just to create a new one, but it's
	// effective and ensures a perfect clone including user/pass if they exist
	var newU, _ = url.Parse(base.String())
	newU.Path = "/knight-library"
	return newU.String()
}
