package commit

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/mesuutt/git-mirror/pkg/config"
)

const defaultCommitFormat = "%d insertion(s), %d deletion(s)"

type Generator struct {
	conf *config.Config
}

func NewGenerator(configPath string) Generator {
	conf, err := config.ReadConfig(configPath)
	if err != nil {
		return Generator{conf: nil}
	}

	return Generator{conf: conf}
}

func (f Generator) GenCommit(stat *FileStat) Commit {
	ext := f.findRealExtension(stat)
	filename := fmt.Sprintf("log.%s", ext)

	// handle files without extension. eg: Makefile, Dockerfile etc
	if !strings.HasPrefix(stat.Ext, ".") {
		filename = stat.Ext
	}

	msg, err := f.buildMessage(stat)
	if err != nil {
		panic("TODO")
	}

	return Commit{
		Filename: filename,
		Message:  msg,
	}
}

func (f Generator) buildMessage(stat *FileStat) (string, error) {
	if f.conf == nil {
		return fmt.Sprintf(defaultCommitFormat, stat.Insert, stat.Delete), nil
	}

	// TODO: move template.New to init or once.Do
	commitTempl, err := template.New("test").Parse(f.conf.Commit.Template)
	if err != nil {
		return "", err
	}

	var commitMessage bytes.Buffer
	fileExt := f.removeDot(stat.Ext)
	now := time.Now()

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

func (f Generator) findRealExtension(stat *FileStat) string {
	searchExt := f.removeDot(stat.Ext)

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

func (f Generator) removeDot(ext string) string {
	if strings.HasPrefix(ext, ".") {
		return strings.Replace(ext, ".", "", 1)
	}

	return ext
}
