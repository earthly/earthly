package cloud

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	CreatedAt  time.Time
	AcceptedAt time.Time
}

// InviteToOrg sends an email invitation to a user and asks for them to join an org.
func (c *client) InviteToOrg(ctx context.Context, invite *OrgInvitation) (string, error) {
	u := "/api/v0/invitations"

	req := &secretsapi.CreateInvitationRequest{
		OrgName:    invite.OrgName,
		Email:      invite.Email,
		Permission: invite.Permission,
		Message:    invite.Message,
	}

	status, body, err := c.doCall(ctx, http.MethodPost, u, withAuth(), withJSONBody(req))
	if err != nil {
		return "", err
	}

	if status != http.StatusCreated {
		return "", errors.Errorf("failed to send email invite: %s", body)
	}

	res := &secretsapi.CreateInvitationResponse{}
	err = jsonpb.UnmarshalString(body, res)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

// ListOrgMembers returns a collection of org members.
func (c *client) ListOrgMembers(ctx context.Context, orgName string) ([]*OrgMember, error) {
	u := fmt.Sprintf("/api/v1/organizations/%s/members", orgName)

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list org members: %s", body)
	}

	res := &secretsapi.ListOrgMembersResponse{}

	err = jsonpb.UnmarshalString(body, res)
	if err != nil {
		return nil, err
	}

	var members []*OrgMember

	for _, m := range res.Members {
		members = append(members, &OrgMember{
			UserEmail:  m.Email,
			OrgName:    m.OrgName,
			Permission: m.Permission,
		})
	}

	return members, nil
}

// UpdateOrgMember updates a member's permission in an org.
func (c *client) UpdateOrgMember(ctx context.Context, orgName, userEmail, permission string) error {
	u := fmt.Sprintf("/api/v1/organizations/%s/members/%s", orgName, userEmail)

	req := &secretsapi.UpdateOrgMemberRequest{
		Member: &secretsapi.OrgMember{
			Email:      userEmail,
			OrgName:    orgName,
			Permission: permission,
		},
	}

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth(), withJSONBody(req))
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to update member: %s", body)
	}

	return nil
}

// RemoveOrgMember removes a member from the org.
func (c *client) RemoveOrgMember(ctx context.Context, orgName, userEmail string) error {
	u := fmt.Sprintf("/api/v1/organizations/%s/members/%s", orgName, userEmail)

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to remove member: %s", body)
	}

	return nil
}

// AcceptInvite accepts the org invitation and adds the user to the org.
func (c *client) AcceptInvite(ctx context.Context, inviteCode string) error {
	u := "/api/v0/invitations/" + inviteCode

	status, body, err := c.doCall(ctx, http.MethodPost, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to remove member: %s", body)
	}

	return nil
}

// ListInvites returns a collection of organization invites and their status.
func (c *client) ListInvites(ctx context.Context, org string) ([]*OrgInvitation, error) {
	u := "/api/v0/invitations?org=" + org

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list invites: %s", body)
	}

	res := &secretsapi.ListInvitationsResponse{}

	err = jsonpb.UnmarshalString(body, res)
	if err != nil {
		return nil, err
	}

	var invites []*OrgInvitation

	for _, invite := range res.Invitations {
		in := &OrgInvitation{
			Email:      invite.RecipientEmail,
			OrgName:    org,
			Permission: invite.Permission,
			CreatedAt:  invite.CreatedAt.AsTime(),
		}
		if invite.AcceptedAt != nil {
			in.AcceptedAt = invite.AcceptedAt.AsTime()
		}
		invites = append(invites, in)
	}

	return invites, nil
}
