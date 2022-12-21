package persistence

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var ErrConfigFileNotFound = errors.New("unable to find database config file")

type ConnectionDetails struct {
	Provider string
	Database string
	Host     string
	Port     string
	User     string
	Password string
	URL      string
}

func LoadConfigFile() error {
	f, err := os.Open(filepath.Join("./configs", "database.yml"))
	if err != nil {
		return err
	}
	defer f.Close()

	parsedConfig, err := parseConfig(f)
	if err != nil {
		return err
	}

	for name, connDetails := range parsedConfig {
		conn, err := NewConnection(connDetails)
		if err != nil {
			return err
		}

		Connections[name] = conn
	}

	return nil
}

func parseConfig(r io.Reader) (map[string]*ConnectionDetails, error) {
	tmpl := template.New("tmpl")
	tmpl.Funcs(map[string]interface{}{
		"env": func(name string) string {
			return os.Getenv(name)
		},
	})

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	t, err := tmpl.Parse(string(b))
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	if err := t.Execute(&output, nil); err != nil {
		return nil, err
	}

	details := map[string]*ConnectionDetails{}
	if err := yaml.Unmarshal(output.Bytes(), &details); err != nil {
		return nil, err
	}

	return details, nil
}
