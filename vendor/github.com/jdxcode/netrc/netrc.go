package netrc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ErrInvalidNetrc means there was an error parsing the netrc file
var ErrInvalidNetrc = errors.New("Invalid netrc")

// Netrc file
type Netrc struct {
	Path     string
	machines []*Machine
	tokens   []string
}

// Machine from the netrc file
type Machine struct {
	Name      string
	IsDefault bool
	tokens    []string
}

// New creates and returns a new empty Netrc for writing to.
func New(p string) *Netrc {
	return &Netrc{machines: make([]*Machine, 0, 20), Path: p}
}

// Parse the netrc file at the given path
// It returns a Netrc instance
func Parse(path string) (*Netrc, error) {
	file, err := read(path)
	if err != nil {
		return nil, err
	}
	netrc, err := parse(lex(file))
	if err != nil {
		return nil, err
	}
	netrc.Path = path
	return netrc, nil
}

// ParseString behaves just like Parse, but can be used to parse netrc data
// without loading it from a file.
func ParseString(contents string) (*Netrc, error) {
	r := strings.NewReader(contents)
	netrc, err := parse(lex(r))
	if err != nil {
		return nil, err
	}
	return netrc, nil
}

// Machine gets a machine by name
func (n *Netrc) Machine(name string) *Machine {
	for _, m := range n.machines {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// AddMachine adds a machine
func (n *Netrc) AddMachine(name, login, password string) {
	machine := n.Machine(name)
	if machine == nil {
		machine = &Machine{}
		n.machines = append(n.machines, machine)
	}
	machine.Name = name
	machine.tokens = []string{"machine ", name, "\n"}
	machine.Set("login", login)
	machine.Set("password", password)
}

// RemoveMachine remove a machine
func (n *Netrc) RemoveMachine(name string) {
	for i, machine := range n.machines {
		if machine.Name == name {
			n.machines = append(n.machines[:i], n.machines[i+1:]...)
			// continue removing but start over since the indexes changed
			n.RemoveMachine(name)
			return
		}
	}
}

// Render out the netrc file to a string
func (n *Netrc) Render() string {
	var b bytes.Buffer
	for _, token := range n.tokens {
		b.WriteString(token)
	}
	for _, machine := range n.machines {
		for _, token := range machine.tokens {
			b.WriteString(token)
		}
	}
	return b.String()
}

// Save the file to disk
func (n *Netrc) Save() error {
	body := []byte(n.Render())
	if filepath.Ext(n.Path) == ".gpg" {
		cmd := exec.Command("gpg", "-a", "--batch", "--default-recipient-self", "-e")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		stdin.Write(body)
		stdin.Close()
		cmd.Stderr = os.Stderr
		body, err = cmd.Output()
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(n.Path, body, 0600)
}

func read(path string) (io.Reader, error) {
	if filepath.Ext(path) == ".gpg" {
		cmd := exec.Command("gpg", "--batch", "--quiet", "--decrypt", path)
		cmd.Stderr = os.Stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
		return stdout, nil
	}
	return os.Open(path)
}

func lex(file io.Reader) []string {
	commentRe := regexp.MustCompile("(^#|\\s+#)")
	scanner := bufio.NewScanner(file)
	scanner.Split(func(data []byte, eof bool) (int, []byte, error) {
		var loc []int
		if eof && len(data) == 0 {
			return 0, nil, nil
		}
		inWhitespace := unicode.IsSpace(rune(data[0]))
		for i, c := range data {
			if c == '#' {
				// line might have a comment
				// but if our regexp returns nil, keep going
				loc = commentRe.FindIndex(data)
				if loc != nil && loc[0] == 0 {
					// currently in a comment
					i = bytes.IndexByte(data, '\n')
					if i == -1 {
						// no newline at end
						if !eof {
							return 0, nil, nil
						}
						i = len(data)
					}
					for i < len(data) {
						if !unicode.IsSpace(rune(data[i])) {
							break
						}
						i++
					}
					return i, data[0:i], nil
				}
			}
			if unicode.IsSpace(rune(c)) != inWhitespace {
				return i, data[0:i], nil
			}
		}
		if eof {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	tokens := make([]string, 0, 100)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

func parse(tokens []string) (*Netrc, error) {
	n := &Netrc{}
	n.machines = make([]*Machine, 0, 20)
	var machine *Machine
	for i, token := range tokens {
		// group tokens into machines
		if token == "machine" || token == "default" {
			// start new group
			machine = &Machine{}
			n.machines = append(n.machines, machine)
			if token == "default" {
				machine.IsDefault = true
				machine.Name = "default"
			} else {
				machine.Name = tokens[i+2]
			}
		}
		if machine == nil {
			n.tokens = append(n.tokens, token)
		} else {
			machine.tokens = append(machine.tokens, token)
		}
	}
	return n, nil
}

// Get a property from a machine
func (m *Machine) Get(name string) string {
	i := 4
	if m.IsDefault {
		i = 2
	}
	for {
		if i+2 >= len(m.tokens) {
			return ""
		}
		if m.tokens[i] == name {
			return m.tokens[i+2]
		}
		i = i + 4
	}
}

// Set a property on the machine
func (m *Machine) Set(name, value string) {
	i := 4
	if m.IsDefault {
		i = 2
	}
	for i+2 < len(m.tokens) {
		if m.tokens[i] == name {
			m.tokens[i+2] = value
			return
		}
		i = i + 4
	}
	m.tokens = append(m.tokens, "  ", name, " ", value, "\n")
}
