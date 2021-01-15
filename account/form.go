package account

import "github.com/Alkemic/forms"

func newLoginForm() *forms.Form {
	return forms.New(map[string]*forms.Field{
		"username": {
			Attributes: map[string]interface{}{
				"class":       "form-control",
				"placeholder": "Username",
				"required":    "required",
			},
			Validators: []forms.Validator{
				&forms.Required{},
			},
		},
		"password": {
			Type: &forms.InputPassword{},
			Attributes: map[string]interface{}{
				"class":       "form-control",
				"placeholder": "Password",
				"required":    "required",
			},
			Validators: []forms.Validator{
				&forms.Required{},
			},
		},
	}, forms.Attributes{"id": "login-form"})
}
