package prettyflags

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
)

type CommandLineFlag struct {
	Name     string
	Value    interface{}
	Usage    string
	AltNames *[]string
}

var (
	yellow    = color.New(color.FgHiYellow).FprintfFunc()
	red       = color.New(color.FgHiRed).FprintfFunc()
	green     = color.New(color.FgHiGreen).FprintfFunc()
	cyan      = color.New(color.FgHiCyan).FprintfFunc()
	white     = color.New(color.FgHiWhite).FprintfFunc()
	darkWhite = color.New(color.FgWhite).FprintfFunc()
	blue      = color.New(color.FgHiBlue).FprintfFunc()
	magenta   = color.New(color.FgHiMagenta).FprintfFunc()

	y  = color.New(color.FgHiYellow).SprintFunc()
	r  = color.New(color.FgHiRed).SprintFunc()
	g  = color.New(color.FgHiGreen).SprintFunc()
	c  = color.New(color.FgHiCyan).SprintFunc()
	w  = color.New(color.FgHiWhite).SprintFunc()
	dw = color.New(color.FgWhite).SprintFunc()
	b  = color.New(color.FgHiBlue).SprintFunc()
	m  = color.New(color.FgHiMagenta).SprintFunc()

	level1 = "    "
	level2 = "      "
	level3 = "        "

	//TODO: Figure out how to allow colors to be used as parameter inputs. Currently the length of the ANSI control
	//      sequence is added to the length of the string and messes up the padding since the characters are counted
	//      but aren't printed to the screen.
	parameterFormat = "%s%-20s %-6s %-20v %s\n"
)

// NewFlagHandler creates a new flag handler for the given flag set.
func NewFlagHandler(app string, version, branch, commit, tag *string) *flagHandler {
	na := "N/A"
	if version == nil {
		version = &na
	}
	if branch == nil {
		branch = &na
	}
	if commit == nil {
		commit = &na
	}
	if tag == nil {
		tag = &na
	}
	return &flagHandler{
		flagKeys: make([]string, 0),
		flagMap:  make(map[string][]*CommandLineFlag),
		app:      app,
		version:  version,
		branch:   branch,
		commit:   commit,
		tag:      tag,
	}
}

// flagHandler is a struct that holds the flag map and flag keys.
type flagHandler struct {
	flagMap  map[string][]*CommandLineFlag
	flagKeys []string
	app      string
	version  *string
	branch   *string
	commit   *string
	tag      *string
}

// PrintAppHeader prints the application header
func (f *flagHandler) PrintAppHeader(output io.Writer) {
	// TODO: Make more generic and less specific to how I use this package
	white(output, "[+] %s\n", f.app)
	white(output, "%sVersion: %s [ Branch: %s | Commit: %s | Tag: %s ]\n", level1, y(f.version), dw(f.branch), dw(f.commit), dw(f.tag))
}

// PrintUsageSectionHeader prints a section header
func (f *flagHandler) PrintUsageSectionHeader(output io.Writer, section string) {
	blue(output, "\n%s%s:\n", level1, section)
	white(output, parameterFormat, level1, "Parameter", "Short", "Default", "Description")
}

// PrintUsageLine prints the usage information for the given parameter
func (f *flagHandler) PrintUsageLine(output io.Writer, parameter string, short string, defaultValue interface{}, description string) {
	switch defaultValue.(type) {
	case string:
		if defaultValue.(string) == "" {
			defaultValue = "\"\""
		}
	}
	darkWhite(output, parameterFormat, level1, parameter, short, defaultValue, description)
}

// Usage automatically generates the usage text in a more readable format than the default flag.Usage() the
// flags package provides.
func (f *flagHandler) Usage(output io.Writer) (usage func()) {
	return func() {
		f.PrintAppHeader(output)
		for _, key := range f.flagKeys {
			f.PrintUsageSectionHeader(output, fmt.Sprintf("%s Options", key))
			for _, flg := range f.flagMap[key] {
				param := fmt.Sprintf("--%s", flg.Name)
				short := ""
				if flg.AltNames != nil {
					short = fmt.Sprintf("-%s", (*flg.AltNames)[0])
				}
				f.PrintUsageLine(output, param, short, flg.Value, flg.Usage)
			}
		}
		fmt.Println("")
	}
}

func (f *flagHandler) Parse() {
	flag.Usage = f.Usage(os.Stdout)
	flag.Parse()
}

func (f *flagHandler) AddFlagBool(name string, section interface{}, defaultValue bool, usage string, altNames *[]string) *bool {
	var v bool

	flag.BoolVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.BoolVar(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagString(name string, section interface{}, defaultValue string, usage string, altNames *[]string) *string {
	var v string

	flag.StringVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.StringVar(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagInt(name string, section interface{}, defaultValue int, usage string, altNames *[]string) *int {
	var v int

	flag.IntVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.IntVar(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagInt64(name string, section interface{}, defaultValue int64, usage string, altNames *[]string) *int64 {
	var v int64

	flag.Int64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Int64Var(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagUnint(name string, section interface{}, defaultValue uint, usage string, altNames *[]string) *uint {
	var v uint

	flag.UintVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.UintVar(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagUint64(name string, section interface{}, defaultValue uint64, usage string, altNames *[]string) *uint64 {
	var v uint64

	flag.Uint64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Uint64Var(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) AddFlagFloat64(name string, section interface{}, defaultValue float64, usage string, altNames *[]string) *float64 {
	var v float64

	flag.Float64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Float64Var(&v, altName, defaultValue, usage)
		}
	}

	f.addToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func (f *flagHandler) addToFlagMap(name string, section interface{}, defaultValue interface{}, usage string, altNames *[]string) {

	var sections []string
	switch section.(type) {
	case string:
		sections = []string{section.(string)}
	case []string:
		sections = section.([]string)
	default:
		panic("Invalid section type")
	}

	for _, section := range sections {
		if _, ok := f.flagMap[section]; !ok {
			f.flagMap[section] = make([]*CommandLineFlag, 0)
			f.flagKeys = append(f.flagKeys, section)
		}
		clf := CommandLineFlag{
			Name:     name,
			Value:    defaultValue,
			Usage:    usage,
			AltNames: altNames,
		}
		f.flagMap[section] = append(f.flagMap[section], &clf)
	}
}
