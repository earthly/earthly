package llbutil

import "regexp"

var invalidDockerTagCharsBeginningRe = regexp.MustCompile(`^[^\w]`)
var invalidDockerTagCharsMiddleRe = regexp.MustCompile(`[^\w.-]`)

// DockerTagSafe turns a string into a safe Docker tag.
func DockerTagSafe(tag string) string {
	if len(tag) == 0 {
		return "latest"
	}
	newTag := tag
	if len(tag) > 128 {
		newTag = newTag[:128]
	}
	newTag = invalidDockerTagCharsBeginningRe.ReplaceAllString(newTag, "_")
	if len(newTag) > 1 {
		newTag = string(newTag[0]) + invalidDockerTagCharsMiddleRe.ReplaceAllString(newTag[1:], "_")
	}
	return newTag
}
