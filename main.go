package main

import (
	"bufio"
	"github.com/Inoi-K/Find-Me/pkg/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()

	line, err := in.ReadString('\n')
	if err != nil {
		log.Printf("couldn't read n: %v", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Printf("coldn't convert n: %v", err)
	}

	users := make([]*user, n)
	for i := 0; i < n; i++ {
		users[i], err = newUser()
		if err != nil {
			log.Printf("couldn't create a new user: %v", err)
		}
	}

	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			log.Printf("Similarity (by Jaccard index) between %v and %v is %v", users[i].name, users[j].name, utils.JaccardIndex(users[i].tags, users[j].tags))
		}
	}
}

type user struct {
	name string
	tags map[string]struct{}
}

func newUser() (*user, error) {
	line, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(line)

	tags, err := newTags()
	if err != nil {
		return nil, err
	}

	return &user{
		name: name,
		tags: tags,
	}, nil
}

func newTags() (map[string]struct{}, error) {
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
