## Wrapper for TOML configuration file extraction

This small module provides the ability to read the configuration file in the toml format, if it not exists create from an empty structure or a specified template and read the values ​​into the given structure

<!--
    TODO:
    Documentation: https://godoc.org/github.com/ntngl/tomfig
-->

### Installation

Using `go get`:
```bash
go get github.com/ntngl/tomfig
```

<!--
    TODO
    [![Build Status](https://travis-ci.org/ntngl/tomfig.svg?branch=v1)](https://travis-ci.org/ntngl/tomfig) [![GoDoc](https://godoc.org/github.com/ntngl/tomfig?status.svg)](https://godoc.org/github.com/ntngl/tomfig)
-->

### Examples

Suppose we have a TOML file like this:

```toml
[section_one]
    sprecial_value = "sample"
    digit = 9

[section_two]
    some_samples = [0, 2, 4, 6, 8]
```

Which could be described as Go struct like this:

```go
// Target config struct
type Config struct {
    One SectionOne `toml:"section_one"`
    Two SectionTwo `toml:"section_two"`
}

type SectionOne struct {
    Value string `toml:"special_value"`
    Digit int
}

type SectionTwo struct {
    Samples []int `toml:"some_samples"`
}
```

And file above can be decoded to presented struct with current module that way:

```go
config := tomfig.NewConfig(pathToFile)
instance := &Config{}

// Throws error when something goes wrong
config.Parse(instance)
```

Then, if all goes right, instance will contain TOML file data.

### Using config template

Package can create new config file, if it's not presented. To do this, you should configure the template (TOML formatted text).

```go
config := tomfig.NewConfig(path)
template := tomlText
config.Template = template
instance := &Config{}

// Also saves template into file by path
config.Parse(instance)
```

if file by path not exists and template not presented, saves file with initial walues fot current `instance`.

### Credits

Thanks [BurntSushi](https://github.com/BurntSushi) for the great TOML [parser](https://github.com/BurntSushi/toml).
