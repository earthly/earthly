package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/earthly/earthly/domain"
	"github.com/pkg/errors"
)

// TargetInput represents the conditions in which a target is invoked.
type TargetInput struct {
	// TargetCanonical is the identifier of this target in canonical form.
	TargetCanonical string `json:"targetCanonical"`
	// BuildArgs are the build args used to build this target.
	BuildArgs []BuildArgInput `json:"buildArgs"`
	// Platform is the target platform of the target.
	Platform string `json:"platform"`
}

// WithBuildArgInput returns a clone of the current target input, with a
// BuildArgInput added to it.
func (ti TargetInput) WithBuildArgInput(bai BuildArgInput) TargetInput {
	tiClone := ti.clone()
	for index, existingBai := range tiClone.BuildArgs {
		if existingBai.Name == bai.Name {
			// Existing build arg. Remove it so we can re-add it.
			tiClone.BuildArgs = append(tiClone.BuildArgs[:index], tiClone.BuildArgs[index+1:]...)
			break
		}
	}
	tiClone.BuildArgs = append(tiClone.BuildArgs, bai)
	return tiClone
}

// Equals compares to another TargetInput for equality.
func (ti TargetInput) Equals(other TargetInput) bool {
	if ti.TargetCanonical != other.TargetCanonical {
		return false
	}
	if ti.Platform != other.Platform {
		return false
	}
	if len(ti.BuildArgs) != len(other.BuildArgs) {
		return false
	}
	for index := range ti.BuildArgs {
		if !ti.BuildArgs[index].Equals(other.BuildArgs[index]) {
			return false
		}
	}
	return true
}

func (ti TargetInput) clone() TargetInput {
	tiCopy := TargetInput{
		TargetCanonical: ti.TargetCanonical,
		BuildArgs:       make([]BuildArgInput, 0, len(ti.BuildArgs)),
		Platform:        ti.Platform,
	}
	for _, bai := range ti.BuildArgs {
		tiCopy.BuildArgs = append(tiCopy.BuildArgs, bai.clone())
	}
	return tiCopy
}

func (ti TargetInput) cloneNoTag() (TargetInput, error) {
	targetStr := ""
	if ti.TargetCanonical != "" {
		target, err := domain.ParseTarget(ti.TargetCanonical)
		if err != nil {
			return TargetInput{}, err
		}
		target.Tag = ""
		targetStr = target.StringCanonical()
	}
	tiCopy := TargetInput{
		TargetCanonical: targetStr,
		BuildArgs:       make([]BuildArgInput, 0, len(ti.BuildArgs)),
		Platform:        ti.Platform,
	}
	for _, bai := range ti.BuildArgs {
		tiCopy.BuildArgs = append(tiCopy.BuildArgs, bai.clone())
	}
	return tiCopy, nil
}

// Hash returns a hash of the target input.
func (ti TargetInput) Hash() (string, error) {
	tiBytes, err := json.Marshal(&ti)
	if err != nil {
		return "", errors.Wrap(err, "serialize TargetInput when creating hash")
	}
	digest := sha256.Sum256(tiBytes)
	return hex.EncodeToString(digest[:]), nil
}

// HashNoTag returns a hash of the target input with tag info stripped away.
func (ti TargetInput) HashNoTag() (string, error) {
	tiNoTag, err := ti.cloneNoTag()
	if err != nil {
		return "", err
	}
	tiBytes, err := json.Marshal(&tiNoTag)
	if err != nil {
		return "", errors.Wrap(err, "serialize TargetInput when creating hash no tag")
	}
	digest := sha256.Sum256(tiBytes)
	return hex.EncodeToString(digest[:]), nil
}

// BuildArgInput represents the conditions in which a build arg is passed.
type BuildArgInput struct {
	// Name is the name of the build arg.
	Name string `json:"name"`
	// ConstantValue is the constant value of this build arg.
	ConstantValue string `json:"constantValue"`
	// DefaultValue represents the default value of the build arg.
	DefaultValue string `json:"defaultConstant"`
}

// IsDefaultValue returns whether the value of the BuildArgInput
// is set as the same as the default.
func (bai BuildArgInput) IsDefaultValue() bool {
	return bai.ConstantValue == bai.DefaultValue
}

// Equals compares to another BuildArgInput for equality.
func (bai BuildArgInput) Equals(other BuildArgInput) bool {
	if bai.Name != other.Name {
		return false
	}
	if bai.ConstantValue != other.ConstantValue {
		return false
	}
	if bai.DefaultValue != other.DefaultValue {
		return false
	}
	return true
}

func (bai BuildArgInput) clone() BuildArgInput {
	return BuildArgInput{
		Name:          bai.Name,
		ConstantValue: bai.ConstantValue,
		DefaultValue:  bai.DefaultValue,
	}
}
