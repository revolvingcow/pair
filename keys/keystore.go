package keys

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"
)

const (
	github         = "https://github.com/{username}.keys"
	gitlab         = "https://gitlab.com/{username}.keys"
	token          = "{username}"
	authorizedFile = "~/.ssh/authorized_keys"
)

// Keystore represents a concurrency safe collection of
// public keys associated with usernames.
type Keystore struct {
	sync.Mutex
	Users []User
}

// NewKeystore provides a standard way to instiantiate a
// Keystore.
func NewKeystore() *Keystore {
	store := &Keystore{}
	store.Read()
	return store
}

// Get queries external sources for the public keys
// associated with the specified usernames.
func (store *Keystore) Get(command string, usernames ...string) {
	urls := []string{github, gitlab}

	// Query each service concurrently
	c := make(chan User)
	for _, url := range urls {
		go func(u string) {
			for _, nick := range usernames {
				c <- getKeys(nick, u, command)
			}
		}(url)
	}

	// Make sure everything has completed before quitting
	seconds := time.Duration(5 * len(usernames))
	hits := len(urls) * len(usernames)
	for i := 0; i < hits; i++ {
		select {
		case user := <-c:
			store.Add(user)
		case <-time.After(time.Second * seconds):
		}
	}
}

// Add keys associated with a username safely.
func (store *Keystore) Add(user User) {
	updated := false

	store.Lock()
	if store.Users == nil {
		store.Users = []User{}
	}
	store.Unlock()

	for i, u := range store.Users {
		if u.Name == user.Name {
			for _, k := range user.Keys {
				if !contains(u.Keys, k) {
					store.Lock()
					store.Users[i].Keys = append(store.Users[i].Keys, k)
					store.Unlock()
				}
			}

			store.Lock()
			store.Users[i].Commands = user.Commands
			store.Unlock()

			updated = true
		}
	}

	if !updated {
		store.Lock()
		store.Users = append(store.Users, user)
		store.Unlock()
	}
}

// Read the ```authorized_keys``` file contents to
// identify managed settings.
func (store *Keystore) Read(usernames ...string) {
	path := expand(authorizedFile)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	filter := len(usernames)
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()

		// If the line was not generating using this
		// program then leave alone
		commentToken := "# pair: "
		comment := bytes.Index(line, []byte(commentToken))
		if comment == -1 {
			continue
		}

		ssh := bytes.Index(line, []byte("ssh-"))
		if ssh == -1 {
			continue
		}

		commands := strings.Split(string(bytes.TrimSpace(line[:ssh])), "\n")
		key := string(bytes.TrimSpace(line[ssh:comment]))
		nick := string(bytes.TrimSpace(line[comment+len(commentToken):]))

		if filter == 0 || contains(usernames, nick) {
			store.Add(User{
				Name:     nick,
				Keys:     []string{key},
				Commands: commands,
			})
		}
	}
}

// Remove the usernames from the key collection. If no usernames
// were specified all associated users will be removed.
func (store *Keystore) Remove(usernames ...string) {
	if len(usernames) == 0 {
		store.Users = []User{}
		return
	}

	keep := []User{}
	for _, nick := range usernames {
		for i, user := range store.Users {
			if user.Name != nick {
				keep = append(keep, store.Users[i])
			}
		}
	}

	store.Users = keep
}

// Save writes the key store to the ```authorized_keys``` file
// while attempting to preserve non-generated settings.
func (store *Keystore) Save() {
	file := expand(authorizedFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	data = removeKeys(store, data)
	data = writeKeys(store, data)

	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Sync will requery the given users (or all users if nothing was specified)
// and refresh there public key stores.
func (store *Keystore) Sync(command string, usernames ...string) {
	if len(usernames) == 0 {
		usernames := []string{}
		for _, user := range store.Users {
			usernames = append(usernames, user.Name)
		}
	}

	for i, user := range store.Users {
		if contains(usernames, user.Name) {
			store.Users[i].Keys = []string{}
		}
	}

	store.Get(command, usernames...)
}

// String outputs a pretty-ish display of the public keys affected by the action(s) taken.
func (store *Keystore) String(prefix string, usernames ...string) string {
	digest := ""
	filter := len(usernames)

	for _, user := range store.Users {
		if filter == 0 || contains(usernames, user.Name) {
			for _, k := range user.Keys {
				digest = fmt.Sprintf("%s%s[ %32s ] %s\n", digest, prefix, k[40:72], user.Name)
			}
		}
	}

	if len(digest) == 0 {
		digest = fmt.Sprintln("No keys were affected")
	}

	return digest
}

func contains(arr []string, text string) bool {
	for _, s := range arr {
		if s == text {
			return true
		}
	}
	return false
}

func containsBytes(arr []byte, text []byte) bool {
	return bytes.Index(arr, text) != -1
}

func expand(path string) string {
	u, err := user.Current()
	if err != nil {
		return path
	}
	return strings.Replace(path, "~", u.HomeDir, -1)
}

func getKeys(username, url, command string) User {
	user := User{
		Name: username,
		Keys: []string{},
		Commands: []string{
			"no-port-forwarding",
			"no-X11-forwarding",
			"no-agent-forwarding",
		},
	}

	if command != "" {
		user.Commands = append([]string{fmt.Sprintf(`command="%s"`, command)}, user.Commands...)
	}

	if len(username) == 0 || len(url) == 0 {
		return user
	}

	response, err := http.Get(strings.Replace(url, token, username, -1))
	if err != nil {
		return user
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil || bytes.Contains(contents, []byte("DOCTYPE")) {
		return user
	}

	user.Keys = strings.Split(strings.TrimSpace(string(contents)), "\n")
	return user
}

func removeKeys(store *Keystore, data []byte) []byte {
	out := []byte{}
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()
		skip := false

		pair := bytes.Index(line, []byte("# pair: "))
		if pair != -1 {
			skip = true
		}

		if !skip {
			out = append(out, line...)
			out = append(out, []byte("\n")...)
		}
	}

	return out
}

func writeKeys(store *Keystore, data []byte) []byte {
	for _, user := range store.Users {
		security := strings.Join(user.Commands, ",")
		for _, k := range user.Keys {
			line := []byte(fmt.Sprintf("%s %s # pair: %s\n", security, k, user.Name))
			if !containsBytes(data, []byte(k)) {
				data = append(data, line...)
			}
		}
	}

	return data
}
