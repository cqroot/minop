package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Version info variables
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtWith = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

// Info contains version and build information.
type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltWith string
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuiltWith: builtWith,
	}
}

var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

// String returns a formatted version info string.
func (i Info) String() string {
	sb := strings.Builder{}
	sb.WriteString("\n  ")
	sb.WriteString(labelStyle.Render("• Version:      "))
	sb.WriteString(i.Version)

	sb.WriteString("\n  ")
	sb.WriteString(labelStyle.Render("• Commit:       "))
	sb.WriteString(i.Commit)

	sb.WriteString("\n  ")
	sb.WriteString(labelStyle.Render("• Built at:     "))
	sb.WriteString(i.Date)

	sb.WriteString("\n  ")
	sb.WriteString(labelStyle.Render("• Built with:   "))
	sb.WriteString(i.BuiltWith)
	return sb.String()
}
