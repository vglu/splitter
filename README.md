# splitter
split text file for chunks

Task: Split a big text file for a smallest parts for import, for example

Version: 1.1.2
Build Time: 2025-03-07T00:42:23Z
  -help
        Show help
  -input_file string
        Input file to split
  -lines_per_file int
        Number of lines per file
  -output_dir string
        Output directory by default same as input file
  -use_header
        Use header in each file default is true (default true)
Example:

`splitter -input_file=source.csv -lines_per_file=1000 -use_header=true`
`splitter -input_file=POS_FILE227_2023-06-02_17-54-34.TXT -lines_per_file=50000 -use_header=false`

