package utils

import (
	"github.com/cqroot/minop/pkg/remote"
	"github.com/fatih/color"
	"time"
)

func TimeString() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func FormattedString(fg color.Attribute, emoji string, r *remote.Remote, msg string) string {
	return color.New(fg).Sprintf("[%s] %s [%s@%s] %s", TimeString(), emoji, r.Username, r.Hostname, msg)
}
