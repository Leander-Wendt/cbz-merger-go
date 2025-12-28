# cbz-merger-go
A simple script that merges .cbz files to an output directory.

E.g: the files "001 - 0001 title.cbz", "001 - 0001 title.cbz", "001 - 0001 title.cbz" will be merged into "001_title.cbz" while maintaining the order directed by the second number.  
Usually these are the number of the volume and chapter.


## Usage

- [Install go](https://go.dev/doc/install)
- Move the script into the directory containing the .cbz files
- Run `go run .\merge_cbz.go` in the directory
- Wait until the program finished
- Find the results in the newly created directory "merged"
