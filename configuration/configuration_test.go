package configuration_test

import (
	"testing"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/configuration"
)

func TestNewContactFormConfigurationWithoutEnvVars(t *testing.T) {
	cfg := configuration.NewContactFormConfiguration()
	if cfg == nil {
		t.Error("NewContactFormConfiguration() SHOULD NOT return nil")
	}
	if cfg.Valid() {
		t.Error("NewContactFormConfiguration() SHOULD NOT return valid validator")
	}
}

func TestNewContactFormConfigurationIncorrectEnvVars(t *testing.T) {
	t.Setenv("XYZ_URL", "http://example.com")
	cfg := configuration.NewContactFormConfiguration()
	if cfg == nil {
		t.Error("NewContactFormConfiguration() SHOULD NOT return nil")
	}
	if cfg.Valid() {
		t.Error("NewContactFormConfiguration() SHOULD NOT return valid validator")
	}
}

func TestNewContactFormConfigurationCorrectEnvVars(t *testing.T) {
	t.Setenv("SENDGRID_API_KEY", "valid-api-key")
	cfg := configuration.NewContactFormConfiguration()
	if cfg == nil {
		t.Error("NewContactFormConfiguration() SHOULD NOT return nil")
	}
	if !cfg.Valid() {
		t.Error("NewContactFormConfiguration() SHOULD return valid validator")
	}
	if cfg != nil && cfg.SendGridApiKey != "valid-api-key" {
		t.Error("NewContactFormConfiguration() SHOULD return valid SendGridApiKey")
	}
}
