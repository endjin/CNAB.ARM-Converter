package generator

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

func TestGenerateTemplate(t *testing.T) {

	bundlePath := "testdata/bundle.json"
	generatedOutputPath := "testdata/azuredeploy-generated.json"
	expectedOutputPath := "testdata/azuredeploy.json"

	err := GenerateTemplate(bundlePath, generatedOutputPath, true, true, "latest")
	if err != nil {
		t.Errorf("GenerateTemplate failed: %s", err.Error())
	}

	expectedBytes, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("failed reading expected output: %s", err)
	}
	expected := string(expectedBytes)

	generatedBytes, err := ioutil.ReadFile(generatedOutputPath)
	if err != nil {
		t.Fatalf("failed reading generated output: %s", err)
	}
	generated := string(generatedBytes)

	assert.Equal(t, expected, generated)
}
