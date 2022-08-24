package service

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/ouqiang/gocron/internal/models"
)

type LdapService struct{}

func (LdapService) Match(username, password string) (ldap.Entry, error) {
	settings := models.Setting{}
	l, err := ldap.DialURL(settings.Get(models.LdapCode, models.LdapKeyUrl).Value)
	if err != nil {
		return ldap.Entry{}, err
	}
	defer l.Close()

	err = l.Bind(settings.Get(models.LdapCode, models.LdapKeyBindDn).Value, settings.Get(models.LdapCode, models.LdapKeyBindPassword).Value)
	if err != nil {
		return ldap.Entry{}, err
	}

	attributes := []string{"dn", "sn"}
	emailAttribute := settings.Get(models.LdapCode, models.LdapEmailAttribute).Value
	if emailAttribute != "" {
		attributes = append(attributes, emailAttribute)
	}

	req := ldap.NewSearchRequest(
		settings.Get(models.LdapCode, models.LdapKeyBaseDn).Value,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf(settings.Get(models.LdapCode, models.LdapKeyFilterRule).Value, username),
		attributes,
		nil,
	)

	sr, err := l.Search(req)
	if err != nil {
		return ldap.Entry{}, err
	}

	// 如果没有数据返回或者超过1条数据返回,这对于用户认证而言都是不允许的.
	// 前这意味着没有查到用户,后者意味着存在重复数据
	if len(sr.Entries) != 1 {
		return ldap.Entry{}, errors.New("user does not exist or too many entries returned")
	}

	err = l.Bind(sr.Entries[0].DN, password)
	if err != nil {
		return ldap.Entry{}, err
	}

	return *sr.Entries[0], nil
}

func (LdapService) GetEntryAttribute(entry ldap.Entry, attributeName string) []string {
	for _, attribute := range entry.Attributes {
		if attribute.Name == attributeName {
			return attribute.Values
		}
	}
	return []string{}
}
