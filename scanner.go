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


func check_url(full_url string, portion string, wg *sync.WaitGroup) {
	resp, err := client.Get(full_url)

	if err != nil {
		fmt.Printf("Got an error: %s\n", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Printf("/%s: Found\n", portion)
	}
	resp.Body.Close()
	wg.Done()
	return
}

func main() {
	var waitGroup sync.WaitGroup
	
	args := os.Args
	path := args[1]
	base_url := args[2]
	
	file, err := os.Open(path)
	
	defer file.Close()
	
	if err != nil {
		log.Fatal(err)
	}
	
	reader := bufio.NewReader(file)
	
	for {
		str, err := reader.ReadString(10)
		trimmed_str := strings.TrimSuffix(str, "\n")
		
		if (err == io.EOF) {
			break
		} else if err != nil {
			fmt.Printf("Error: %q\n", err)
		}
		
		full_url := fmt.Sprintf("%s/%s", base_url, trimmed_str)
		waitGroup.Add(1)
		go check_url(full_url, trimmed_str, &waitGroup)
	}

	waitGroup.Wait()
}
