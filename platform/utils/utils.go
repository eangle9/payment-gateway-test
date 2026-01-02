package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	goerrors "errors"
	"fmt"
	"io"
	"math"
	mrand "math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

const (
	SmallLetters   = "abcdefghijklmnopqrstuvwxyz"
	CapitalLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits         = "0987654321"
)

func GenerateCustomRandomString(set string, length int) string {
	randString := make([]byte, length)
	_, _ = io.ReadAtLeast(rand.Reader, randString, length) //nolint:errcheck // since length = len(randString)

	for i := 0; i < len(randString); i++ {
		randString[i] = set[int(randString[i])%len(set)]
	}

	return string(randString)
}

// PrepareMultipartFormFile creates a multipart/form-data body.
//
// NOTE: be sure to close the writer before sending the request
func PrepareMultipartFormFile(filePath, fieldName string) (*bytes.Buffer, *multipart.Writer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, file.Name())
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, nil, err
	}

	return body, writer, nil
}

func CheckForNullUUID(errorMessage string) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(uuid.UUID)
		if !ok {
			return goerrors.New("faild to parse value to uuid")
		}
		nuuid := uuid.NullUUID{}
		if s == nuuid.UUID {
			return goerrors.New(errorMessage)
		}
		return nil
	}
}

func CheckForNullUUIDString(errorMessage string) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return goerrors.New("faild to parse value to uuid string")
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return goerrors.New("faild to parse value to uuid string")
		}
		nuuid := uuid.NullUUID{}
		if id == nuuid.UUID {
			return goerrors.New(errorMessage)
		}
		return nil
	}
}

// Calculate n percentage of x
func PercentageChange(n, x decimal.Decimal) (percentage decimal.Decimal) {
	if !x.IsZero() && !n.IsZero() {
		percentage = (n.Mul(x)).Div(decimal.NewFromInt(100))
		return percentage
	}
	return decimal.Zero
}

func Contains[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func GenerateTimestamp(currentTime time.Time) string {
	year := currentTime.Year()
	month := currentTime.Month()
	day := currentTime.Day()
	hour := currentTime.Hour()
	minute := currentTime.Minute()
	second := currentTime.Second()

	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", year, month, day, hour, minute, second)
}

func GenerateBase64EncodedString(input string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return encoded
}

func GenerateHash(input string, length int) (string, error) {

	// Create a new SHA-256 hash instance
	hash := sha256.New()

	// Write the input string to the hash instance
	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}

	// Get the raw SHA-256 hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	// Determine if the input contains alphabets, numbers, or both
	hasAlphabets := false
	hasNumbers := false
	for _, char := range input {
		if unicode.IsLetter(char) {
			hasAlphabets = true
		} else if unicode.IsDigit(char) {
			hasNumbers = true
		}
	}

	// Filter out non-alphabetic characters if input has only alphabets
	// Filter out non-numeric characters if input has only numbers
	if hasAlphabets && !hasNumbers {
		hashString = filterNonAlphabetic(hashString)
	} else if !hasAlphabets && hasNumbers {
		hashString = filterNonNumeric(hashString)
	}

	// Take the first 'length' characters from the hash string
	if len(hashString) > length {
		hashString = hashString[:length]
	}

	return hashString, nil
}

func filterNonAlphabetic(s string) string {
	result := ""
	for _, char := range s {
		if unicode.IsLetter(char) {
			result += string(char)
		}
	}
	return result
}

func filterNonNumeric(s string) string {
	result := ""
	for _, char := range s {
		if unicode.IsDigit(char) {
			result += string(char)
		}
	}
	return result
}

func ParseMixedDuration(input string) (time.Duration, error) {
	parts := strings.Fields(input)
	var days, hours, minutes, seconds int

	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if strings.Contains(part, "day") {
			value, err := strconv.Atoi(parts[i-1])
			if err != nil {
				return 0, err
			}
			days = value
		} else if strings.Contains(part, ":") {
			timeParts := strings.Split(part, ":")
			var timeValues []int
			for _, tp := range timeParts {
				value, err := strconv.Atoi(tp)
				if err != nil {
					return 0, err
				}
				timeValues = append(timeValues, value)
			}

			switch len(timeValues) {
			case 2: // Hours and Minutes
				hours = timeValues[0]
				minutes = timeValues[1]
			case 3: // Hours, Minutes, and Seconds
				hours = timeValues[0]
				minutes = timeValues[1]
				seconds = timeValues[2]
			}
		}
	}

	duration := time.Duration(days*24+hours)*time.Hour + time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second
	return duration, nil

}
func GenerateTime() string {
	t := time.Now()
	return fmt.Sprintf("%4d%02d%02v%02v%02v%02v", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func PublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	certificate, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(certificate)
}

func PrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyFile)
}

func MapJSONOrNull(t []byte) pgtype.JSON {
	if string(t) == "{}" || t == nil || string(t) == "null" || string(t) == "" {
		return pgtype.JSON{
			Status: pgtype.Null,
		}
	}

	return pgtype.JSON{
		Bytes:  t,
		Status: pgtype.Present,
	}
}

// Function to check if a string is a valid email format using ozzo-validation
func EmailValidator(email string) bool {
	err := validation.Validate(email, is.Email)
	return err == nil
}

// Custom validation rule to reject spaces
func NoSpaces(value interface{}) error {
	s, _ := value.(string)
	if len(s) == 0 {
		return nil
	}
	if containsSpace(s) {
		return fmt.Errorf("cannot contain spaces")
	}
	return nil
}

// Helper function to check if a string contains spaces
func containsSpace(s string) bool {
	for _, char := range s {
		if char == ' ' {
			return true
		}
	}
	return false
}

// GenerateCode generates a code with 5 digits followed by a capital letter.
func GenerateCode() string {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	digits := r.Intn(100000) // Generates a number between 0 and 99999
	letter := CapitalLetters[r.Intn(len(CapitalLetters))]
	return fmt.Sprintf("%05d%c", digits, letter)
}
func TimeToCronSpecWithTimezone(t time.Time, tz string) (string, error) {
	// Load the specified timezone
	location, err := time.LoadLocation(tz)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %v", err)
	}

	// Convert time to the specified timezone
	tInLocation := t.In(location)

	// Generate the cron spec string
	cronSpec := fmt.Sprintf("%d %d %d %d %d %d",
		tInLocation.Second(),       // Seconds
		tInLocation.Minute(),       // Minutes
		tInLocation.Hour(),         // Hours
		tInLocation.Day(),          // Day of month
		int(tInLocation.Month()),   // Month
		int(tInLocation.Weekday())) // Day of week (0 = Sunday)

	return fmt.Sprintf("CRON_TZ=%s %s", tz, cronSpec), nil
}
func TimeToCronSpecNoSeconds(t time.Time) string {
	return fmt.Sprintf("%d %d %d %d %d",
		t.Minute(),       // Minutes
		t.Hour(),         // Hours
		t.Day(),          // Day of month
		int(t.Month()),   // Month
		int(t.Weekday())) // Day of week (0 = Sunday)
}
func IsValidTimeZone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

// GetCurrentDateInTimezone returns the current date in the specified timezone.
// GetCurrentDateInTimezone returns the current date in the specified timezone.
func GetCurrentDateInTimezone(timeZoneName string, duration time.Duration) (time.Time, error) {
	// Load the specified timezone
	location, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %v", err)
	}

	// Get the current time in the specified timezone
	currentTime := time.Now().Add(duration).In(location)
	fmt.Printf("Current time in %s: %v\n", timeZoneName, currentTime)

	// Return the current date (year, month, day) in the specified timezone
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, location)
	fmt.Printf("Current date in %s: %v\n", timeZoneName, currentDate)

	return currentDate, nil
}

// GetOffsetFromUTC returns the offset from UTC in minutes for the specified timezone.
func GetOffsetFromUTC(timeZoneName string) (int, error) {
	// Load the specified timezone
	location, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return 0, fmt.Errorf("invalid timezone: %v", err)
	}

	// Get the current time in the specified timezone
	currentTime := time.Now().In(location)

	// Calculate the offset from UTC in minutes
	_, offset := currentTime.Zone()
	offsetMinutes := offset / 60

	return offsetMinutes, nil
}
func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}
func GenerateRandomNDigitNumber(n int) int {
	if n <= 0 {
		return 0
	}
	min := int(math.Pow10(n - 1))
	max := int(math.Pow10(n)) - 1
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func FormatToISO(t time.Time) time.Time {
	return t.UTC() // Ensures the time is in UTC
}
func GetWeekStart(t time.Time, tzName string) (time.Time, error) {
	// Load the timezone
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.Time{}, err
	}

	// Convert the time to the specified timezone
	localTime := t.In(loc)

	// Calculate Monday as the start of the week
	weekday := int(localTime.Weekday())
	if weekday == 0 {
		weekday = 7 // Adjust Sunday to be the last day of the week
	}
	weekStart := localTime.AddDate(0, 0, -weekday+1)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, loc)
	return weekStart, nil
}

func GetWeekEnd(t time.Time, tzName string) (time.Time, error) {
	// Get the start of the week first
	start, err := GetWeekStart(t, tzName)
	if err != nil {
		return time.Time{}, err
	}

	// Then, add 6 days to get to Sunday and set the time to the end of the day
	weekEnd := start.AddDate(0, 0, 6)
	weekEnd = time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 0, weekEnd.Location())
	return weekEnd, nil
}

// RandomUniqueNumbers generates up to 'count' unique random numbers in [0, max) range
// Handles edge cases where:
// - max <= 0: returns empty slice
// - count <= 0: returns empty slice
// - count > max: returns all numbers in range (shuffled)
func RandomUniqueNumbers(max, count int) []int {
	switch {
	case max <= 0 || count <= 0:
		return []int{}
	case max == 1:
		return []int{0}
	}

	// Create slice with all numbers in range
	numbers := make([]int, max)
	for i := range numbers {
		numbers[i] = i
	}

	// Shuffle the numbers
	mrand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	// Return requested count or maximum available
	if count > max {
		return numbers
	}
	return numbers[:count]
}
