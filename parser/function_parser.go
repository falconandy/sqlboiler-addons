package parser

import (
	"regexp"
	"strings"

	"github.com/falconandy/sqlboiler-addons/model"
)

var (
	signatureRe = regexp.MustCompile(`^(\w+)\s*\(([^)]*)\)(.*)$`)
)

func ParseFunctionSignature(rawSignature string) (model.FunctionSignature, bool) {
	match := signatureRe.FindStringSubmatch(strings.TrimSpace(rawSignature))
	if match == nil {
		return model.FunctionSignature{}, false
	}

	functionName := match[1]
	functionParameters := strings.TrimSpace(match[2])
	functionReturn := strings.TrimSpace(match[3])

	parameters, ok := parseFunctionParameters(functionParameters)
	if !ok {
		return model.FunctionSignature{}, false
	}

	return model.FunctionSignature{
		Name:       functionName,
		Parameters: parameters,
		Return:     functionReturn,
	}, true
}

func parseFunctionParameters(rawParameters string) ([]model.FunctionParameter, bool) {
	if len(rawParameters) == 0 {
		return nil, true
	}

	var parameters []model.FunctionParameter
	var pendingParameterNames []string
	for _, p := range strings.Split(rawParameters, ",") {
		p = strings.TrimSpace(p)
		parts := strings.Fields(p)
		switch len(parts) {
		case 1:
			pendingParameterNames = append(pendingParameterNames, parts[0])
		case 2:
			parameters = append(parameters, model.FunctionParameter{
				Names: append(pendingParameterNames, parts[0]),
				Type:  parts[1],
			})
			pendingParameterNames = nil
		default:
			return nil, false
		}
	}
	if len(pendingParameterNames) > 0 {
		return nil, false
	}
	return parameters, true
}
