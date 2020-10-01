package termcolor

import (
	"go.uber.org/zap"
	"math"
	"strconv"
	"strings"
)

var (
	// ANSI terminal control codes

	// ColorPrefix contains the ANSI control code prefix
	ColorPrefix = "\u001b["
	// ColorSuffix contains the ANSI control code suffix
	ColorSuffix = "m"

	// Reset contains the ANSI control code to reset the format
	Reset = ColorPrefix + "0" + ColorSuffix

	// Bold contains the ANSI control code to enable bold text
	Bold = "1" // Bold effect applies to text only (not background)

	// Underline contains the ANSI control code to enable underlined text
	Underline = "4" // Underline text

	// Reverse contains the ANSI control code to enable reverse text
	Reverse = "7" // Use reverse text (swap foreground and background)
)

// Fg returns ANSI color code for color "hex", applies to foreground
func Fg(hex string) string {
	return "38;5;" + strconv.Itoa(int(HexToColor256(hex)))
}

// Bg return ANSI color code for color "hex", applies to background
func Bg(hex string) string {
	return "48;5;" + strconv.Itoa(int(HexToColor256(hex)))
}

// HexToColor256 converts input "hex" into an ANSI color code (xterm-256)
// hex may be a set of 3 or 6 hexadecimal digits for RGB values
func HexToColor256(hex string) uint8 {
	return RgbToColor256(HexToRgb(hex))
}

// RgbToColor256 converts the red, green, blue tuple into an ANSI color
// code (xterm-256)
func RgbToColor256(red, green, blue uint8) uint8 {
	// red, green, blue range 0-255 on input

	if red == green && red == blue {
		// Grayscale
		switch {
		case red == 255:
			// Bright white
			return 15
		case red == 0:
			// Black
			return 0
		}

		return 232 + uint8(math.Round(float64(red)/10.65))
	}

	return (36*ColorToAnsiIndex(red) +
		6*ColorToAnsiIndex(green) +
		ColorToAnsiIndex(blue)) + 16
}

// ColorToAnsiIndex converts a uint8 color value (0-255) into an ANSI
// color value (0-5)
func ColorToAnsiIndex(c uint8) uint8 {
	return uint8(math.Round(float64(c) / 51.0))
}

// HexToRgb converts a 3 or 6 digit hexadecimal string, representing RGB
// values, into separate R,G,B values
func HexToRgb(hex string) (r, g, b uint8) {
	switch len(hex) {
	case 6:
		r = HexToUint8(hex[:2])
		g = HexToUint8(hex[2:4])
		b = HexToUint8(hex[4:])
	case 3:
		r = HexToUint8(hex[:1] + hex[:1])
		g = HexToUint8(hex[1:2] + hex[1:2])
		b = HexToUint8(hex[2:3] + hex[2:3])
	default:
		zap.S().Errorf("Unsupported color %s", hex)
		return uint8(255), uint8(255), uint8(255)
	}
	return
}

// HexToUint8 converts a hexadecimal string to a uint8
func HexToUint8(hexbyte string) uint8 {
	val, _ := strconv.ParseUint(hexbyte, 16, 8)
	return uint8(val)
}

// ExpandStyle expands a style string into an ANSI control code
func ExpandStyle(style string) string {
	switch {
	case strings.HasPrefix(style, "Fg#"):
		fgstyle := style[3:]
		return Fg(fgstyle)
	case strings.HasPrefix(style, "Bg#"):
		bgstyle := style[3:]
		return Bg(bgstyle)
	case style == "Bold":
		return Bold
	case style == "Underline":
		return Underline
	case style == "Reverse":
		return Reverse
	default:
		zap.S().Warnf("Unhandled style [%s]", style)
	}
	return style
}

/* Colorize "message" with "style".
 *
 * Style may contain multiple conditions, delimited by "?".  Such
 * conditional parts are called "clauses".
 *
 * Clauses may contain an "=", in which case the part before the "="
 * is the pattern which "message" must match, and the part after the
 * "=" is the style to use.
 *
 * If no "?" is present, then there is no pattern or clause; all of
 * "message" will be output in the specified "style".
 *
 * The style is formatted like "part[|part[|..]].".
 *
 * Each Part may be a color specification or an effect specification.
 *
 * Color specifications look like "Fg#rgb" or "Bg#rgb".  Fg applies
 * color rgb to the foreground; Gb to the background.  Rgb consists of
 * red, green, and blur components in hexadecimal format.  Rgb may be
 * 3 or 6 digits long, consisting of 1 or 2 digits for r, g, and b.
 * Effects are listed above, such as Bold or Underline.
 *
 * Examples:
 *   Fg#fff changes the foreground color to white
 *   Fg#ffff00|Bold changes the foreground color to bold yellow
 *   ?Y=Fg#0f0|Bold?N=Fg#f00 uses a green foreground if the message
 *       matches "Y", or red if it matches "N"
**/

// Colorize applies style st to message, returning the colorized string
func Colorize(st Style, message string) string {
	var sb strings.Builder

	style := string(st)

	if len(message) == 0 {
		return message
	}

	for _, clause := range strings.Split(style, "?") {
		// ignore whitespace
		clause = strings.TrimSpace(clause)

		// Skip if there is an empty clause
		if len(clause) == 0 {
			continue
		}

		/* If we need to match a specific pattern, skip any patterns
		 * that don't match.
		**/
		if strings.Contains(clause, "=") {
			eq := strings.Index(clause, "=")
			pattern := strings.TrimSpace(clause[:eq])
			style = strings.TrimSpace(clause[eq+1:])

			if pattern != message {
				style = ""
				continue
			}
			break
		}
	}

	if len(style) == 0 {
		return message
	}

	parts := make([]string, 0)

	sb.WriteString(ColorPrefix)
	for _, s := range strings.Split(style, "|") {
		parts = append(parts, ExpandStyle(strings.TrimSpace(s)))
	}
	sb.WriteString(strings.Join(parts, ";"))
	sb.WriteString(ColorSuffix)
	sb.WriteString(message)
	sb.WriteString(Reset)

	return sb.String()
}
