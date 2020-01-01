package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	// voiceDict := `http://inter.youdao.com/dictvoice?audio=The+congregation+had+11+publishers%2C+and+Brother+Mills+was+company+%28congregation%29+servant.&type=1&le=en`
	pathfile := "word.txt"

	file, err := os.Open(pathfile)
	ErrCheck(err)
	defer file.Close()

	fileReader := bufio.NewReader(file)

	// Extract all word from the file
	for {
		lineWord, err := fileReader.ReadString('\n')
		if err == io.EOF {
			log.Print(lineWord)
			break
		}
		ErrCheck(err)

		var (
			query = strings.Trim(lineWord, "\n")
			from  = "en"
			to    = "id"
		)

		// API for dictionary
		dictAPI := fmt.Sprintf("http://inter.youdao.com/intersearch?q=%s&from=%s&to=%s", query, from, to)

		// Getting content of query with GET http
		resp, err := http.Get(dictAPI)
		ErrCheck(err)
		defer resp.Body.Close()
		log.Print("Getting query content.\n")

		// Reading content of query in bytes form
		body, err := ioutil.ReadAll(resp.Body)
		ErrCheck(err)

		// Convert bytes content of query to string
		bodyStr := string(body)

		// Create directorh
		dirDict := CreateDir("enid-dicts")
		// log.Printf("Create `%s` directory.", dirDict)

		// Create json file
		JSONFile, err := CreateJSON(query, dirDict)
		ErrCheck(err)
		defer JSONFile.Close()
		// log.Printf("Create %s.json file.", query)

		// Writting query content to JSON file
		writeStr, err := JSONFile.WriteString(bodyStr)
		ErrCheck(err)
		log.Printf("Wrote %d bytes\n", writeStr)

		// Commit content of query
		JSONFile.Sync()

		log.Printf("File %s.json created.", query)
	}

}

// CreateDir make a directory with return directory name
func CreateDir(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err == nil {
			log.Printf("%s directory created.", dir)
		}
		ErrCheck(err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		log.Printf("%s directory does exist.", dir)
	}

	return dir
}

// CreateJSON received query as file name and directory name. Return *osFile and error
func CreateJSON(query string, dir string) (*os.File, error) {
	fileQueryName := fmt.Sprintf("%s.json", query)
	pathFile := path.Join(dir, fileQueryName)

	if _, err := os.Stat(pathFile); !os.IsNotExist(err) {
		log.Printf("%s.json file does exist.", query)
	}

	file, err := os.Create(pathFile)
	log.Printf("%s.json file created.", file.Name())

	return file, err
}

// ErrCheck return error from error recived.
func ErrCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
