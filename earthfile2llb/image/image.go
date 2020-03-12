package image

import specs "github.com/opencontainers/image-spec/specs-go/v1"

// Image is a partial of the standard Image struct defined as part of the image opencontainers spec
// at https://github.com/opencontainers/image-spec/blob/master/specs-go/v1/config.go#L82
type Image struct {
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	Config       specs.ImageConfig `json:"config"`
}

// NewImage returns a new image.
func NewImage() *Image {
	return &Image{}
}

// Clone creates a copy of the image.
func (img *Image) Clone() *Image {
	if img == nil {
		return NewImage()
	}
	clone := &Image{
		Architecture: img.Architecture,
		OS:           img.OS,
		Config: specs.ImageConfig{
			User:         img.Config.User,
			Env:          make([]string, len(img.Config.Env)),
			Entrypoint:   make([]string, len(img.Config.Entrypoint)),
			Cmd:          make([]string, len(img.Config.Cmd)),
			WorkingDir:   img.Config.WorkingDir,
			StopSignal:   img.Config.StopSignal,
			ExposedPorts: make(map[string]struct{}),
			Volumes:      make(map[string]struct{}),
			Labels:       make(map[string]string),
		},
	}
	copy(clone.Config.Env, img.Config.Env)
	copy(clone.Config.Entrypoint, img.Config.Entrypoint)
	copy(clone.Config.Cmd, img.Config.Cmd)
	if img.Config.ExposedPorts != nil {
		for k, v := range img.Config.ExposedPorts {
			clone.Config.ExposedPorts[k] = v
		}
	}
	if img.Config.Volumes != nil {
		for k, v := range img.Config.Volumes {
			clone.Config.Volumes[k] = v
		}
	}
	if img.Config.Labels != nil {
		for k, v := range img.Config.Labels {
			clone.Config.Labels[k] = v
		}
	}
	return clone
}
