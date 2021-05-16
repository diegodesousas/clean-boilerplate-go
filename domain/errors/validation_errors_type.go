package errors

type ValidationType string

var (
	Invalid     ValidationType = "invalid"
	Required    ValidationType = "required"
	GreaterThen ValidationType = "greater_then"
)

var Messages = map[ValidationType]string{
	Invalid:     "The field %s is invalid",
	Required:    "The field %s is required",
	GreaterThen: "The %s must be greater than %d",
}
