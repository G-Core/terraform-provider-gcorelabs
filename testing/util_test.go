package testing

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gcloud/gcorecloud-go"
	th "gcloud/gcorecloud-go/testhelper"
)

func TestWaitFor(t *testing.T) {
	err := gcorecloud.WaitFor(2, func() (bool, error) {
		return true, nil
	})
	th.CheckNoErr(t, err)
}

func TestWaitForTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gcorecloud.WaitFor(1, func() (bool, error) {
		return false, nil
	})
	require.Error(t, err)
	th.AssertEquals(t, "a timeout occurred", err.Error())
}

func TestWaitForError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gcorecloud.WaitFor(2, func() (bool, error) {
		return false, errors.New("error has occurred")
	})
	require.Error(t, err)
	th.AssertEquals(t, "error has occurred", err.Error())
}

func TestWaitForPredicateExceed(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gcorecloud.WaitFor(1, func() (bool, error) {
		time.Sleep(4 * time.Second)
		return false, errors.New("just wasting time")
	})
	require.Error(t, err)
	th.AssertEquals(t, "a timeout occurred", err.Error())
}

func TestNormalizeURL(t *testing.T) {
	urls := []string{
		"NoSlashAtEnd",
		"SlashAtEnd/",
	}
	expected := []string{
		"NoSlashAtEnd/",
		"SlashAtEnd/",
	}
	for i := 0; i < len(expected); i++ {
		th.CheckEquals(t, expected[i], gcorecloud.NormalizeURL(urls[i]))
	}
}

func TestNormalizePathURL(t *testing.T) {
	baseDir, _ := os.Getwd()

	rawPath := "template.yaml"
	basePath, _ := filepath.Abs(".")
	result, _ := gcorecloud.NormalizePathURL(basePath, rawPath)
	expected := strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "template.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	googleURL := "http://www.google.com"
	testPath := "very/nested/file.yaml"

	rawPath = googleURL
	basePath, _ = filepath.Abs(".")
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = googleURL
	th.CheckEquals(t, expected, result)

	rawPath = testPath
	basePath, _ = filepath.Abs(".")
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	rawPath = testPath
	basePath = googleURL
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml/"
	basePath = "http://www.google.com/"
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = testPath
	basePath = "http://www.google.com/even/more"
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/even/more/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = testPath
	basePath = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more"}, "/")
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more/very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml/"
	basePath = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more"}, "/")
	result, _ = gcorecloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more/very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)
}

func TestStripLastSlashURL(t *testing.T) {
	testCase := []map[string]string{
		{
			"url":    "http://test.com/1/1//////",
			"result": "http://test.com/1/1",
		},
		{
			"url":    "http://test.com/1/1",
			"result": "http://test.com/1/1",
		},
		{
			"url":    "http://test.com/1/1/",
			"result": "http://test.com/1/1",
		},
		{
			"url":    "",
			"result": "",
		},
		{
			"url":    "/",
			"result": "",
		},
		{
			"url":    "///////",
			"result": "",
		},
	}

	for _, m := range testCase {
		require.Equal(t, m["result"], gcorecloud.StripLastSlashURL(m["url"]))
	}
}
