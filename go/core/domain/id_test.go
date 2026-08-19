package domain

import "testing"

func TestNewUUIDReturnsCanonicalVersion4UUID(t *testing.T) {
	value, err := NewUUID()
	if err != nil {
		t.Fatalf("NewUUID returned error: %v", err)
	}
	if err := ValidateUUID(value); err != nil {
		t.Fatalf("ValidateUUID(%q) returned error: %v", value, err)
	}
	if value[14] != '4' {
		t.Fatalf("UUID version = %q, want 4", value[14])
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "version 4", value: "0f8fad5b-d9cb-469f-a165-70867728950e", valid: true},
		{name: "version 7", value: "01890f7a-3cda-7cc1-8f7d-8f7d8f7d8f7d", valid: true},
		{name: "uppercase", value: "0F8FAD5B-D9CB-469F-A165-70867728950E"},
		{name: "wrong separators", value: "0f8fad5b_d9cb_469f_a165_70867728950e"},
		{name: "extra separator", value: "0f8fad5b-d9cb-469f-a165-7086772-950e"},
		{name: "wrong variant", value: "0f8fad5b-d9cb-469f-7165-70867728950e"},
		{name: "missing version", value: "0f8fad5b-d9cb-069f-a165-70867728950e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateUUID(test.value)
			if test.valid && err != nil {
				t.Fatalf("ValidateUUID returned error: %v", err)
			}
			if !test.valid && err == nil {
				t.Fatal("ValidateUUID succeeded")
			}
		})
	}
}
