package content

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/helper"
)

type VariableHandler func(args ...string) string

var variableHandlers = map[string]VariableHandler{
	"rel": relVarHandler,
}

// this function handles the $rel variable
func relVarHandler(args ...string) string {
	log := helper.AcquireLogger()
	p := args[0]
	cmd := fmt.Sprintf("$rel(%s)", p)
	conf := config.Get()
	if conf.PublicURL == "" {
		log.Warn().Str("could not replace", cmd).Msg("CENTRA_PUBLIC_URL not set")
		return cmd
	}
	u, err := url.Parse(conf.PublicURL)
	if err != nil {
		log.Error().Err(err).Msg("could not parse url")
		return cmd
	}
	u.Path = path.Join(append([]string{u.Path}, strings.Split(p, ",")...)...)
	return u.String()
}

func HandleVariable(name string, args ...string) VariableHandler {
	name = strings.ToLower(name)

	if vh, ok := variableHandlers[name]; ok {
		return vh
	}

	return handleVariableIgnore
}

func handleVariableIgnore(args ...string) string { return "" }

var varRegex = regexp.MustCompile(`\$([a-zA-Z0-9]+)\((.*?)\)`)

func ProcessVariables(input string) string {
	return varRegex.ReplaceAllStringFunc(input, func(match string) string {
		submatches := varRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		name := submatches[1]
		argsRaw := submatches[2]

		var args []string
		if len(strings.TrimSpace(argsRaw)) > 0 {
			args = strings.Split(argsRaw, ",")
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}
		}

		handler := HandleVariable(name)

		return handler(args...)
	})
}

func ProcessMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	result := make(map[string]any)

	for k, v := range input {
		switch val := v.(type) {
		case string:
			result[k] = ProcessVariables(val)

		case map[string]any:
			result[k] = ProcessMap(val)

		case []any:
			newSlice := make([]any, len(val))
			for i, item := range val {
				switch itemVal := item.(type) {
				case string:
					newSlice[i] = ProcessVariables(itemVal)
				case map[string]any:
					newSlice[i] = ProcessMap(itemVal)
				default:
					newSlice[i] = itemVal
				}
			}
			result[k] = newSlice

		default:
			result[k] = v
		}
	}

	return result
}
