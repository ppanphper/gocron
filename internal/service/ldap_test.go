package service

import (
	"github.com/ouqiang/gocron/internal/models"
	"testing"
)

func TestLdapService_Match(t *testing.T) {
	s := Ldap{}
	entry, err := s.Match("user01", "123456", models.LDAPSetting{})

	t.Log(entry.Attributes, err, s.GetEntryAttribute(entry, "mail"))
}
