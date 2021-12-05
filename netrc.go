package netrc

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type credentials struct {
	Login    string
	Password string
}
type File struct {
	File io.Reader
}

type option func(f *File) error

func NewFile(options ...option) (*File, error) {
	file := File{}
	for _, opt := range options {
		err := opt(&file)
		if err != nil {
			return &File{}, err
		}
	}
	return &file, nil
}

func WithFile(path string) option {
	return func(f *File) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		f.File = file
		return nil
	}
}

func (f *File) Fetch(host string) credentials {
	scanner := bufio.NewScanner(f.File)
	credentials := credentials{}
	machine_found := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "machine ") {
			line_parts := strings.Split(line, " ")
			machine := line_parts[1]
			if machine == host {
				machine_found = true
				if len(line_parts) >= 6 {
					// Netrc entry is using one-line format
					login_found := false
					password_found := false
					for _, part := range line_parts[2:] {
						switch part {
						case "login":
							login_found = true
						case "password":
							password_found = true
						default:
							if login_found && credentials.Login == "" {
								credentials.Login = part
							}
							if password_found && credentials.Password == "" {
								credentials.Password = part
							}
						}
					}
				}
			}
		}
		// Found everything in one-line format (above)
		if credentials.Login != "" && credentials.Password != "" {
			break
		}
		if machine_found && strings.HasPrefix(line, "login ") {
			line_parts := strings.Split(line, " ")
			credentials.Login = line_parts[1]
		}
		if machine_found && strings.HasPrefix(line, "password") {
			line_parts := strings.Split(line, "password ")
			credentials.Password = line_parts[1]
		}
		// Found everything on multiple lines
		if credentials.Login != "" && credentials.Password != "" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		panic("Unable to read file")
	}
	return credentials
}
