package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()

	tags1, err := readTags()
	if err != nil {
		log.Printf("couldn't read tags: %v", err)
	}
	tags2, err := readTags()
	if err != nil {
		log.Printf("couldn't read tags: %v", err)
	}

	intersection := intersect(tags1, tags2)

	fmt.Printf("Tags 1: %v\nTags 2: %v\nIntersection: %v", tags1, tags2, intersection)
}

func readTags() (map[string]struct{}, error) {
	line, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}

	tags := make(map[string]struct{})
	line = strings.TrimSpace(line)
	for _, tag := range strings.Split(line, " ") {
		tags[tag] = struct{}{}
	}
	delete(tags, "")

	return tags, nil
}

func intersect(s1, s2 map[string]struct{}) map[string]struct{} {
	s3 := make(map[string]struct{})
	for k := range s1 {
		if _, consists := s2[k]; consists {
			s3[k] = struct{}{}
		}
	}
	return s3
}
