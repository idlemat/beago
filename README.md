# beago
A simple command-line application for querying artist album rankings on BestEverAlbums.
## Usage:

  beago [options] [artist name or ID]

Running without options returns the ranking for the first artist found. The flags are:

 -s
>		returns up to ten search results with associated artist IDs
 -c
>		get album results by ID
 -p=1
>		shows the requested album page

## Examples:

* beago -s the fall
* beago -p=2 the fall
* beago -c 1351
