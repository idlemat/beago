// Beago is a simple command-line application for querying artist album rankings on BestEverAlbums.
// Usage:
//
//	  beago [options] [artist name or ID]
//
// Running without options returns the ranking for the first artist found. The flags are:
//
//	  -s
//			returns up to ten search results with associated artist IDs
//	  -c
//			get album results by ID
//	  -p=1
//			shows the requested album page
//
// Examples:
//
// beago -s the fall
// beago -p=2 the fall
// beago -c 1351
//
package beago

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	isSearch := flag.Bool("s", false, "Perform a search with the given argument.")
	useId := flag.Bool("c", false, "Fetch the artist using id.")
	pageNum := flag.Int("p", 1, "Specifies the artist page number.")
	flag.Parse()

	arg := strings.Join(flag.Args(), "+")
	var out string

	if *useId {
		out = processArtistPage(arg, *pageNum)
	} else if *isSearch {
		out = processSearchPage(arg, false)
	} else {
		out = processArtistPage(processSearchPage(arg, true), *pageNum)
	}

	fmt.Println(out)
}

func processSearchPage(query string, lucky bool) string {
	htmlUnprocessed := stringifyPage(query, "search")

	if lucky {
		first := strings.SplitN(htmlUnprocessed, "title=\"Click to see further details regarding this artist\" href=\"/thechart.php?b=", 2)
		return strings.SplitN(first[1], "\"", 2)[0]
	}

	divide := strings.Split(htmlUnprocessed, "title=\"Click to see further details regarding this artist\" href=\"/thechart.php?b=")
	n := len(divide)

	results := make([]string, n, n)
	results[0] = "Most relevant results are:"
	for i, ss := range divide[1:] {
		artistWithId := strings.SplitN(ss, "<", 2)[0]
		results[i+1] = strings.ReplaceAll(artistWithId, "\">", "\t - \t")
	}

	return strings.Join(results, "\n")
}

func processArtistPage(artistID string, pageNum int) string {
	url := fmt.Sprintf("https://www.besteveralbums.com/thechart.php?b=%v&f=&fv=&orderby=InfoRankScore&sortdir=DESC&page=%v", artistID, pageNum)
	htmlUnprocessed := stringifyPage(url, "artist")
	divide := strings.Split(htmlUnprocessed, "title=\"Click to see further details regarding this album.\">")
	n := len(divide)

	topAlbums := make([]string, n, n)
	navInfo := strings.SplitN(divide[n-1], "<span class=\"current\">T", 2)[1]
	topAlbums[0] = "T" + strings.SplitN(navInfo, "<", 2)[0]

	for i, ss := range divide[1:] {
		albumName := strings.SplitN(ss, "<", 2)[0]
		topAlbums[i+1] = fmt.Sprintf("%v. %s", (pageNum-1)*10+i+1, albumName)
	}

	return strings.Join(topAlbums, "\n")
}

func stringifyPage(info string, t string) string {
	var err error
	var r *http.Response
	switch t {
	case "search":
		r, err = http.PostForm("https://www.besteveralbums.com/search.php?r=", url.Values{"txtSearch": {info}, "searchfield": {"band"}, "type": {"all"}})
	case "artist":
		r, err = http.Get(info)
	}
	if err != nil {
		fmt.Println("Unable to get response:")
		fmt.Println(err)
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	return fmt.Sprintf("%s", body)
}
