package commit

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/mesuutt/git-mirror/pkg/config"
)

const defaultCommitFormat = "%d insertion(s), %d deletion(s)"

type diffGen struct {
	conf *config.Config
}

func NewDiffGenerator(configPath string) diffGen {
	conf, err := config.ReadConfig(configPath)
	if err != nil {
		return diffGen{conf: nil}
	}

	return diffGen{conf: conf}
}

func (f diffGen) GenDiff(stats []FileStat) Diff {
	var diff Diff
	dayParts := strings.Split(time.Now().Format("2006-01-02"), "-")
	for _, stat := range stats {
		ext := f.findExtensionAlias(stat.Ext)
		filename := fmt.Sprintf("log.%s", ext)

		// handle files without extension. eg: Makefile, Dockerfile etc
		if !strings.HasPrefix(stat.Ext, ".") {
			filename = stat.Ext
		}

		msg, err := f.buildMessage(&stat, ext)
		if err != nil {
			panic("TODO")
		}

		diff.Changes = append(diff.Changes, Change{
			Dir:       filepath.Join(dayParts[0], dayParts[1], dayParts[2]),
			Filename:  filename,
			Text:      msg,
			Insertion: stat.Insert,
			Deletion:  stat.Delete,
		})
	}

	return diff
}

func (f diffGen) buildMessage(stat *FileStat, ext string) (string, error) {
	if f.conf == nil {
		return fmt.Sprintf(defaultCommitFormat, stat.Insert, stat.Delete), nil
	}

	// TODO: move template.New to init or once.Do
	commitTempl, err := template.New("test").Parse(f.conf.Commit.Template)
	if err != nil {
		return "", err
	}

	var commitMessage bytes.Buffer
	fileExt := f.removeDot(ext)
	now := time.Now() // TODO: get commit time instead now(usable especially with --dry-run)

	commonVars := struct {
		HM     string
		Hour   string
		Minute string
	}{
		HM:     now.Format("15:04"),
		Hour:   now.Format("15"),
		Minute: now.Format("04"),
	}

	commitTemplateVarMap := map[string]interface{}{
		"InsertCount": stat.Insert,
		"DeleteCount": stat.Delete,
		"Ext":         fileExt,
		"HM":          commonVars.HM,
		"Hour":        commonVars.Hour,
		"Minute":      commonVars.Minute,
	}

	if err := commitTempl.Execute(&commitMessage, commitTemplateVarMap); err != nil {
		return "", err
	}

	if tmpl, ok := f.conf.Templates[fileExt]; ok {
		wrapperTempl, err := template.New("wrapper").Parse(tmpl)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		varMap := map[string]interface{}{
			"Message": commitMessage.String(),
			"HM":      commonVars.HM,
			"Hour":    commonVars.Hour,
			"Minute":  commonVars.Minute,
		}

		if err := wrapperTempl.Execute(&buf, varMap); err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return fmt.Sprintf(defaultCommitFormat, stat.Insert, stat.Delete), nil
}

func (f diffGen) findExtensionAlias(ext string) string {
	if f.conf == nil {
		return ext
	}

	searchExt := f.removeDot(ext)

	for typ, aliases := range f.conf.Aliases {
		if typ != "default" {
			for ext, alias := range aliases {
				if ext == searchExt {
					return alias
				}
			}
		}
	}

	if aliases, ok := f.conf.Aliases["default"]; ok {
		for ext, alias := range aliases {
			if ext == searchExt {
				return alias
			}
		}
	}

	return searchExt
}

func (f diffGen) removeDot(ext string) string {
	if strings.HasPrefix(ext, ".") {
		return strings.Replace(ext, ".", "", 1)
	}

	return ext
}
