package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "strings"
import "sync"
import "net/http"

var client = &http.Client{}

func worker(jobs <-chan string, wg *sync.WaitGroup) {
	for url := range jobs {
		resp, err := client.Get(url)

		if err != nil {
			fmt.Printf("Got an error: %s\n", err)
			wg.Done()
			continue
		}
		if resp.StatusCode == 200 {
			fmt.Printf("%s: Found\n", url)
		}
		resp.Body.Close()
		wg.Done()
	}
}

func main() {
	var waitGroup sync.WaitGroup

	args := os.Args
	path := args[1]
	base_url := args[2]

	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	jobs := make(chan string, 4613)

	for w := 1; w <= 250; w++ {
		go worker(jobs, &waitGroup)
	}

	for {
		str, err := reader.ReadString(10)
		trimmed_str := strings.TrimSuffix(str, "\n")

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error: %q\n", err)
		}
		waitGroup.Add(1)
		full_url := fmt.Sprintf("%s/%s", base_url, trimmed_str)
		jobs <- full_url
	}

	close(jobs)

	waitGroup.Wait()

}
