package jira

import "fmt"

// LinkService handles JIRA issue linking operations
type LinkService struct {
	client *Client
}

// NewLinkService creates a new link service
func NewLinkService(client *Client) *LinkService {
	return &LinkService{client: client}
}

// LinkBlocks creates a "Blocks" link between two issues
// blocker blocks blocked
func (s *LinkService) LinkBlocks(blocker, blocked string) error {
	return s.LinkIssues("Blocks", blocker, blocked)
}

// LinkIssues creates a link of the specified type between two issues
// outwardKey --[type]--> inwardKey
func (s *LinkService) LinkIssues(linkType, outwardKey, inwardKey string) error {
	req := LinkIssueRequest{
		Type: LinkType{
			Name: linkType,
		},
		OutwardIssue: &Issue{
			Key: outwardKey,
		},
		InwardIssue: &Issue{
			Key: inwardKey,
		},
	}

	if err := s.client.Do("POST", "/rest/api/2/issueLink", req, nil); err != nil {
		return fmt.Errorf("failed to link issues %s --[%s]--> %s: %w",
			outwardKey, linkType, inwardKey, err)
	}

	return nil
}

// LinkRelates creates a "Relates" link between two issues
func (s *LinkService) LinkRelates(fromKey, toKey string) error {
	return s.LinkIssues("Relates", fromKey, toKey)
}

// LinkDuplicates creates a "Duplicates" link between two issues
func (s *LinkService) LinkDuplicates(fromKey, toKey string) error {
	return s.LinkIssues("Duplicates", fromKey, toKey)
}

// LinkClones creates a "Clones" link between two issues
func (s *LinkService) LinkClones(fromKey, toKey string) error {
	return s.LinkIssues("Clones", fromKey, toKey)
}
