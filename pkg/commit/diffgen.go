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

const defaultLogMsgFormat = "%d insertion(s), %d deletion(s)"

type diffGen struct {
	conf *config.Config
}

func NewDiffGenerator(conf *config.Config) diffGen {
	return diffGen{conf: conf}
}

func (f diffGen) GenDiff(stats []FileStat, commitInfo CommitInfo) (*Diff, error) {
	var diff Diff
	dayParts := strings.Split(time.Now().Format("2006-01-02"), "-")
	for _, stat := range stats {
		fileExt := f.decideFileExtension(f.clearDot(stat.Ext))
		filename := fmt.Sprintf("log.%s", fileExt)

		// handle files without extension. eg: Makefile, Dockerfile etc
		if !strings.HasPrefix(stat.Ext, ".") {
			filename = stat.Ext
		}

		text, err := f.generateLog(&stat, fileExt, commitInfo)
		if err != nil {
			return nil, fmt.Errorf("log text generation failed, err: %w", err)
		}

		diff.Changes = append(diff.Changes, Change{
			Dir:       filepath.Join(dayParts[0], dayParts[1], dayParts[2]),
			Filename:  filename,
			Text:      text,
			Insertion: stat.Insert,
			Deletion:  stat.Delete,
		})
	}

	return &diff, nil
}

func (f diffGen) generateLog(stat *FileStat, fileExt string, commitInfo CommitInfo) (string, error) {
	if f.conf == nil {
		return fmt.Sprintf(defaultLogMsgFormat, stat.Insert, stat.Delete), nil
	}

	// TODO: we can move template.New to constructor, init or once.Do
	commitTemplate, err := template.New("log").Parse(f.conf.Commit.Template)
	if err != nil {
		return "", err
	}

	commonVars := struct {
		HM     string
		Hour   string
		Minute string
	}{
		HM:     commitInfo.Time.Format("15:04"),
		Hour:   commitInfo.Time.Format("15"),
		Minute: commitInfo.Time.Format("04"),
	}

	commitTemplateVarMap := map[string]interface{}{
		"InsertCount": stat.Insert,
		"DeleteCount": stat.Delete,
		"Ext":         fileExt,
		"HM":          commonVars.HM,
		"Hour":        commonVars.Hour,
		"Minute":      commonVars.Minute,
	}

	var contentText bytes.Buffer
	if err := commitTemplate.Execute(&contentText, commitTemplateVarMap); err != nil {
		return "", fmt.Errorf("code content generate failed. Please check commit template in config. err: %w", err)
	}

	if tmpl, ok := f.conf.Templates[fileExt]; ok {
		codeTemplate, err := template.New("code").Parse(tmpl)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		varMap := map[string]interface{}{
			"Message": contentText.String(),
			"HM":      commonVars.HM,
			"Hour":    commonVars.Hour,
			"Minute":  commonVars.Minute,
		}

		if err := codeTemplate.Execute(&buf, varMap); err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return contentText.String(), nil
}

func (f diffGen) decideFileExtension(ext string) string {
	if f.conf == nil {
		return ext
	}

	for typ, aliases := range f.conf.Overwrites {
		if typ != "default" {
			for fileType, alias := range aliases {
				if ext == fileType {
					return alias
				}
			}
		}
	}

	if aliases, ok := f.conf.Overwrites["default"]; ok {
		for fileType, alias := range aliases {
			if ext == fileType {
				return alias
			}
		}
	}

	return ext
}

func (f diffGen) clearDot(ext string) string {
	if strings.HasPrefix(ext, ".") {
		return strings.Replace(ext, ".", "", 1)
	}

	return ext
}
