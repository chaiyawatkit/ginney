package ginney

const (
	GinContextKey          = "ad1ad1b903a4711506a2bfd6a8fd9086d2aaee36fc267b9be847963b9412b95e"
	CorrelationIdHeaderKey = "X-Correlation-ID"
	ContentTypeHeaderKey   = "Content-Type"
	CensoredFieldText      = "[HIDDEN_FIELD]"
)

var (
	RequestBodyKeyCensoredList = []string{
		"password",
		"privatekey",
		"secretkey",
		"file",
		"phoneNumber",
	}
)
