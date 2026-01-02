package constant

type ContextKey any
type Role string

const (
	Admin    Role = "ADMIN"
	HR       Role = "HR"
	Employee Role = "EMPLOYEE"
)

type Status string

const (
	Active   Status = "ACTIVE"
	Inactive Status = "INACTIVE"
	Pending  Status = "PENDING"
	Deleted  Status = "DELETED"
	Yes      Status = "YES"
	No       Status = "NO"
)
const (
	AuthorizationHeaderkey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
	REQUESTTIME             = "2006-01-02 15:04:05"
)

type TokenProvider string

const (
	Apple    TokenProvider = "APPLE"
	Google   TokenProvider = "GOOGLE"
	Facebook TokenProvider = "FACEBOOK"
	Normal   TokenProvider = "NORMAL"
)

type TokenType string

const (
	AccessToken           TokenType = "ACCESS_TOKEN"
	RefreshToken          TokenType = "REFRESH_TOKEN"
	ResetPassword         TokenType = "RESET_PASSWORD"
	VerificationTokenType TokenType = "VERIFICATION_TOKEN"
	InviteTokenType       TokenType = "INVITE_TOKEN"
	SecretToken           TokenType = "SECRET_TOKEN"
)

type PaymentType string

const (
	PaymentTypeOnetime = "ONETIME"
)

type Currency string

const (
	CurrencyETB Currency = "ETB"
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
	CurrencyGBP Currency = "GBP"
)
