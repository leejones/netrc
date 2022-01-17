package netrc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Credentials struct {
	Username string
	Password string
}

type File struct {
	File io.Reader
}

type option func(f *File) error

func DefaultNetrcPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to find user's home directory")
	}
	return path.Join(homeDir, ".netrc")
}

func Get(machine string) (Credentials, error) {
	netrcFile, err := NewFile(
		WithFile(DefaultNetrcPath()),
	)
	if err != nil {
		return Credentials{}, err
	}
	creds, err := netrcFile.Get(machine)
	if err != nil {
		return Credentials{}, err
	}
	return creds, nil
}

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

func (f *File) Get(host string) (Credentials, error) {
	credentials := Credentials{}
	credentialsFound := false
	machineFound := false
	scanner := bufio.NewScanner(f.File)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "machine ") {
			lineParts := strings.Split(line, " ")
			machine := lineParts[1]
			if machine == host {
				machineFound = true
				// Check for the single-line format, for example:
				// 		machine foo login bar password baz
				if len(lineParts) >= 6 {
					loginFound := false
					passwordFound := false
					for _, part := range lineParts[2:] {
						switch part {
						case "login":
							loginFound = true
						case "password":
							passwordFound = true
						default:
							if loginFound && credentials.Username == "" {
								credentials.Username = part
							}
							if passwordFound && credentials.Password == "" {
								credentials.Password = part
							}
						}
					}
				}
			}
		}
		// Found everything in single-line format (above)
		if credentials.Username != "" && credentials.Password != "" {
			credentialsFound = true
			break
		}
		if machineFound && strings.HasPrefix(line, "login ") {
			lineParts := strings.Split(line, " ")
			credentials.Username = lineParts[1]
		}
		if machineFound && strings.HasPrefix(line, "password") {
			lineParts := strings.Split(line, "password ")
			credentials.Password = lineParts[1]
		}
		// Found everything on multiple lines
		if credentials.Username != "" && credentials.Password != "" {
			credentialsFound = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return Credentials{}, fmt.Errorf("Error reading file: %v", err)
	}

	if !credentialsFound {
		return Credentials{}, fmt.Errorf("Credentials for machine: %v not found", host)
	}
	return credentials, nil
}
