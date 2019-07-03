package tomfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config describes object that implements toml config file parser
type Config struct {
	Path       string
	Template   string
	workingDir string
	absPath    string
}

// NewConfig sets up config with filesystem path fallback
func NewConfig(path string) *Config {
	return &Config{Path: path}
}

// Parse is for parsing toml config file by given path, and saving config
// template if file by path not found.
func (cnf *Config) Parse(instance interface{}) error {
	// Setting up default config file name from binary's name if that not setted up before
	if len(cnf.Path) == 0 {
		_, path := filepath.Split(os.Args[0])
		ext := ".toml"
		cnf.Path = path + ext
	}

	cnf.absPath = cnf.Path
	if !filepath.IsAbs(cnf.absPath) {
		// Detecting working directory
		workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return fmt.Errorf("Unable to receive working directory. %s", err)
		}

		// Setting up absolute path
		cnf.workingDir = workingDir
		cnf.absPath = filepath.Join(workingDir, cnf.absPath)
	}

	// Handling a situation when a file does not exist
	_, err := os.Stat(cnf.absPath)
	if os.IsNotExist(err) {
		if len(cnf.Template) != 0 {
			// If template setted up, saving it to file
			if err = saveTemplate(cnf.absPath, cnf.Template); err != nil {
				return fmt.Errorf("Unable to save config template by path %s. %s", cnf.absPath, err)
			}
		} else {
			// Saving empty file initialized by given struct default values
			if err = saveEmpty(cnf.absPath, instance); err != nil {
				return fmt.Errorf("Unable to initialize empty config file from struct. %s", err)
			}
		}
	}

	// Decoding file got by absolute path to given config instance with custom struct
	if _, err := toml.DecodeFile(cnf.absPath, instance); err != nil {
		return fmt.Errorf("Unable to parse config file by path '%s'. %s", cnf.absPath, err)
	}
	return nil
}

// saving config template to target path
func saveTemplate(path, template string) error {
	configFile, err := newConfigFile(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	if _, err = configFile.Write([]byte(template)); err != nil {
		return err
	}
	return nil
}

func saveEmpty(path string, empty interface{}) error {
	configFile, err := newConfigFile(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := toml.NewEncoder(configFile)
	return encoder.Encode(empty)
}

func newConfigFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	configFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return configFile, nil
}
