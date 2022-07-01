package cloud

import (
	"context"
	"net/http"

	secretsapi "github.com/earthly/cloud-api/secrets"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
)

// OrgInvitation can be used to invite a user to become a member in an org.
type OrgInvitation struct {
	Name       string
	Email      string
	Permission string
	Message    string
	OrgName    string
}

// InviteToOrg sends an email invitation to a user and asks for them to join an org.
func (c *client) InviteToOrg(ctx context.Context, invite *OrgInvitation) (string, error) {
	u := "/api/v1/invitations"

	req := &secretsapi.CreateInvitationRequest{
		Name:       invite.Name,
		OrgName:    invite.OrgName,
		Email:      invite.Email,
		Permission: invite.Permission,
		Message:    invite.Message,
	}

	status, body, err := c.doCall(ctx, http.MethodPost, u, withAuth(), withJSONBody(req))
	if err != nil {
		return "", err
	}

	if status != http.StatusOK {
		return "", errors.Errorf("failed to send email invite: %s", body)
	}

	res := &secretsapi.CreateInvitationResponse{}
	err = jsonpb.UnmarshalString(body, res)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}
