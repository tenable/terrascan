package termcolor

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	// ColorPatterns contains the coloring rules
	ColorPatterns map[*regexp.Regexp]FieldStyle
	patternFile   string

	defaultValuePattern = `.*?`
)

// Style contains a style spec for termcolor.Colorize()
type Style string

// FieldStyle contains the styles for a particular output field
type FieldStyle struct {
	KeyStyle   Style
	ValueStyle Style
}

// FieldSpec defines a key/value pattern that the color patterns look for.
type FieldSpec struct {
	KeyPattern   string
	ValuePattern string
}

type colorPatternSerialized struct {
	KeyStyle     string `json:"key-style"`
	ValueStyle   string `json:"value-style"`
	KeyPattern   string `json:"key-pattern"`
	ValuePattern string `json:"value-pattern"`
}

/* -------------------------------------------
 * Patterns which define the output to color.
 *
 * The patterns are applied to the rendered output.  Currently YAML
 * and JSON are supported.
 *
 * The format is roughly:
 * {<key-pattern>, <value-pattern>}: {<key-style>, <value-style>}
 *
 * Where <*-pattern> is a regexp and <*-style> is a style appropriate
 * for Colorize()
**/

var defaultColorPatterns = map[FieldSpec]FieldStyle{
	{"description", defaultValuePattern}:   {"", "Fg#0c0"},
	{"severity", defaultValuePattern}:      {"", "?HIGH=Fg#f00?MEDIUM=Fg#c84?LOW=Fg#cc0"},
	{"resource_name", defaultValuePattern}: {"", "Fg#0ff|Bold"},
	{"resource_type", defaultValuePattern}: {"", "Fg#0cc"},
	{"file", defaultValuePattern}:          {"", "Fg#fff|Bold"},
	{"low", `\d+`}:                         {"Fg#cc0", "Fg#cc0"},
	{"medium", `\d+`}:                      {"Fg#c84", "Fg#c84"},
	{"high", `\d+`}:                        {"Fg#f00", "Fg#f00"},

	{"count", ""}:                      {"Bg#ccc|Fg#000", ""},
	{"rule_name", defaultValuePattern}: {"Bg#ccc|Fg#000", ""},
}

func init() {
	cf := os.Getenv("TERRASCAN_COLORS_FILE")
	if len(cf) > 0 {
		patternFile = cf
	}
}

// GetColorPatterns loads the map used by the colorizer
func GetColorPatterns() map[*regexp.Regexp]FieldStyle {
	var patterns map[FieldSpec]FieldStyle
	var pdata []byte

	if len(ColorPatterns) > 0 {
		return ColorPatterns
	}

	if len(patternFile) > 0 {
		var err error
		pdata, err = ioutil.ReadFile(patternFile)
		if err != nil {
			zap.S().Warnf("Unable to read color patterns: %v", err)
			zap.S().Warn("Will proceed with defaults")
		}
	}

	if len(pdata) > 0 {
		patterns = make(map[FieldSpec]FieldStyle)
		var pd = make([]colorPatternSerialized, 0)

		err := json.Unmarshal(pdata, &pd)
		if err != nil {
			zap.S().Warnf("Unable to process color patterns from %s: %v", patternFile, err)
			zap.S().Warn("Will proceed with defaults")
			patterns = defaultColorPatterns
		}

		for _, item := range pd {
			fsp := FieldSpec{
				KeyPattern:   item.KeyPattern,
				ValuePattern: item.ValuePattern,
			}
			fs := FieldStyle{
				KeyStyle:   Style(item.KeyStyle),
				ValueStyle: Style(item.ValueStyle),
			}

			if len(fsp.ValuePattern) == 0 {
				fsp.ValuePattern = defaultValuePattern
			} else if fsp.ValuePattern == "-" {
				fsp.ValuePattern = ""
			}
			patterns[fsp] = fs
		}
	} else {
		patterns = defaultColorPatterns
	}

	ColorPatterns = make(map[*regexp.Regexp]FieldStyle, len(patterns))

	/* Build the regexp needed for the different patterns */
	for ptn, fmts := range patterns {
		var rePtn string

		/* rePtn should process a whole line and have 5 subgroups */
		if len(ptn.ValuePattern) == 0 {
			rePtn = fmt.Sprintf(`^([-\s]*"?)(%s)("?:\s*?)()(.*?)\s*$`, ptn.KeyPattern)
		} else {
			rePtn = fmt.Sprintf(`^([-\s]*"?)(%s)("?: "?)(%s)("?,?)\s*$`, ptn.KeyPattern, ptn.ValuePattern)
		}
		ColorPatterns[regexp.MustCompile("(?m)"+rePtn)] = fmts
	}

	return ColorPatterns
}
