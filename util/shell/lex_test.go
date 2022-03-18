package shell

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellParserMandatoryEnvVars(t *testing.T) {
	var newWord string
	var err error
	shlex := NewLex('\\')
	setEnvs := []string{"VAR=plain", "ARG=x"}
	emptyEnvs := []string{"VAR=", "ARG=x"}
	unsetEnvs := []string{"ARG=x"}

	noEmpty := "${VAR:?message here$ARG}"
	noUnset := "${VAR?message here$ARG}"

	// disallow empty
	newWord, err = shlex.ProcessWord(noEmpty, setEnvs)
	require.NoError(t, err)
	require.Equal(t, "plain", newWord)

	_, err = shlex.ProcessWord(noEmpty, emptyEnvs)
	require.Error(t, err)
	require.Contains(t, err.Error(), "message herex")

	_, err = shlex.ProcessWord(noEmpty, unsetEnvs)
	require.Error(t, err)
	require.Contains(t, err.Error(), "message herex")

	// disallow unset
	newWord, err = shlex.ProcessWord(noUnset, setEnvs)
	require.NoError(t, err)
	require.Equal(t, "plain", newWord)

	newWord, err = shlex.ProcessWord(noUnset, emptyEnvs)
	require.NoError(t, err)
	require.Equal(t, "", newWord)

	_, err = shlex.ProcessWord(noUnset, unsetEnvs)
	require.Error(t, err)
	require.Contains(t, err.Error(), "message herex")
}

func TestShellParser4EnvVars(t *testing.T) {
	fn := "envVarTest"
	lineCount := 0

	file, err := os.Open(fn)
	require.NoError(t, err)
	defer file.Close()

	shlex := NewLex('\\')
	scanner := bufio.NewScanner(file)
	envs := []string{"PWD=/home", "SHELL=bash", "KOREAN=한국어", "NULL="}
	envsMap := BuildEnvs(envs)
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Skip comments and blank lines
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		words := strings.Split(line, "|")
		require.Equal(t, 3, len(words))

		platform := strings.TrimSpace(words[0])
		source := strings.TrimSpace(words[1])
		expected := strings.TrimSpace(words[2])

		// Key W=Windows; A=All; U=Unix
		if platform != "W" && platform != "A" && platform != "U" {
			t.Fatalf("Invalid tag %s at line %d of %s. Must be W, A or U", platform, lineCount, fn)
		}

		if ((platform == "W" || platform == "A") && runtime.GOOS == "windows") ||
			((platform == "U" || platform == "A") && runtime.GOOS != "windows") {
			newWord, err := shlex.ProcessWord(source, envs)
			if expected == "error" {
				require.Errorf(t, err, "input: %q, result: %q", source, newWord)
			} else {
				require.NoError(t, err, "at line %d of %s", lineCount, fn)
				require.Equal(t, expected, newWord, "at line %d of %s", lineCount, fn)
			}

			newWord, err = shlex.ProcessWordWithMap(source, envsMap)
			if expected == "error" {
				require.Errorf(t, err, "input: %q, result: %q", source, newWord)
			} else {
				require.NoError(t, err, "at line %d of %s", lineCount, fn)
				require.Equal(t, expected, newWord, "at line %d of %s", lineCount, fn)
			}
		}
	}
}

func TestShellParser4Words(t *testing.T) {
	fn := "wordsTest"

	file, err := os.Open(fn)
	if err != nil {
		t.Fatalf("Can't open '%s': %s", err, fn)
	}
	defer file.Close()

	const (
		modeNormal = iota
		modeOnlySetEnv
	)
	for _, mode := range []int{modeNormal, modeOnlySetEnv} {
		var envs []string
		shlex := NewLex('\\')
		if mode == modeOnlySetEnv {
			shlex.RawQuotes = true
			shlex.SkipUnsetEnv = true
		}
		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			line := scanner.Text()
			lineNum = lineNum + 1

			if strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "ENV ") {
				line = strings.TrimLeft(line[3:], " ")
				envs = append(envs, line)
				continue
			}

			words := strings.Split(line, "|")
			if len(words) != 2 {
				t.Fatalf("Error in '%s'(line %d) - should be exactly one | in: %q", fn, lineNum, line)
			}
			test := strings.TrimSpace(words[0])
			expected := strings.Split(strings.TrimLeft(words[1], " "), ",")

			// test for ProcessWords
			result, err := shlex.ProcessWords(test, envs)

			if err != nil {
				result = []string{"error"}
			}

			if len(result) != len(expected) {
				t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
			}
			for i, w := range expected {
				if w != result[i] {
					t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
				}
			}

			// test for ProcessWordsWithMap
			result, err = shlex.ProcessWordsWithMap(test, BuildEnvs(envs))

			if err != nil {
				result = []string{"error"}
			}

			if len(result) != len(expected) {
				t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
			}
			for i, w := range expected {
				if w != result[i] {
					t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
				}
			}
		}
	}
}

func TestGetEnv(t *testing.T) {
	sw := &shellWord{envs: nil}

	getEnv := func(name string) string {
		value, _ := sw.getEnv(name)
		return value
	}
	sw.envs = BuildEnvs([]string{})
	if getEnv("foo") != "" {
		t.Fatal("2 - 'foo' should map to ''")
	}

	sw.envs = BuildEnvs([]string{"foo"})
	if getEnv("foo") != "" {
		t.Fatal("3 - 'foo' should map to ''")
	}

	sw.envs = BuildEnvs([]string{"foo="})
	if getEnv("foo") != "" {
		t.Fatal("4 - 'foo' should map to ''")
	}

	sw.envs = BuildEnvs([]string{"foo=bar"})
	if getEnv("foo") != "bar" {
		t.Fatal("5 - 'foo' should map to 'bar'")
	}

	sw.envs = BuildEnvs([]string{"foo=bar", "car=hat"})
	if getEnv("foo") != "bar" {
		t.Fatal("6 - 'foo' should map to 'bar'")
	}
	if getEnv("car") != "hat" {
		t.Fatal("7 - 'car' should map to 'hat'")
	}

	// Make sure we grab the last 'car' in the list
	sw.envs = BuildEnvs([]string{"foo=bar", "car=hat", "car=bike"})
	if getEnv("car") != "bike" {
		t.Fatal("8 - 'car' should map to 'bike'")
	}
}
