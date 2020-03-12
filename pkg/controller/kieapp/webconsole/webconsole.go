package webconsole

//go:generate go run -mod=vendor ../defaults/.packr/packr.go

import (
	"errors"
	"fmt"
	"github.com/RHsyseng/operator-utils/pkg/logs"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/webconsole/utils"
	"strconv"
	"strings"

	consolev1 "github.com/openshift/api/console/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CustomResourceDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

var log = logs.GetLogger("kieapp.webconsole")

func CreateSampleYAMLs() {
	log.Info("Loading Console yaml samples.")
	files, err := loadFilesOnlyWithBox("consoleyamlsamples")
	if err != nil {
		log.Error(err, "Failed to load consoleyamlsamples")
		return
	}

	resMap, err := applyMultipleWebConsoleYamls(files)
	if err != nil {
		log.Error(err, "consoleyamlsamples not successfully applied")
		return
	}

	for k, v := range resMap {
		log.Info("Creating YAML sample. file: % %s ", files[k], v)
	}
}

func applyMultipleWebConsoleYamls(yamlString []string) (map[int]string, error) {
	if yamlString == nil || len(yamlString) == 0 {
		return nil, errors.New("Empty yaml list")
	}
	resMap := make(map[int]string)
	for idx, yamlStr := range yamlString {
		_, err := create(yamlStr)
		if err != nil {
			resMap[idx] = err.Error()
		} else {
			resMap[idx] = "Applied"
		}
	}
	return resMap, nil
}

func create(yamlStr string) (bool, error) {
	client, err := utils.GetClient()
	if err != nil {
		return false, err
	}
	obj := &CustomResourceDefinition{}
	err = yaml.Unmarshal([]byte(yamlStr), obj)
	if err != nil {
		return false, err
	}

	snippetStr := obj.Annotations["consolesnippet"]
	var snippet bool = false
	if snippetStr != "" {
		tmp, err := strconv.ParseBool(snippetStr)
		if err != nil {
			return false, errors.New("Unable to parse snippet as boolean.")
		}
		snippet = tmp
	}
	title, _ := obj.Annotations["consoletitle"]
	desc, _ := obj.Annotations["consoledesc"]
	name, _ := obj.Annotations["consolename"]

	yamlSample := &consolev1.ConsoleYAMLSample{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: utils.TargetNamespace,
		},
		Spec: consolev1.ConsoleYAMLSampleSpec{
			TargetResource: v1.TypeMeta{
				APIVersion: obj.TypeMeta.APIVersion,
				Kind:       obj.TypeMeta.Kind,
			},
			Title:       consolev1.ConsoleYAMLSampleTitle(title),
			Description: consolev1.ConsoleYAMLSampleDescription(desc),
			YAML:        consolev1.ConsoleYAMLSampleYAML(yamlStr),
			Snippet:     snippet,
		},
	}
	yamlSample, err = client.ConsoleYAMLSample.Create(yamlSample)
	if err != nil {
		return false, err
	}
	return true, nil
}

func loadFilesOnlyWithBox(boxname string, path ...string) ([]string, error) {
	fullpath := "noop" //call for binary embedded files
	if len(path) == 2 {
		fullpath = strings.Join([]string{path[0], path[1]}, "/")
	}
	box := packr.New(boxname, fullpath)
	if box.List() == nil {
		return nil, fmt.Errorf("%s not found ", fullpath)
	}
	var files []string
	for _, filename := range box.List() {
		yamlStr, err := box.FindString(filename)
		if err != nil {
			log.Error("Failed to load file:", filename, err)
			continue
		}
		files = append(files, yamlStr)
	}
	return files, nil
}

//used by packr
func convertYamlToBinary() {
	box := packr.New("consoleyamlsamples", "./examples/consoleyamlsamples")
	if box.List() == nil {
		log.Error("Failed to create packr box on yaml sample folder.")
	}
}
