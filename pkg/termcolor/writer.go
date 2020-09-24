package termcolor

import (
    "io"
    "fmt"
    "regexp"
)

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
var color_patterns = map[struct { key,value string}][2]string {
  {"description",".*?"}: {"", "Fg#0c0"},
  {"severity",".*?"}: {"", "?HIGH=Fg#f00?MEDIUM=Fg#c84?LOW=Fg#cc0"},
  {"resource_name",".*?"}: {"", "Fg#0ff|Bold"},
  {"resource_type",".*?"}: {"", "Fg#0cc"},
  {"file",".*?"}: {"", "Fg#fff|Bold"},
  {"low",".*?"}: {"Fg#cc0", "Fg#cc0"},
  {"medium",".*?"}: {"Fg#c84", "Fg#c84"},
  {"high",".*?"}: {"Fg#f00", "Fg#f00"},

  {"count",""}: {"Bg#ccc|Fg#000", ""},
  {"rule_name",".*?"}: {"Bg#ccc|Fg#000", ""},
}

var (
    patterns map[*regexp.Regexp][2]string
)

type ColorizedWriter struct {
    writer io.Writer
}

func NewColorizedWriter(w io.Writer) ColorizedWriter {
    return ColorizedWriter{w}
}

func init() {
    patterns = make(map[*regexp.Regexp][2]string, len(color_patterns))

    /* Build the regexp needed for the different patterns */
    for ptn,fmts := range color_patterns {
        var rePtn string

        /* rePtn should process a whole line and have 5 subgroups */
        if len(ptn.value) == 0 {
            rePtn = fmt.Sprintf(`^([-\s]*"?)(%s)("?:\s*?)()(.*?)\s*$`, ptn.key)
        } else {
            rePtn = fmt.Sprintf(`^([-\s]*"?)(%s)("?: "?)(%s)("?,?)\s*$`, ptn.key, ptn.value)
        }
        patterns[regexp.MustCompile("(?m)"+rePtn)] = fmts
    }
}

func (me ColorizedWriter) Write(p []byte) (n int, err error) {
    /* Before output is written, perform color substitutions */
    for ptn,style := range patterns {
        p = ptn.ReplaceAllFunc(p, func(matched []byte) []byte {
            groups := ptn.FindStringSubmatch(string(matched))
            return []byte(fmt.Sprintf( "%s%s%s%s%s",
                    groups[1],
                    Colorize(style[0], string(groups[2])),
                    groups[3],
                    Colorize(style[1], string(groups[4])),
                    groups[5]  ))
        });
    }

    return me.writer.Write( p )
}
