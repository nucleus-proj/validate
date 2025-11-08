package fluent_validator

import (
	"encoding/base64"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// String rules
func NonEmpty(s string) ValidatorFunc {
	return func() ValidationResult {
		if s == "" {
			return Fail("must not be empty")
		}
		return Success()
	}
}

func MinLen(s string, n int) ValidatorFunc {
	return func() ValidationResult {
		if len(s) < n {
			return Fail("too short: min " + strconv.Itoa(n))
		}
		return Success()
	}
}

func MaxLen(s string, n int) ValidatorFunc {
	return func() ValidationResult {
		if len(s) > n {
			return Fail("too long: max " + strconv.Itoa(n))
		}
		return Success()
	}
}

func LenBetween(s string, min, max int) ValidatorFunc {
	return func() ValidationResult {
		l := len(s)
		if l < min || l > max {
			return Fail("length must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max))
		}
		return Success()
	}
}

func Matches(s string, re *regexp.Regexp) ValidatorFunc {
	return func() ValidationResult {
		if !re.MatchString(s) {
			return Fail("must match pattern")
		}
		return Success()
	}
}

func OneOf(s string, allowed []string, caseSensitive bool) ValidatorFunc {
	return func() ValidationResult {
		if !caseSensitive {
			s = strings.ToLower(s)
		}
		for _, a := range allowed {
			if !caseSensitive {
				a = strings.ToLower(a)
			}
			if s == a {
				return Success()
			}
		}
		return Fail("must be one of: " + strings.Join(allowed, ", "))
	}
}

// Number rules
func IntMin(v, min int) ValidatorFunc {
	return func() ValidationResult {
		if v < min {
			return Fail("must be >= " + strconv.Itoa(min))
		}
		return Success()
	}
}
func IntMax(v, max int) ValidatorFunc {
	return func() ValidationResult {
		if v > max {
			return Fail("must be <= " + strconv.Itoa(max))
		}
		return Success()
	}
}
func IntBetween(v, min, max int) ValidatorFunc {
	return func() ValidationResult {
		if v < min || v > max {
			return Fail("must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max))
		}
		return Success()
	}
}
func IntNonZero(v int) ValidatorFunc {
	return func() ValidationResult {
		if v == 0 {
			return Fail("must not be zero")
		}
		return Success()
	}
}

func FloatMin(v, min float64) ValidatorFunc {
	return func() ValidationResult {
		if v < min {
			return Fail("must be >= " + trimFloatZeros(min))
		}
		return Success()
	}
}
func FloatMax(v, max float64) ValidatorFunc {
	return func() ValidationResult {
		if v > max {
			return Fail("must be <= " + trimFloatZeros(max))
		}
		return Success()
	}
}
func FloatBetween(v, min, max float64) ValidatorFunc {
	return func() ValidationResult {
		if v < min || v > max {
			return Fail("must be between " + trimFloatZeros(min) + " and " + trimFloatZeros(max))
		}
		return Success()
	}
}
func FloatNonZero(v float64) ValidatorFunc {
	return func() ValidationResult {
		if v == 0 {
			return Fail("must not be zero")
		}
		return Success()
	}
}

// Number extras
func IntPositive(v int) ValidatorFunc {
	return func() ValidationResult {
		if v <= 0 {
			return Fail("must be > 0")
		}
		return Success()
	}
}
func IntNonNegative(v int) ValidatorFunc {
	return func() ValidationResult {
		if v < 0 {
			return Fail("must be >= 0")
		}
		return Success()
	}
}
func IntGreaterThan(v, min int) ValidatorFunc {
	return func() ValidationResult {
		if v <= min {
			return Fail("must be > " + strconv.Itoa(min))
		}
		return Success()
	}
}
func IntLessThan(v, max int) ValidatorFunc {
	return func() ValidationResult {
		if v >= max {
			return Fail("must be < " + strconv.Itoa(max))
		}
		return Success()
	}
}
func IntMultipleOf(v, m int) ValidatorFunc {
	return func() ValidationResult {
		if m == 0 || v%m != 0 {
			return Fail("must be a multiple of " + strconv.Itoa(m))
		}
		return Success()
	}
}

func FloatGreaterThan(v, min float64) ValidatorFunc {
	return func() ValidationResult {
		if !(v > min) {
			return Fail("must be > " + trimFloatZeros(min))
		}
		return Success()
	}
}
func FloatLessThan(v, max float64) ValidatorFunc {
	return func() ValidationResult {
		if !(v < max) {
			return Fail("must be < " + trimFloatZeros(max))
		}
		return Success()
	}
}
func FloatMultipleOf(v, m float64) ValidatorFunc {
	return func() ValidationResult {
		if m == 0 {
			return Fail("must be a multiple of 0 is undefined")
		}
		q := v / m
		qi := float64(int64(q))
		r := v - qi*m
		if r < 0 {
			r = -r
		}
		if r > 1e-9 {
			return Fail("must be a multiple of " + trimFloatZeros(m))
		}
		return Success()
	}
}

// Time rules
func TimeNotZero(t time.Time) ValidatorFunc {
	return func() ValidationResult {
		if t.IsZero() {
			return Fail("must not be zero time")
		}
		return Success()
	}
}
func TimeBefore(t, cutoff time.Time) ValidatorFunc {
	return func() ValidationResult {
		if !t.Before(cutoff) {
			return Fail("must be before cutoff")
		}
		return Success()
	}
}
func TimeAfter(t, cutoff time.Time) ValidatorFunc {
	return func() ValidationResult {
		if !t.After(cutoff) {
			return Fail("must be after cutoff")
		}
		return Success()
	}
}
func TimeBetween(t, start, end time.Time) ValidatorFunc {
	return func() ValidationResult {
		if t.Before(start) || t.After(end) {
			return Fail("must be between start and end")
		}
		return Success()
	}
}

// Time extras
func InPast(t time.Time) ValidatorFunc {
	return func() ValidationResult {
		if !t.Before(time.Now()) {
			return Fail("must be in the past")
		}
		return Success()
	}
}
func InFuture(t time.Time) ValidatorFunc {
	return func() ValidationResult {
		if !t.After(time.Now()) {
			return Fail("must be in the future")
		}
		return Success()
	}
}
func IsWeekday(t time.Time) ValidatorFunc {
	return func() ValidationResult {
		wd := t.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			return Fail("must be a weekday")
		}
		return Success()
	}
}
func IsWeekend(t time.Time) ValidatorFunc {
	return func() ValidationResult {
		wd := t.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			return Fail("must be a weekend day")
		}
		return Success()
	}
}

// Duration rules
func DurationMin(d, min time.Duration) ValidatorFunc {
	return func() ValidationResult {
		if d < min {
			return Fail("duration too small: min " + min.String())
		}
		return Success()
	}
}
func DurationMax(d, max time.Duration) ValidatorFunc {
	return func() ValidationResult {
		if d > max {
			return Fail("duration too large: max " + max.String())
		}
		return Success()
	}
}

// Collection rules (length-based via explicit length parameter)
func NotEmptyLen(n int) ValidatorFunc {
	return func() ValidationResult {
		if n == 0 {
			return Fail("must not be empty")
		}
		return Success()
	}
}
func LenMin(n, min int) ValidatorFunc {
	return func() ValidationResult {
		if n < min {
			return Fail("size too small: min " + strconv.Itoa(min))
		}
		return Success()
	}
}
func LenMax(n, max int) ValidatorFunc {
	return func() ValidationResult {
		if n > max {
			return Fail("size too large: max " + strconv.Itoa(max))
		}
		return Success()
	}
}
func LenBetweenSize(n, min, max int) ValidatorFunc {
	return func() ValidationResult {
		if n < min || n > max {
			return Fail("size must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max))
		}
		return Success()
	}
}

func ContainsString(list []string, elem string) ValidatorFunc {
	return func() ValidationResult {
		for _, v := range list {
			if v == elem {
				return Success()
			}
		}
		return Fail("must contain " + elem)
	}
}

func UniqueStrings(list []string) ValidatorFunc {
	return func() ValidationResult {
		seen := make(map[string]struct{}, len(list))
		for _, v := range list {
			if _, ok := seen[v]; ok {
				return Fail("must be unique")
			}
			seen[v] = struct{}{}
		}
		return Success()
	}
}

// Email and phone
var reEmailLight = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var reE164 = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func EmailValid(s string) ValidatorFunc {
	return func() ValidationResult {
		if s == "" {
			return Fail("must not be empty")
		}
		if !reEmailLight.MatchString(s) {
			return Fail("invalid email")
		}
		return Success()
	}
}

func PhoneE164(s string) ValidatorFunc {
	return func() ValidationResult {
		if !reE164.MatchString(s) {
			return Fail("invalid phone (use E.164, e.g. +15551234567)")
		}
		return Success()
	}
}

// PhoneWithCountryCode validates E.164 phone number with a required country code prefix, e.g., "+251".
func PhoneWithCountryCode(s string, countryCode string) ValidatorFunc {
	return func() ValidationResult {
		if !strings.HasPrefix(s, countryCode) {
			return Fail("invalid phone: must start with " + countryCode)
		}
		if !reE164.MatchString(s) {
			return Fail("invalid phone (use E.164, e.g. +15551234567)")
		}
		return Success()
	}
}

// Additional string classifiers
func HasPrefix(s, prefix string) ValidatorFunc {
	return func() ValidationResult {
		if !strings.HasPrefix(s, prefix) {
			return Fail("must start with " + prefix)
		}
		return Success()
	}
}
func HasSuffix(s, suffix string) ValidatorFunc {
	return func() ValidationResult {
		if !strings.HasSuffix(s, suffix) {
			return Fail("must end with " + suffix)
		}
		return Success()
	}
}
func Contains(s, substr string) ValidatorFunc {
	return func() ValidationResult {
		if !strings.Contains(s, substr) {
			return Fail("must contain " + substr)
		}
		return Success()
	}
}
func Trimmed(s string) ValidatorFunc {
	return func() ValidationResult {
		if strings.TrimSpace(s) != s {
			return Fail("must not have leading/trailing spaces")
		}
		return Success()
	}
}
func IsAlpha(s string) ValidatorFunc {
	return func() ValidationResult {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return Fail("must contain only letters")
			}
		}
		return Success()
	}
}
func IsNumeric(s string) ValidatorFunc {
	return func() ValidationResult {
		if s == "" {
			return Fail("must be numeric")
		}
		for _, r := range s {
			if !unicode.IsDigit(r) {
				return Fail("must be numeric")
			}
		}
		return Success()
	}
}
func IsAlnum(s string) ValidatorFunc {
	return func() ValidationResult {
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return Fail("must be alphanumeric")
			}
		}
		return Success()
	}
}

var reHex = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func IsHex(s string) ValidatorFunc {
	return func() ValidationResult {
		if !reHex.MatchString(s) {
			return Fail("must be hex")
		}
		return Success()
	}
}
func IsBase64(s string) ValidatorFunc {
	return func() ValidationResult {
		if _, err := base64.StdEncoding.DecodeString(s); err != nil {
			return Fail("must be base64")
		}
		return Success()
	}
}

var reSlug = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func IsSlug(s string) ValidatorFunc {
	return func() ValidationResult {
		if !reSlug.MatchString(s) {
			return Fail("must be a slug")
		}
		return Success()
	}
}

var reUUIDv4 = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func IsUUIDv4(s string) ValidatorFunc {
	return func() ValidationResult {
		if !reUUIDv4.MatchString(s) {
			return Fail("must be UUID v4")
		}
		return Success()
	}
}

var reULID = regexp.MustCompile(`^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)

func IsULID(s string) ValidatorFunc {
	return func() ValidationResult {
		if !reULID.MatchString(s) {
			return Fail("must be ULID")
		}
		return Success()
	}
}

// URL/Hostname/IP
func IsURL(s string) ValidatorFunc {
	return func() ValidationResult {
		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return Fail("must be URL")
		}
		return Success()
	}
}

var reHostname = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)(?:\.(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?))*$`)

func IsHostname(s string) ValidatorFunc {
	return func() ValidationResult {
		if len(s) > 253 || !reHostname.MatchString(s) {
			return Fail("must be hostname")
		}
		return Success()
	}
}
func IsIP(s string) ValidatorFunc {
	return func() ValidationResult {
		if net.ParseIP(s) == nil {
			return Fail("must be IP")
		}
		return Success()
	}
}
func IsIPv4(s string) ValidatorFunc {
	return func() ValidationResult {
		ip := net.ParseIP(s)
		if ip == nil || ip.To4() == nil {
			return Fail("must be IPv4")
		}
		return Success()
	}
}
func IsIPv6(s string) ValidatorFunc {
	return func() ValidationResult {
		ip := net.ParseIP(s)
		if ip == nil || ip.To4() != nil {
			return Fail("must be IPv6")
		}
		return Success()
	}
}
func IsCIDR(s string) ValidatorFunc {
	return func() ValidationResult {
		if _, _, err := net.ParseCIDR(s); err != nil {
			return Fail("must be CIDR")
		}
		return Success()
	}
}

// Email domain policies (simple split)
func EmailDomainAllowlist(s string, allowed []string) ValidatorFunc {
	return func() ValidationResult {
		at := strings.LastIndexByte(s, '@')
		if at == -1 {
			return Fail("invalid email")
		}
		dom := strings.ToLower(s[at+1:])
		for _, d := range allowed {
			if dom == strings.ToLower(d) {
				return Success()
			}
		}
		return Fail("email domain not allowed")
	}
}
func EmailDomainBlocklist(s string, blocked []string) ValidatorFunc {
	return func() ValidationResult {
		at := strings.LastIndexByte(s, '@')
		if at == -1 {
			return Fail("invalid email")
		}
		dom := strings.ToLower(s[at+1:])
		for _, d := range blocked {
			if dom == strings.ToLower(d) {
				return Fail("email domain blocked")
			}
		}
		return Success()
	}
}

// Luhn checksum (e.g., credit card numbers); input should be digits only (spaces allowed)
func LuhnValid(s string) ValidatorFunc {
	return func() ValidationResult {
		sum := 0
		alt := false
		digits := 0
		for i := len(s) - 1; i >= 0; i-- {
			ch := s[i]
			if ch == ' ' {
				continue
			}
			if ch < '0' || ch > '9' {
				return Fail("must be numeric")
			}
			d := int(ch - '0')
			if alt {
				d *= 2
				if d > 9 {
					d -= 9
				}
			}
			sum += d
			alt = !alt
			digits++
		}
		if digits == 0 || sum%10 != 0 {
			return Fail("invalid luhn")
		}
		return Success()
	}
}

func trimFloatZeros(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	// trim trailing zeros and optional dot
	i := len(s)
	for i > 0 && s[i-1] == '0' {
		i--
	}
	if i > 0 && s[i-1] == '.' {
		i--
	}
	return s[:i]
}
