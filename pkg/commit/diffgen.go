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

type diffGen struct {
	conf *config.Config
}

func NewDiffGenerator(conf *config.Config) diffGen {
	return diffGen{conf: conf}
}

func (f diffGen) GenDiff(stats []FileStat, commitMeta Meta) (*Diff, error) {
	mergedFileStats := make(map[string]FileStat)
	for _, stat := range stats {
		fileType := f.decideNewExtension(f.clearDot(stat.Ext))
		if _, ok := mergedFileStats[fileType]; ok {
			mergedFileStats[fileType] = FileStat{
				Insert: mergedFileStats[fileType].Insert + stat.Insert,
				Delete: mergedFileStats[fileType].Delete + stat.Delete,
				Ext:    stat.Ext,
			}
		} else {
			mergedFileStats[fileType] = FileStat{
				Insert: stat.Insert,
				Delete: stat.Delete,
				Ext:    stat.Ext,
			}
		}
	}

	var diff Diff
	dayParts := strings.Split(time.Now().Format("2006-01-02"), "-")

	for fileType, stat := range mergedFileStats {
		// handle files without extension. eg: Makefile, Dockerfile etc
		filename := fileType
		if strings.HasPrefix(stat.Ext, ".") {
			filename = fmt.Sprintf("log.%s", fileType)
		}

		contentText, err := f.generateContentText(&stat, fileType, &commitMeta)
		if err != nil {
			return nil, err
		}

		text, err := f.generateCode(fileType, contentText, commitMeta)
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

func (f diffGen) generateCode(fileType string, contentText string, commitMeta Meta) (string, error) {
	if tmpl, ok := f.conf.Templates[fileType]; ok {
		codeTemplate, err := template.New("code").Parse(tmpl)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		varMap := map[string]interface{}{
			"Message": contentText,
			"HM":      commitMeta.Time.Format("15:04"),
			"Hour":    commitMeta.Time.Format("15"),
			"Minute":  commitMeta.Time.Format("04"),
		}

		if err := codeTemplate.Execute(&buf, varMap); err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return contentText, nil
}

func (f diffGen) generateContentText(stat *FileStat, fileType string, commitMeta *Meta) (string, error) {
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
		HM:     commitMeta.Time.Format("15:04"),
		Hour:   commitMeta.Time.Format("15"),
		Minute: commitMeta.Time.Format("04"),
	}

	commitTemplateVarMap := map[string]interface{}{
		"InsertCount":   stat.Insert,
		"DeleteCount":   stat.Delete,
		"FileExtension": stat.Ext,
		"FileType":      fileType,
		"HM":            commonVars.HM,
		"Hour":          commonVars.Hour,
		"Minute":        commonVars.Minute,
	}

	var contentText bytes.Buffer
	if err := commitTemplate.Execute(&contentText, commitTemplateVarMap); err != nil {
		return "", fmt.Errorf("code content generate failed. Please check commit template in config. err: %w", err)
	}

	return contentText.String(), nil
}

func (f diffGen) decideNewExtension(ext string) string {
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
