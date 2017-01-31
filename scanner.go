package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "strings"
import "net/http"


var client = &http.Client{}


func worker(jobs <-chan string, results chan <- string) {
	for url := range jobs {
		resp, err := client.Get(url)
		if err != nil {
			results <- fmt.Sprintf("Got an error: %s\n", err)
			continue
		}

		if resp.StatusCode == 200 {
			results <- fmt.Sprintf("/%s: Found\n", url)
		} else {
			results <- fmt.Sprintf("&%s: Not found\n", url)
		}
		resp.Body.Close()
	}
}


func main() {
	args := os.Args
	path := args[1]
	base_url := args[2]
	
	file, err := os.Open(path)
	
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	jobs := make(chan string)
	results := make(chan string)

	for w := 1; w <= 50; w++ {
		go worker(jobs, results)
	}
	
	for {
		str, err := reader.ReadString(10)
		trimmed_str := strings.TrimSuffix(str, "\n")
		
		if (err == io.EOF) {
			break
		} else if err != nil {
			fmt.Printf("Error: %q\n", err)
		}
		
		full_url := fmt.Sprintf("%s/%s", base_url, trimmed_str)
		jobs <- full_url
	}

	close(jobs)

	for range results {
		<- results
	}
}
