package templateutil_test

import (
	"testing"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

func TestNormalizeNames(t *testing.T) {
	cases := []struct {
		input       string
		pascal      string
		snake       string
		kebab       string
		pluralKebab string
	}{
		{"product", "Product", "product", "product", "products"},
		{"order-item", "OrderItem", "order_item", "order-item", "order-items"},
		{"OrderItem", "OrderItem", "order_item", "order-item", "order-items"},
		{"order_item", "OrderItem", "order_item", "order-item", "order-items"},
		{"userProfile", "UserProfile", "user_profile", "user-profile", "user-profiles"},
	}

	for _, tc := range cases {
		n := templateutil.NormalizeName(tc.input)
		if n.Pascal != tc.pascal {
			t.Errorf("input=%q: Pascal got %q want %q", tc.input, n.Pascal, tc.pascal)
		}
		if n.Snake != tc.snake {
			t.Errorf("input=%q: Snake got %q want %q", tc.input, n.Snake, tc.snake)
		}
		if n.Kebab != tc.kebab {
			t.Errorf("input=%q: Kebab got %q want %q", tc.input, n.Kebab, tc.kebab)
		}
		if n.PluralKebab != tc.pluralKebab {
			t.Errorf("input=%q: PluralKebab got %q want %q", tc.input, n.PluralKebab, tc.pluralKebab)
		}
	}
}
