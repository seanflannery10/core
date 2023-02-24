package validator

type Validator struct {
	Errors map[string]string `json:",omitempty"`
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) != 0
}

func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = map[string]string{}
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
