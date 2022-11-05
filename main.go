package main

import (
	"bufio"
	"fmt"
	"github.com/Inoi-K/Find-Me/configs/consts"
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

	fmt.Print("Spheres: ")
	line, err := in.ReadString('\n')
	if err != nil {
		log.Printf("couldn't read spheres: %v", err)
	}
	spheres := strings.Fields(line)

	fmt.Print("Current main sphere: ")
	mainSphere, err := in.ReadString('\n')
	if err != nil {
		log.Printf("couldn't read spheres: %v", err)
	}

	fmt.Print("Number of users: ")
	line, err = in.ReadString('\n')
	if err != nil {
		log.Printf("couldn't read number of users: %v", err)
	}
	usersCount, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Printf("coldn't convert usersCount: %v", err)
	}

	users := make([]*user, usersCount)
	for i := 0; i < usersCount; i++ {
		users[i], err = newUser(i, spheres)
		if err != nil {
			log.Printf("couldn't create a new user: %v", err)
		}
	}

	showSimilarity(users, spheres, mainSphere)
}

type user struct {
	name              string
	sphereDescription map[string]string
	sphereTags        map[string]map[string]struct{}
}

func newUser(id int, spheres []string) (*user, error) {
	fmt.Printf("Name of the user %v: ", id)
	line, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(line)

	sphereDescription := make(map[string]string)
	sphereTags := make(map[string]map[string]struct{})
	for _, sphere := range spheres {
		fmt.Printf("Description for sphere %v: ", sphere)
		desc, err := in.ReadString('\n')
		if err != nil {
			return nil, err
		}
		sphereDescription[sphere] = desc

		fmt.Printf("Unique tags for sphere %v: ", sphere)
		tags, err := newTags()
		if err != nil {
			return nil, err
		}
		sphereTags[sphere] = tags
	}

	usr := &user{
		name:              name,
		sphereDescription: sphereDescription,
		sphereTags:        sphereTags,
	}

	usr.processDescription()

	return usr, nil
}

func (u *user) processDescription() {
	for sphere, desc := range u.sphereDescription {
		words := strings.Fields(desc)
		for _, word := range words {
			word = strings.ToLower(word)
			if tag, exists := consts.Synonyms[word]; exists {
				u.sphereTags[sphere][tag] = struct{}{}
			}
		}
	}
}

func newTags() (map[string]struct{}, error) {
	line, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}

	tags := make(map[string]struct{})
	line = strings.TrimSpace(line)
	for _, tag := range strings.Split(line, " ") {
		tag = strings.ToLower(tag)
		tags[tag] = struct{}{}
	}
	delete(tags, "")

	return tags, nil
}

func showSimilarity(users []*user, spheres []string, mainSphere string) {
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			mainSimilarity := 1.0

			log.Printf("Similarity (by Jaccard index) between %v and %v solely by:\n", users[i].name, users[j].name)
			for _, sphere := range spheres {
				tags1 := users[i].sphereTags[sphere]
				tags2 := users[j].sphereTags[sphere]
				similarity := utils.JaccardIndex(tags1, tags2)
				log.Printf(" - %v: %v", sphere, similarity)

				coefficient := consts.OtherSphereCoefficient
				if sphere == mainSphere {
					coefficient = consts.MainSphereCoefficient
				}
				mainSimilarity += similarity * coefficient
			}

			log.Printf(" MAIN SIMILARITY: %v", mainSimilarity)
		}
	}
}
