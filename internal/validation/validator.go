package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError — equivalente ao ValidationException do Laravel
type ValidationError struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validação falhou: %v", e.Errors)
}

func (e *ValidationError) Add(field, msg string) {
	e.Errors[field] = append(e.Errors[field], msg)
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validator — equivalente ao Validator::make() do Laravel
type Validator struct {
	data   map[string]interface{}
	errors *ValidationError
}

func New(data map[string]interface{}) *Validator {
	return &Validator{
		data:   data,
		errors: &ValidationError{Errors: make(map[string][]string)},
	}
}

// Validate — executa as regras e retorna erro se inválido
// Equivalente ao $request->validate([...]) do Laravel
func (v *Validator) Validate(rules map[string][]string) error {
	for field, fieldRules := range rules {
		value, exists := v.data[field]
		for _, rule := range fieldRules {
			v.applyRule(field, value, exists, rule)
		}
	}
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}

func (v *Validator) applyRule(field string, value interface{}, exists bool, rule string) {
	parts := strings.SplitN(rule, ":", 2)
	ruleName := parts[0]
	var param string
	if len(parts) == 2 {
		param = parts[1]
	}

	switch ruleName {
	case "required":
		if !exists || isEmpty(value) {
			v.errors.Add(field, fmt.Sprintf("O campo %s é obrigatório.", field))
		}

	case "min":
		if !exists || isEmpty(value) {
			return
		}
		var n int
		fmt.Sscanf(param, "%d", &n)
		if str, ok := value.(string); ok && len(str) < n {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ter no mínimo %d caracteres.", field, n))
		}
		if num, ok := toFloat(value); ok && num < float64(n) {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser no mínimo %s.", field, param))
		}

	case "max":
		if !exists || isEmpty(value) {
			return
		}
		var n int
		fmt.Sscanf(param, "%d", &n)
		if str, ok := value.(string); ok && len(str) > n {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ter no máximo %d caracteres.", field, n))
		}
		if num, ok := toFloat(value); ok && num > float64(n) {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser no máximo %s.", field, param))
		}

	case "email":
		if !exists || isEmpty(value) {
			return
		}
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if str, ok := value.(string); ok && !re.MatchString(str) {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser um e-mail válido.", field))
		}

	case "numeric":
		if !exists || isEmpty(value) {
			return
		}
		if _, ok := toFloat(value); !ok {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser numérico.", field))
		}

	case "gt": // greater than — equivalente ao Rule::gt() do Laravel
		if !exists || isEmpty(value) {
			return
		}
		var n float64
		fmt.Sscanf(param, "%f", &n)
		if num, ok := toFloat(value); ok && num <= n {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser maior que %s.", field, param))
		}

	case "in": // equivalente ao Rule::in([...]) do Laravel
		if !exists || isEmpty(value) {
			return
		}
		allowed := strings.Split(param, ",")
		str := fmt.Sprintf("%v", value)
		found := false
		for _, a := range allowed {
			if strings.EqualFold(str, strings.TrimSpace(a)) {
				found = true
				break
			}
		}
		if !found {
			v.errors.Add(field, fmt.Sprintf("O campo %s deve ser um dos valores: %s.", field, param))
		}
	}
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}
