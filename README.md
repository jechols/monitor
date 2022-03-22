# Offsite Monitor: Get Website Timeouts and Failures

OM: GWTF is a tool to do offsite monitoring of some of our sites on the cheap
because free services tend to be iffy at best (a simple "ping" sort of process
a few times a day), and even paid services don't look that great.

## Build

Go must be installed to build this, but the binary can be copied just about
anywhere post-build, and can be run against a remote URL directly, so there is
no need to have the Go toolchain in production.  Process:

```bash
cd om-gwtf       # if you're not already here
make             # or just run "go build" - I'm just really lazy
./search-test
```

## Tools

Once built, you'll have the multi-tool binary, `om-gwtf`. This has several
sub-tools built in:

- `om-gwtf -tool oregonnews` hits the oregonnews site with random search terms and
  reports success, time to retrieve the page (in milliseconds), and other data
  for determining success.
- `om-gwtf -tool libweb` hits the library homepage and simply verifies that it has
  basic text in the response, and report success or failure.

Each tool allows a hostname parameter so we aren't stuck recompiling if the
hostnames change, or we want to run a tool against a staging server or
something. Example: `om-gwtf -host oregonnews.uoregon.edu -tool oregonnews`. If
a host is not specified, localhost is assumed. This is not likely what you want
to use.

### Oregon News search test

The search test for Oregon News runs randomized queries (to break the cache)
against our ONI website, verifying that the site is not only up and running,
but properly responding to searches, as those tend to invoke the most complex
pieces of the ONI stack.

The exit status will be zero on success (defined below) and non-zero otherwise.
For a simple "it works" test, that's all that's needed.

Success is defined as hitting the search URL with random words and getting a
valid response body back which contains the text "\d\d\d+ results" somewhere.
We ensure three digits minimum because with any "or" query against five common
English words, we're guaranteed to have a huge result set.  If it's less than
100, something went horribly wrong.  Chances are if it's less than 10k,
something is wrong, but we don't want the cutoff to be too high.

STDERR gets a bit of informational logging, and STDOUT gets some more useful
JSON logging which may be handy for more advanced monitoring someday, and could
be captured to a file in the shorter term.

Sample JSON output (formatted for readability):

```json
{
  "Start": "2020-12-02T04:54:15.7078599-08:00",
  "Error": "",
  "Success": true,
  "URL": "https://oregonnews.uoregon.edu/search/pages/results?andtext=&city=&county=&date1=1846-01-01&date2=2020-12-31&frequency=&language=&ortext=squirrel+doll+table+back+corn&phrasetext=&proxdistance=5&proxtext=&rows=20&searchType=advanced",
  "Words": "squirrel doll table back corn",
  "ResultCount": 984681,
  "DurationMillis": 1442
}
```
