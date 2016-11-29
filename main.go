package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const ghRawURL = "https://raw.githubusercontent.com/%s/%s/master/README.md"

type Config struct {
	Organization string  `json:"organization"`
	Stacks       []Stack `json:"stacks"`
}

type Stack struct {
	Title      string   `json:"title"`
	TagLine    string   `json:"tagLine"`
	Complexity int      `json:"complexity"`
	Price      int      `json:"price"`
	Providers  []string `json:"providers"`
}

func (s *Stack) Write(cfg *Config) error {
	resp, err := http.Get(fmt.Sprintf(ghRawURL, cfg.Organization, s.Title))
	if err != nil {
		return err
	}
	readme, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	metadata, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	if _, err := os.Stat("./content/stack"); err != nil {
		os.Mkdir("./content/stack", 0755)
	}
	path := fmt.Sprintf("./content/stack/%s.md", s.Title)
	fmt.Printf("Writing %s\n", path)
	return ioutil.WriteFile(path, bytes.Join([][]byte{metadata, readme}, []byte("\n")), 0755)
}

func config(path string) *Config {
	raw, err := ioutil.ReadFile(path)
	failOnErr(err)
	cfg := &Config{}
	failOnErr(json.Unmarshal(raw, cfg))
	return cfg
}

func failOnErr(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	cfg := config("./site.json")
	for _, stack := range cfg.Stacks {
		failOnErr(stack.Write(cfg))
	}
}
