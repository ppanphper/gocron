package service

import "testing"

func TestLdapService_Match(t *testing.T) {
	s := LdapService{}
	entry, err := s.Match("user01", "123456")

	t.Log(entry.Attributes, err, s.GetEntryAttribute(entry, "mail"))
}
