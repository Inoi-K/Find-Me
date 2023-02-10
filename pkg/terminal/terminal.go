package terminal

import (
	"bufio"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/user"
	"log"
	"os"
	"strconv"
	"strings"
)

func ProcessInput() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
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

	users := make([]*user.User, usersCount)
	for i := 0; i < usersCount; i++ {
		users[i], err = readUser(in, i, spheres)
		if err != nil {
			log.Printf("couldn't create a new user: %v", err)
		}
	}

	user.ShowSimilarityAll(users, mainSphere)
}

// readUser gets the required information about the user from a terminal
func readUser(in *bufio.Reader, id int, spheres []string) (*user.User, error) {
	fmt.Printf("Name of the User %v: ", id)
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
		tags, err := readTags(in)
		if err != nil {
			return nil, err
		}
		sphereTags[sphere] = tags
	}

	return user.NewUser(name, sphereDescription, sphereTags)
}

// readTags gets the user's tags from a terminal
func readTags(in *bufio.Reader) (map[string]struct{}, error) {
	line, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return user.NewTags(line)
}
