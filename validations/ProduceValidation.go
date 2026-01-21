package validations

type RequestQueueValidation struct {
	TargetUrl string                   `json:"targetUrl" validate:"required,max=128"`
	Method    string                   `json:"httpMethod" validate:"required,max=16"`
	Headers   map[string]interface{}   `json:"headers"`
	Body      []map[string]interface{} `json:"body" validate:"required"`
}

type TestHitHeaderValidation struct {
	Authorization string `json:"Authorization" validate:"required"`
	ContentType   string `json:"Content-Type" validate:"required"`
}

type TestHitRequestValidation struct {
	Ping string `json:"ping" validate:"required"`
}
