# goscanner
A simple scanner in go to scan urls at high concurrency

# How to use
You should have go locally to be able to use/compile this program.
The program expect 2 args, first is the path to the file containing a list of urls, the
second one is the domain you want to scan

If you want to run it straight use:
`go run scanner.go /path/to/the/file http://something.com`

otherwise compile it with:
`go build -o program_name scanner.go`
and launch it with `./program_name file domain`
