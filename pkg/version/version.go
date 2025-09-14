package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtWith = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltWith string
}

func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuiltWith: builtWith,
	}
}

func (i Info) String() string {
	sb := strings.Builder{}
	sb.WriteString("\n  ")
	sb.WriteString(color.HiBlueString("• Version:      "))
	sb.WriteString(i.Version)

	sb.WriteString("\n  ")
	sb.WriteString(color.HiBlueString("• Commit:       "))
	sb.WriteString(i.Commit)

	sb.WriteString("\n  ")
	sb.WriteString(color.HiBlueString("• Built at:     "))
	sb.WriteString(i.Date)

	sb.WriteString("\n  ")
	sb.WriteString(color.HiBlueString("• Built with:   "))
	sb.WriteString(i.BuiltWith)
	return sb.String()
}
