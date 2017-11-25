package portal

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	"kope.io/auth/pkg/oauth/session"
)

func (s *HTTPServer) mapUser(session *session.Session, info *session.UserInfo) (types.UID, error) {
	providerID := session.ProviderId

	if providerID == "" {
		return "", fmt.Errorf("providerID not specified")
	}

	conf, err := s.config.AuthProvider(providerID)
	if err != nil {
		return "", fmt.Errorf("error reading configuration for %q: %v", providerID, err)
	}

	if conf == nil {
		return "", fmt.Errorf("no provider configuration for %q", providerID)
	}

	validator, err := buildValidator(conf.PermitEmails)
	if err != nil {
		return "", fmt.Errorf("error building email validator: %v", err)
	}

	email := info.Email
	if email == "" {
		return "", fmt.Errorf("rejected login attempt without email: %s", info)
	}

	if !validator(email) {
		glog.Infof("rejected login attempt from %q", email)
		return "", fmt.Errorf("rejected login attempt, email %q not permitted", email)
	}

	user, err := s.tokenStore.MapToUser(info, true)
	if err != nil {
		return "", err
	}

	glog.Infof("mapped %s to %s", email, user.UID)
	return user.UID, nil
}

func buildValidator(permitEmails []string) (func(string) bool, error) {
	allowAll := false
	var exact []string
	var suffixes []string
	for _, permitEmail := range permitEmails {
		wildcardCount := strings.Count(permitEmail, "*")
		if wildcardCount == 0 {
			if permitEmail == "" {
				// TODO: Move to validation?
				// TODO: Maybe ignore invalid rules?
				return nil, fmt.Errorf("empty permitEmail not allowed")
			}
			exact = append(exact, permitEmail)
		} else if wildcardCount == 1 && strings.HasPrefix(permitEmail, "*") {
			if permitEmail == "*" {
				allowAll = true
			} else {
				// TODO: Block dangerous things i.e. require *@ or *. ?
				suffixes = append(suffixes, permitEmail[1:])
			}
		} else {
			return nil, fmt.Errorf("Cannot parse permitEmail rule: %q", permitEmail)
		}
	}

	validator := func(email string) bool {
		if email == "" {
			return false
		}
		email = strings.TrimSpace(strings.ToLower(email))
		if allowAll {
			return true
		}
		for _, s := range exact {
			if s == email {
				return true
			}
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(email, suffix) {
				return true
			}
		}

		return false
	}
	return validator, nil
}
