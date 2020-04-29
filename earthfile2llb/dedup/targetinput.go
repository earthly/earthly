package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/earthly/earthly/domain"
)

// TargetInput represents the conditions in which a target is invoked.
type TargetInput struct {
	// TargetCanonical is the identifier of this target in canonical form.
	TargetCanonical string `json:"targetCanonical"`
	// BuildArgs are the build args used to build this target.
	BuildArgs []BuildArgInput `json:"buildArgs"`
}

// Equals compares to another TargetInput for equality.
func (ti TargetInput) Equals(other TargetInput) bool {
	if ti.TargetCanonical != other.TargetCanonical {
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
	}
	for _, bai := range ti.BuildArgs {
		baiCopy, err := bai.cloneNoTag()
		if err != nil {
			return TargetInput{}, err
		}
		tiCopy.BuildArgs = append(tiCopy.BuildArgs, baiCopy)
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
	// IsConstant represents whether the build arg has a constant value.
	IsConstant bool `json:"isConstant"`
	// ConstantValue is the constant value of this build arg. Only set if IsConstant is true.
	ConstantValue string `json:"constantValue"`
	// DefaultValue represents the default value of the build arg.
	DefaultValue string `json:"defaultConstant"`
	// VariableFromInput is the definition of the conditions in which the build arg was formed.
	// It is only set if IsConstant is false.
	VariableFromInput VariableFromInput `json:"variableFromInput"`
}

// IsDefaultValue returns whether the BuildArgInput is constant and its value
// is set as the same as the default.
func (bai BuildArgInput) IsDefaultValue() bool {
	return bai.IsConstant && bai.ConstantValue == bai.DefaultValue
}

// Equals compares to another BuildArgInput for equality.
func (bai BuildArgInput) Equals(other BuildArgInput) bool {
	if bai.Name != other.Name {
		return false
	}
	if bai.IsConstant != other.IsConstant {
		return false
	}
	if bai.ConstantValue != other.ConstantValue {
		return false
	}
	if bai.DefaultValue != other.DefaultValue {
		return false
	}
	if !bai.VariableFromInput.Equals(other.VariableFromInput) {
		return false
	}
	return true
}

func (bai BuildArgInput) cloneNoTag() (BuildArgInput, error) {
	vfiCopy, err := bai.VariableFromInput.cloneNoTag()
	if err != nil {
		return BuildArgInput{}, err
	}
	return BuildArgInput{
		Name:              bai.Name,
		IsConstant:        bai.IsConstant,
		ConstantValue:     bai.ConstantValue,
		DefaultValue:      bai.DefaultValue,
		VariableFromInput: vfiCopy,
	}, nil
}

// VariableFromInput represents the conditions in which a variable build arg is passed.
type VariableFromInput struct {
	// TargetInput is the target where this variable's expression is executed.
	TargetInput TargetInput `json:"targetInput"`
	// Index is the order number of the expression for the build arg. Starts at 0 for each target.
	Index int `json:"index"`
}

// Equals compares to another VariableFromInput for equality.
func (vfi VariableFromInput) Equals(other VariableFromInput) bool {
	if !vfi.TargetInput.Equals(other.TargetInput) {
		return false
	}
	if vfi.Index != other.Index {
		return false
	}
	return true
}

func (vfi VariableFromInput) cloneNoTag() (VariableFromInput, error) {
	tiCopy, err := vfi.TargetInput.cloneNoTag()
	if err != nil {
		return VariableFromInput{}, err
	}
	return VariableFromInput{
		TargetInput: tiCopy,
		Index:       vfi.Index,
	}, nil
}
