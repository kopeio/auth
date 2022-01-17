package httpserver

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	"kope.io/auth/pkg/session"
)

func (s *HTTPServer) mapUser(ctx context.Context, session *session.Session, info *session.UserInfo) error {
	providerID := session.ProviderId

	if providerID == "" {
		return fmt.Errorf("providerID not specified")
	}

	conf, err := s.config.AuthProvider(ctx, providerID)
	if err != nil {
		return fmt.Errorf("error reading configuration for %q: %w", providerID, err)
	}

	if conf == nil {
		return fmt.Errorf("no provider configuration for %q", providerID)
	}

	validator, err := buildValidator(conf.Spec.PermittedEmails)
	if err != nil {
		return fmt.Errorf("error building email validator: %w", err)
	}

	email := info.Email
	if email == "" {
		return fmt.Errorf("rejected login attempt without email: %v", info)
	}

	if !validator(email) {
		klog.Infof("rejected login attempt from %q", email)
		return fmt.Errorf("rejected login attempt, email %q not permitted", email)
	}

	// user, err := s.tokenStore.MapToUser(info, true)
	// if err != nil {
	// 	return "", err
	// }

	// glog.Infof("mapped %s to %s", email, user.UID)
	// return user.UID, nil

	session.ProviderEmail = email

	return nil
}

func buildValidator(permittedEmails []string) (func(string) bool, error) {
	allowAll := false
	var exact []string
	var suffixes []string
	for _, permittedEmail := range permittedEmails {
		wildcardCount := strings.Count(permittedEmail, "*")
		if wildcardCount == 0 {
			if permittedEmail == "" {
				// TODO: Move to validation?
				// TODO: Maybe ignore invalid rules?
				return nil, fmt.Errorf("empty permitEmail not allowed")
			}
			exact = append(exact, permittedEmail)
		} else if wildcardCount == 1 && strings.HasPrefix(permittedEmail, "*") {
			if permittedEmail == "*" {
				allowAll = true
			} else {
				// TODO: Block dangerous things i.e. require *@ or *. ?
				suffixes = append(suffixes, permittedEmail[1:])
			}
		} else {
			return nil, fmt.Errorf("Cannot parse permittedEmail rule: %q", permittedEmail)
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
