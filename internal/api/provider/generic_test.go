package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"EmailVerified", "email_verified"},
		{"PhoneVerified", "phone_verified"},
		{"FamilyName", "family_name"},
		{"GivenName", "given_name"},
		{"NickName", "nick_name"},
		{"MiddleName", "middle_name"},
		{"PreferredUsername", "preferred_username"},
		{"EmailPrimary", "email_primary"},
		{"Email", "email"},
		{"Name", "name"},
		{"Subject", "subject"},
		{"Issuer", "issuer"},
		{"Profile", "profile"},
		{"Picture", "picture"},
		{"Website", "website"},
		{"Gender", "gender"},
		{"Birthdate", "birthdate"},
		{"ZoneInfo", "zone_info"},
		{"Locale", "locale"},
		{"UpdatedAt", "updated_at"},
		{"Phone", "phone"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMappingField(t *testing.T) {
	t.Run("returns configured value when set", func(t *testing.T) {
		mapping := map[string]string{
			"Email":        "custom_email_field",
			"EmailVerified": "is_verified",
		}

		assert.Equal(t, "custom_email_field", getMappingField(mapping, "Email"))
		assert.Equal(t, "is_verified", getMappingField(mapping, "EmailVerified"))
	})

	t.Run("returns snake_case default when not configured", func(t *testing.T) {
		mapping := map[string]string{}

		assert.Equal(t, "email_verified", getMappingField(mapping, "EmailVerified"))
		assert.Equal(t, "phone_verified", getMappingField(mapping, "PhoneVerified"))
		assert.Equal(t, "family_name", getMappingField(mapping, "FamilyName"))
		assert.Equal(t, "given_name", getMappingField(mapping, "GivenName"))
		assert.Equal(t, "email", getMappingField(mapping, "Email"))
		assert.Equal(t, "name", getMappingField(mapping, "Name"))
	})

	t.Run("returns snake_case default when configured value is empty", func(t *testing.T) {
		mapping := map[string]string{
			"Email": "",
		}

		assert.Equal(t, "email", getMappingField(mapping, "Email"))
	})

	t.Run("returns configured value for some fields and default for others", func(t *testing.T) {
		mapping := map[string]string{
			"Email": "user.email_address",
			// EmailVerified not configured, should default to email_verified
			"Name": "user.fullname",
		}

		assert.Equal(t, "user.email_address", getMappingField(mapping, "Email"))
		assert.Equal(t, "email_verified", getMappingField(mapping, "EmailVerified"))
		assert.Equal(t, "user.fullname", getMappingField(mapping, "Name"))
		assert.Equal(t, "phone", getMappingField(mapping, "Phone"))
	})
}

func TestGetBooleanFieldByPath(t *testing.T) {
	t.Run("reads boolean field from map", func(t *testing.T) {
		obj := map[string]interface{}{
			"data": map[string]interface{}{
				"email_verified": true,
			},
		}

		result, err := getBooleanFieldByPath(obj, "data.email_verified", false)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("returns fallback when field not found", func(t *testing.T) {
		obj := map[string]interface{}{
			"data": map[string]interface{}{},
		}

		result, err := getBooleanFieldByPath(obj, "data.missing_field", true)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("returns fallback when path is snake_case default for EmailVerified", func(t *testing.T) {
		obj := map[string]interface{}{
			"email_verified": true,
		}

		mapping := map[string]string{}
		path := getMappingField(mapping, "EmailVerified")

		result, err := getBooleanFieldByPath(obj, path, false)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("uses custom mapping path when configured", func(t *testing.T) {
		obj := map[string]interface{}{
			"custom": map[string]interface{}{
				"verified_status": true,
			},
		}

		mapping := map[string]string{
			"EmailVerified": "custom.verified_status",
		}
		path := getMappingField(mapping, "EmailVerified")

		result, err := getBooleanFieldByPath(obj, path, false)
		require.NoError(t, err)
		assert.True(t, result)
	})
}

func TestGetStringFieldByPath(t *testing.T) {
	t.Run("reads string field from map", func(t *testing.T) {
		obj := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John Doe",
			},
		}

		result, err := getStringFieldByPath(obj, "user.name", "")
		require.NoError(t, err)
		assert.Equal(t, "John Doe", result)
	})

	t.Run("returns fallback when field not found", func(t *testing.T) {
		obj := map[string]interface{}{
			"user": map[string]interface{}{},
		}

		result, err := getStringFieldByPath(obj, "user.missing", "default_value")
		require.NoError(t, err)
		assert.Equal(t, "default_value", result)
	})

	t.Run("returns fallback when path is snake_case default for Name", func(t *testing.T) {
		obj := map[string]interface{}{
			"name": "Jane Doe",
		}

		mapping := map[string]string{}
		path := getMappingField(mapping, "Name")

		result, err := getStringFieldByPath(obj, path, "")
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", result)
	})

	t.Run("uses custom mapping path when configured", func(t *testing.T) {
		obj := map[string]interface{}{
			"profile": map[string]interface{}{
				"full_name": "Jane Doe",
			},
		}

		mapping := map[string]string{
			"Name": "profile.full_name",
		}
		path := getMappingField(mapping, "Name")

		result, err := getStringFieldByPath(obj, path, "")
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", result)
	})

	t.Run("converts int to string", func(t *testing.T) {
		obj := map[string]interface{}{
			"age": 25,
		}

		result, err := getStringFieldByPath(obj, "age", "")
		require.NoError(t, err)
		assert.Equal(t, "25", result)
	})

	t.Run("converts float64 to string", func(t *testing.T) {
		obj := map[string]interface{}{
			"score": 95.7,
		}

		result, err := getStringFieldByPath(obj, "score", "")
		require.NoError(t, err)
		assert.Equal(t, "96", result)
	})

	t.Run("returns empty string for nil value", func(t *testing.T) {
		obj := map[string]interface{}{
			"nullable": nil,
		}

		result, err := getStringFieldByPath(obj, "nullable", "")
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
}
