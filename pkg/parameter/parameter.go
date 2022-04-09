package parameter

import (
	"strconv"
	"strings"
)

var params map[string]string = make(map[string]string)

func Parse(parameter string) (map[string]string, error) {
	for k := range params {
		delete(params, k)
	}

	if strings.TrimSpace(parameter) == "" {
		return params, nil
	}

	for _, p := range strings.Split(strings.TrimSpace(parameter), ",") {
		if strings.Contains(p, "=") {
			parts := strings.Split(p, "=")
			params[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		} else if strings.TrimSpace(p) == "" {
			continue
		} else {
			params[p] = "true"
		}
	}

	return params, nil
}

func Bool(key string, def bool) bool {
	if val, ok := params[key]; !ok {
		return def
	} else if val, err := strconv.ParseBool(val); err != nil {
		return def
	} else {
		return val
	}
}

func Int(key string, def int) int {
	if val, ok := params[key]; !ok {
		return def
	} else if val, err := strconv.Atoi(val); err != nil {
		return def
	} else {
		return val
	}
}

func String(key string, def string) string {
	if val, ok := params[key]; !ok {
		return def
	} else {
		return val
	}
}
