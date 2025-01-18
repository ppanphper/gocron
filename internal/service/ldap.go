package service

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/ouqiang/gocron/internal/models"
	"strings"
)

type Ldap struct{}

var LdapService Ldap

func (Ldap) Match(username, password string, ldapSetting models.LDAPSetting) (ldap.Entry, error) {
	l, err := ldap.DialURL(ldapSetting.Url)
	if err != nil {
		return ldap.Entry{}, err
	}
	defer l.Close()

	err = l.Bind(ldapSetting.BindDn, ldapSetting.BindPassword)
	if err != nil {
		return ldap.Entry{}, err
	}

	attributes := []string{"dn", "sn"}
	emailAttribute := ldapSetting.LdapEmailAttribute
	if emailAttribute != "" {
		attributes = append(attributes, emailAttribute)
	}

	req := ldap.NewSearchRequest(
		ldapSetting.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf(ldapSetting.FilterRule, username),
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
		return ldap.Entry{}, errors.New("LDAP数据异常,用户不存在或者存在多个对应的用户")
	}

	err = l.Bind(sr.Entries[0].DN, password)
	if err != nil {
		return ldap.Entry{}, err
	}

	return *sr.Entries[0], nil
}

func (Ldap) GetEntryAttribute(entry ldap.Entry, attributeName string) []string {
	for _, attribute := range entry.Attributes {
		if strings.ToLower(attribute.Name) == strings.ToLower(attributeName) {
			return attribute.Values
		}
	}
	return []string{}
}

// Enable LDAP是否可用
func (Ldap) Enable(ldapSetting models.LDAPSetting) bool {
	return ldapSetting.Enable == "1"
}
