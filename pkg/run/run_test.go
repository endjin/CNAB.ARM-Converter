package run

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestBuildPorterCommandParams(t *testing.T) {
	cnabBundleTag := "myregistry.io/mybundle:0.1.0"
	cnabAction := "install"
	cnabInstallationName := "mybundle1"

	cmdParams := buildPorterCommandParams(cnabInstallationName, cnabAction, cnabBundleTag)

	expectedPattern :=
		`install mybundle1 -d azure --tag myregistry.io\/mybundle:0\.1\.0 --cred \/tmp\/cnabarmdriver(.*)\/mybundle1-creds\.json --parameter-set \/tmp\/cnabarmdriver(.*)\/mybundle1-params\.json`

	cmdParamsStr := strings.Join(cmdParams, " ")
	t.Log(cmdParamsStr)
	match, _ := regexp.MatchString(expectedPattern, strings.Join(cmdParams, " "))

	assert.Equal(t, match, true)
}

func TestGenerateParamsFile(t *testing.T) {
	os.Setenv("CNAB_PARAM_foo", "1")
	os.Setenv("CNAB_PARAM_bar", "2")

	cnabInstallationName := "mybundle1"
	path, err := generateParamsFile(cnabInstallationName)

	assert.NilError(t, err)

	content, _ := ioutil.ReadFile(path)
	text := string(content)

	expected :=
		`{"schemaVersion":"","name":"mybundle1","created":"0001-01-01T00:00:00Z","modified":"0001-01-01T00:00:00Z","parameters":[{"name":"foo","source":{"value":"1"}},{"name":"bar","source":{"value":"2"}}]}`

	assert.Equal(t, expected, text)
}

func TestGenerateCredsFile(t *testing.T) {
	os.Setenv("CNAB_CRED_foo", "1")
	os.Setenv("CNAB_CRED_FILE_bar", base64.StdEncoding.EncodeToString([]byte("2")))

	cnabInstallationName := "mybundle1"
	path, err := generateCredsFile(cnabInstallationName)

	assert.NilError(t, err)

	content, _ := ioutil.ReadFile(path)
	text := string(content)

	expectedPattern :=
		`{"schemaVersion":"","name":"mybundle1","created":"0001-01-01T00:00:00Z","modified":"0001-01-01T00:00:00Z","credentials":\[{"name":"foo","source":{"env":"CNAB_CRED_foo"}},{"name":"bar","source":{"path":"/tmp/cnabarmdriver(.*)/bar"}}]}`

	match, _ := regexp.MatchString(expectedPattern, text)

	assert.Equal(t, match, true)
}
