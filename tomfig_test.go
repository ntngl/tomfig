package tomfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

var tmpl = `# Sample config
[sample]
test = "test"

[another]
digit = 9
`

type testConfig struct {
	Sample  sample
	Another another
}

type sample struct {
	Test string
}

type another struct {
	Digit int
}

type testDict struct {
	Path         string
	AwaitingPath string
	Template     string
	Instance     testConfig
}

func TestNewConfig(t *testing.T) {
	path := "test.toml"
	cnf := NewConfig(path)

	if cnf.Path != path {
		t.Error("How it's posible?!")
	}
}

func TestParse(t *testing.T) {
	workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Error("Woops! ", err)
	}

	_, name := filepath.Split(os.Args[0])
	ext := ".toml"
	path := name + ext

	testTable := []testDict{
		{
			Path:         "",
			AwaitingPath: filepath.Join(workingDir, path),
			Template:     "",
			Instance:     testConfig{Sample: sample{""}, Another: another{0}},
		},
		{
			Path:         "test.toml",
			AwaitingPath: filepath.Join(workingDir, "test.toml"),
			Template:     tmpl,
			Instance:     testConfig{Sample: sample{"test"}, Another: another{9}},
		},
		{
			Path:         "config/test.toml",
			AwaitingPath: filepath.Join(workingDir, "config/test.toml"),
			Template:     "",
			Instance:     testConfig{Sample: sample{""}, Another: another{0}},
		},
	}

	for _, entry := range testTable {
		func(path, wpath, template string, conf testConfig) {
			//defer os.Remove(path)
			cnf := NewConfig(path)
			cnf.Template = template

			instance := &testConfig{}
			if err := cnf.Parse(instance); err != nil {
				t.Error("Wrong! ", err)
				return
			}

			if _, err := os.Stat(wpath); os.IsNotExist(err) {
				t.Errorf("File saved by wrong path, awaits %s, got %s", wpath, cnf.absPath)
				return
			}

			if instance.Sample.Test != conf.Sample.Test || instance.Another.Digit != conf.Another.Digit {
				t.Errorf("Parsed values not match same from sample")
			}
		}(entry.Path, entry.AwaitingPath, entry.Template, entry.Instance)
	}
}

func Test_saveTemplate(t *testing.T) {
	path := "test.toml"
	template := tmpl
	instance := &testConfig{}
	defer os.Remove(path)

	if err := saveTemplate(path, template); err != nil {
		t.Errorf("Unable to save config template by path %s, %s", path, err)
		return
	}

	// Decoding file got by absolute path to given config instance with custom struct
	if _, err := toml.DecodeFile(path, instance); err != nil {
		t.Errorf("Unable to parse config file by path '%s'. %s", path, err)
		return
	}

	one := instance.Sample.Test
	two := instance.Another.Digit
	test1 := "test"
	test2 := 9
	if one != test1 || two != test2 {
		t.Errorf("Something goes wrong. Awaits '%s' and '%d', got '%s' and '%d'", test1, test2, one, two)
	}
}

func Test_saveEmpty(t *testing.T) {
	path := "test.toml"
	instance := &testConfig{}
	defer os.Remove(path)

	if err := saveEmpty(path, instance); err != nil {
		t.Errorf("Unable to create new empty config by path %s, %s", path, err)
		return
	}

	// Decoding file got by absolute path to given config instance with custom struct
	if _, err := toml.DecodeFile(path, instance); err != nil {
		t.Errorf("Unable to parse config file by path '%s'. %s", path, err)
		return
	}

	one := instance.Sample.Test
	two := instance.Another.Digit
	test1 := ""
	test2 := 0
	if one != test1 || two != test2 {
		t.Errorf("Something goes wrong. Awaits '%s' and '%d', got '%s' and '%d'", test1, test2, one, two)
	}
}

func Test_newConfigFile(t *testing.T) {
	path := "test.toml"
	defer os.Remove(path)

	cFile, err := newConfigFile(path)
	if err != nil {
		t.Errorf("Unable to create new config file by path %s, %s", path, err)
		return
	}
	cFile.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("New config file not exists by path %s, %s", path, err)
		return
	}
}
