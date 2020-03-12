package webconsole

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestLoadSampleYamlFiles(t *testing.T) {
	files, err := loadFilesOnlyWithBox("consoleyamlsamples", "./examples", "consoleyamlsamples")
	if err != nil {
		log.Error(err, "Failed ot load consoleyamlsamples")
		return
	}

	resMap, err := applyMultipleWebConsoleYamls(files)
	if err != nil {
		log.Error(err, "consoleyamlsamples not successfully applied")
	}

	for k, v := range resMap {
		if v != "Applied" {
			log.Warnf("file % was not successfully applied ", files[k], v)
		}
	}
}

func TestConsoleYamlSamples(t *testing.T) {
	files, err := loadTestFiles("name", "./examples", "consoleyamlsamples")

	size := len(files)
	resMap, err := applyMultipleWebConsoleYamls(files)
	if err != nil {
		t.Errorf("error : not able to process yaml files %v", err)
	}
	assert.Equal(t, size, len(resMap), "Number of processed file doesn't match expected.")
}

func TestWebConsole(t *testing.T) {
	loadTestFiles("name", "./examples", "webconsole")
}

func loadTestFiles(boxname string, path string, folder string) ([]string, error) {
	fullpath := strings.Join([]string{path, folder}, "/")

	fileList, err := ioutil.ReadDir(fullpath)
	if err != nil {
		fmt.Println(fmt.Errorf("%s not found with io ", fullpath))
		return nil, fmt.Errorf("%s not found with io ", fullpath)
	}

	var files []string
	for _, filename := range fileList {
		yamlStr, err := ioutil.ReadFile(fullpath + "/" + filename.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		files = append(files, string(yamlStr))
	}
	//for _, f := range files {
	//	a := []rune(f)
	//	fmt.Println("filename: ", string(a[0: 10]))
	//
	//	err := ApplyWebConsoleYaml(f)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	return files, nil
}
