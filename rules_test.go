package fluent_validator

import (
	"encoding/base64"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestStringRules(t *testing.T) {
	t.Parallel()
	re := regexp.MustCompile(`^[a-z]+$`)
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"NonEmpty ok", NonEmpty("x"), true, nil},
		{"NonEmpty fail", NonEmpty(""), false, []string{"must not be empty"}},
		{"MinLen ok", MinLen("abcd", 3), true, nil},
		{"MinLen fail", MinLen("ab", 3), false, []string{"too short: min 3"}},
		{"MaxLen ok", MaxLen("ab", 3), true, nil},
		{"MaxLen fail", MaxLen("abcd", 3), false, []string{"too long: max 3"}},
		{"LenBetween ok", LenBetween("abc", 2, 3), true, nil},
		{"LenBetween fail", LenBetween("a", 2, 3), false, []string{"length must be between 2 and 3"}},
		{"Matches ok", Matches("abc", re), true, nil},
		{"Matches fail", Matches("ab1", re), false, []string{"must match pattern"}},
		{"OneOf ok", OneOf("b", []string{"a", "b"}, true), true, nil},
		{"OneOf fail", OneOf("c", []string{"a", "b"}, true), false, []string{"must be one of: a, b"}},
		{"OneOf case-insensitive ok", OneOf("B", []string{"a", "b"}, false), true, nil},
		{"HasPrefix ok", HasPrefix("foobar", "foo"), true, nil},
		{"HasPrefix fail", HasPrefix("bar", "foo"), false, []string{"must start with foo"}},
		{"HasSuffix ok", HasSuffix("foobar", "bar"), true, nil},
		{"HasSuffix fail", HasSuffix("foo", "bar"), false, []string{"must end with bar"}},
		{"Contains ok", Contains("hello world", "world"), true, nil},
		{"Contains fail", Contains("hello", "world"), false, []string{"must contain world"}},
		{"Trimmed ok", Trimmed("abc"), true, nil},
		{"Trimmed fail", Trimmed(" abc "), false, []string{"must not have leading/trailing spaces"}},
		{"IsAlpha ok", IsAlpha("abcXYZ"), true, nil},
		{"IsAlpha fail", IsAlpha("abc123"), false, []string{"must contain only letters"}},
		{"IsNumeric ok", IsNumeric("123"), true, nil},
		{"IsNumeric fail", IsNumeric("12a"), false, []string{"must be numeric"}},
		{"IsAlnum ok", IsAlnum("abc123"), true, nil},
		{"IsAlnum fail", IsAlnum("abc-123"), false, []string{"must be alphanumeric"}},
		{"IsHex ok", IsHex("0A1b"), true, nil},
		{"IsHex fail", IsHex("g001"), false, []string{"must be hex"}},
		{"IsBase64 ok", IsBase64(base64.StdEncoding.EncodeToString([]byte("hi"))), true, nil},
		{"IsBase64 fail", IsBase64("not-base64"), false, []string{"must be base64"}},
		{"IsSlug ok", IsSlug("hello-world"), true, nil},
		{"IsSlug fail", IsSlug("Hello World"), false, []string{"must be a slug"}},
		{"IsUUIDv4 ok", IsUUIDv4("550e8400-e29b-41d4-a716-446655440000"), true, nil},
		{"IsUUIDv4 fail", IsUUIDv4("550e8400-e29b-21d4-a716-446655440000"), false, []string{"must be UUID v4"}},
		{"IsULID ok", IsULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"), true, nil},
		{"IsULID fail", IsULID("Z1ARZ3NDEKTSV4RRFFQ69G5FAV"), false, []string{"must be ULID"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestNumberRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"IntMin ok", IntMin(5, 3), true, nil},
		{"IntMin fail", IntMin(2, 3), false, []string{"must be >= 3"}},
		{"IntMax ok", IntMax(3, 5), true, nil},
		{"IntMax fail", IntMax(6, 5), false, []string{"must be <= 5"}},
		{"IntBetween ok", IntBetween(4, 3, 5), true, nil},
		{"IntBetween fail", IntBetween(2, 3, 5), false, []string{"must be between 3 and 5"}},
		{"IntNonZero ok", IntNonZero(1), true, nil},
		{"IntNonZero fail", IntNonZero(0), false, []string{"must not be zero"}},
		{"IntPositive ok", IntPositive(1), true, nil},
		{"IntPositive fail", IntPositive(0), false, []string{"must be > 0"}},
		{"IntNonNegative ok", IntNonNegative(0), true, nil},
		{"IntNonNegative fail", IntNonNegative(-1), false, []string{"must be >= 0"}},
		{"IntGreaterThan ok", IntGreaterThan(6, 5), true, nil},
		{"IntGreaterThan fail", IntGreaterThan(5, 5), false, []string{"must be > 5"}},
		{"IntLessThan ok", IntLessThan(4, 5), true, nil},
		{"IntLessThan fail", IntLessThan(5, 5), false, []string{"must be < 5"}},
		{"IntMultipleOf ok", IntMultipleOf(10, 5), true, nil},
		{"IntMultipleOf fail", IntMultipleOf(11, 5), false, []string{"must be a multiple of 5"}},

		{"FloatMin ok", FloatMin(3.2, 3.1), true, nil},
		{"FloatMin fail", FloatMin(3.0, 3.1), false, []string{"must be >= 3.1"}},
		{"FloatMax ok", FloatMax(3.2, 3.3), true, nil},
		{"FloatMax fail", FloatMax(3.4, 3.3), false, []string{"must be <= 3.3"}},
		{"FloatBetween ok", FloatBetween(3.2, 3.1, 3.3), true, nil},
		{"FloatBetween fail", FloatBetween(3.4, 3.1, 3.3), false, []string{"must be between 3.1 and 3.3"}},
		{"FloatNonZero ok", FloatNonZero(0.1), true, nil},
		{"FloatNonZero fail", FloatNonZero(0.0), false, []string{"must not be zero"}},
		{"FloatGreaterThan ok", FloatGreaterThan(3.2, 3.1), true, nil},
		{"FloatGreaterThan fail", FloatGreaterThan(3.1, 3.1), false, []string{"must be > 3.1"}},
		{"FloatLessThan ok", FloatLessThan(3.2, 3.3), true, nil},
		{"FloatLessThan fail", FloatLessThan(3.3, 3.3), false, []string{"must be < 3.3"}},
		{"FloatMultipleOf ok", FloatMultipleOf(10.0, 2.5), true, nil},
		{"FloatMultipleOf fail", FloatMultipleOf(10.1, 2.5), false, []string{"must be a multiple of 2.5"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestTimeRules(t *testing.T) {
	t.Parallel()
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"TimeNotZero ok", TimeNotZero(now), true, nil},
		{"TimeNotZero fail", TimeNotZero(time.Time{}), false, []string{"must not be zero time"}},
		{"TimeBefore ok", TimeBefore(past, future), true, nil},
		{"TimeBefore fail", TimeBefore(future, past), false, []string{"must be before cutoff"}},
		{"TimeAfter ok", TimeAfter(future, past), true, nil},
		{"TimeAfter fail", TimeAfter(past, future), false, []string{"must be after cutoff"}},
		{"TimeBetween ok", TimeBetween(now, past, future), true, nil},
		{"TimeBetween fail", TimeBetween(past, now, future), false, []string{"must be between start and end"}},
		{"InPast ok", InPast(past), true, nil},
		{"InPast fail", InPast(future), false, []string{"must be in the past"}},
		{"InFuture ok", InFuture(future), true, nil},
		{"InFuture fail", InFuture(past), false, []string{"must be in the future"}},
		{"IsWeekday ok", IsWeekday(time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)), true, nil},
		{"IsWeekday fail", IsWeekday(time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC)), false, []string{"must be a weekday"}},
		{"IsWeekend ok", IsWeekend(time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC)), true, nil},
		{"IsWeekend fail", IsWeekend(time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)), false, []string{"must be a weekend day"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestCollectionRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"NotEmptyLen ok", NotEmptyLen(1), true, nil},
		{"NotEmptyLen fail", NotEmptyLen(0), false, []string{"must not be empty"}},
		{"LenMin ok", LenMin(3, 2), true, nil},
		{"LenMin fail", LenMin(1, 2), false, []string{"size too small: min 2"}},
		{"LenMax ok", LenMax(2, 3), true, nil},
		{"LenMax fail", LenMax(4, 3), false, []string{"size too large: max 3"}},
		{"LenBetweenSize ok", LenBetweenSize(3, 2, 4), true, nil},
		{"LenBetweenSize fail", LenBetweenSize(1, 2, 4), false, []string{"size must be between 2 and 4"}},
		{"ContainsString ok", ContainsString([]string{"a", "b"}, "b"), true, nil},
		{"ContainsString fail", ContainsString([]string{"a", "b"}, "c"), false, []string{"must contain c"}},
		{"UniqueStrings ok", UniqueStrings([]string{"a", "b"}), true, nil},
		{"UniqueStrings fail", UniqueStrings([]string{"a", "b", "a"}), false, []string{"must be unique"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestEmailPhoneRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"EmailValid ok", EmailValid("user@example.com"), true, nil},
		{"EmailValid empty", EmailValid(""), false, []string{"must not be empty"}},
		{"EmailValid bad", EmailValid("user@"), false, []string{"invalid email"}},
		{"PhoneE164 ok", PhoneE164("+15551234567"), true, nil},
		{"PhoneE164 bad", PhoneE164("5551234567"), false, []string{"invalid phone (use E.164, e.g. +15551234567)"}},
		{"PhoneWithCountryCode ok", PhoneWithCountryCode("+251912345678", "+251"), true, nil},
		{"PhoneWithCountryCode wrong prefix", PhoneWithCountryCode("+15551234567", "+251"), false, []string{"invalid phone: must start with +251"}},
		{"PhoneWithCountryCode malformed", PhoneWithCountryCode("251912345678", "+251"), false, []string{"invalid phone: must start with +251"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestNetAndIdRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"IsURL ok", IsURL("https://example.com/path"), true, nil},
		{"IsURL fail", IsURL("not a url"), false, []string{"must be URL"}},
		{"IsHostname ok", IsHostname("example.com"), true, nil},
		{"IsHostname fail", IsHostname("-bad-.com"), false, []string{"must be hostname"}},
		{"IsIP v4 ok", IsIPv4("192.168.1.1"), true, nil},
		{"IsIP v4 fail", IsIPv4("abcd"), false, []string{"must be IPv4"}},
		{"IsIP v6 ok", IsIPv6("2001:db8::1"), true, nil},
		{"IsIP v6 fail", IsIPv6("192.168.1.1"), false, []string{"must be IPv6"}},
		{"IsCIDR ok", IsCIDR("10.0.0.0/8"), true, nil},
		{"IsCIDR fail", IsCIDR("10.0.0.0"), false, []string{"must be CIDR"}},
		{"EmailDomainAllowlist ok", EmailDomainAllowlist("a@ex.com", []string{"ex.com"}), true, nil},
		{"EmailDomainAllowlist fail", EmailDomainAllowlist("a@ex.com", []string{"other.com"}), false, []string{"email domain not allowed"}},
		{"EmailDomainBlocklist ok", EmailDomainBlocklist("a@ex.com", []string{"other.com"}), true, nil},
		{"EmailDomainBlocklist fail", EmailDomainBlocklist("a@ex.com", []string{"ex.com"}), false, []string{"email domain blocked"}},
		{"LuhnValid ok", LuhnValid("4539 1488 0343 6467"), true, nil},
		{"LuhnValid fail", LuhnValid("4539 1488 0343 6468"), false, []string{"invalid luhn"}},
	}
	_ = net.IPv4(0, 0, 0, 0) // keep net import
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}

func TestDurationRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		v         Validator
		wantValid bool
		wantMsg   []string
	}{
		{"DurationMin ok", DurationMin(5*time.Second, 3*time.Second), true, nil},
		{"DurationMin fail", DurationMin(2*time.Second, 3*time.Second), false, []string{"duration too small: min 3s"}},
		{"DurationMax ok", DurationMax(2*time.Second, 3*time.Second), true, nil},
		{"DurationMax fail", DurationMax(4*time.Second, 3*time.Second), false, []string{"duration too large: max 3s"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.v.Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("valid=%v want %v", res.IsValid, tc.wantValid)
			}
			if tc.wantMsg != nil && !reflect.DeepEqual(res.Message, tc.wantMsg) {
				t.Fatalf("msg=%v want %v", res.Message, tc.wantMsg)
			}
		})
	}
}
