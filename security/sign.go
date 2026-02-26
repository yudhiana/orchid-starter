package security

import (
	"fmt"
	"orchid-starter/internal/common"
	"time"
)

const (
	AsymmetricSignPattern = "{{method}}:{{relative-path}}:{{body-hash}}:{{timestamp}}"
)

type SignType string

const (
	AsymmetricSign SignType = "asymmetric"
)

type AsymetricPatternInput struct {
	Method       string
	RelativePath string
	BodyHash     string
	Timestamp    time.Time
	formatTime   string
}

func GenerateSignPattern(signType SignType, input any) (pattern string, err error) {
	switch signType {
	case AsymmetricSign:
		in, ok := input.(AsymetricPatternInput)
		if !ok {
			return "", fmt.Errorf("invalid input for asymmetric sign pattern: %T", input)
		}

		// build the string based on the predefined pattern
		newSign, errRender := common.Render(AsymmetricSignPattern, map[string]any{
			"method":        in.Method,
			"relative-path": in.RelativePath,
			"body-hash":     in.BodyHash,
			"timestamp":     in.Timestamp.Format(in.formatTime),
		})

		if errRender != nil {
			return "", fmt.Errorf("failed to render asymmetric sign pattern Error: %w", errRender)
		}

		pattern = newSign
	default:
		return "", fmt.Errorf("unsupported sign type: %s", signType)
	}
	return
}
