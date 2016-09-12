package keys

// User represents a Git user and there publicly accessible
// SSH keys.
type User struct {
	Name     string
	Keys     []string
	Commands []string
}
