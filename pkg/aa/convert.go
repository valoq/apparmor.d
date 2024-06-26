// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

// Must is a helper that wraps a call to a function returning (any, error) and
// panics if the error is non-nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// cmpFileAccess compares two access strings for file rules.
// It is aimed to be used in slices.SortFunc.
func cmpFileAccess(i, j string) int {
	if slices.Contains(requirements[FILE]["access"], i) &&
		slices.Contains(requirements[FILE]["access"], j) {
		return requirementsWeights[FILE]["access"][i] - requirementsWeights[FILE]["access"][j]
	}
	if slices.Contains(requirements[FILE]["transition"], i) &&
		slices.Contains(requirements[FILE]["transition"], j) {
		return requirementsWeights[FILE]["transition"][i] - requirementsWeights[FILE]["transition"][j]
	}
	if slices.Contains(requirements[FILE]["access"], i) {
		return -1
	}
	return 1
}

func validateValues(kind Kind, key string, values []string) error {
	for _, v := range values {
		if v == "" {
			continue
		}
		if !slices.Contains(requirements[kind][key], v) {
			return fmt.Errorf("invalid mode '%s'", v)
		}
	}
	return nil
}

// Helper function to convert a string to a slice of rule values according to
// the rule requirements as defined in the requirements map.
func toValues(kind Kind, key string, input string) ([]string, error) {
	req, ok := requirements[kind][key]
	if !ok {
		return nil, fmt.Errorf("unrecognized requirement '%s' for rule %s", key, kind)
	}

	res := tokenToSlice(input)
	for idx := range res {
		res[idx] = strings.Trim(res[idx], `" `)
		if res[idx] == "" {
			res = slices.Delete(res, idx, idx+1)
			continue
		}
		if !slices.Contains(req, res[idx]) {
			return nil, fmt.Errorf("unrecognized %s: %s", key, res[idx])
		}
	}
	slices.SortFunc(res, func(i, j string) int {
		return requirementsWeights[kind][key][i] - requirementsWeights[kind][key][j]
	})
	return slices.Compact(res), nil
}

// Helper function to convert an access string to a slice of access according to
// the rule requirements as defined in the requirements map.
func toAccess(kind Kind, input string) ([]string, error) {
	var res []string

	switch kind {
	case FILE:
		raw := strings.Split(input, "")
		trans := []string{}
		for _, access := range raw {
			if slices.Contains(requirements[FILE]["access"], access) {
				res = append(res, access)
			} else {
				trans = append(trans, access)
			}
		}

		transition := strings.Join(trans, "")
		if len(transition) > 0 {
			if slices.Contains(requirements[FILE]["transition"], transition) {
				res = append(res, transition)
			} else {
				return nil, fmt.Errorf("unrecognized transition: %s", transition)
			}
		}

	case FILE + "-log":
		raw := strings.Split(input, "")
		for _, access := range raw {
			if slices.Contains(requirements[FILE]["access"], access) {
				res = append(res, access)
			} else if maskToAccess[access] != "" {
				res = append(res, maskToAccess[access])
			} else {
				return nil, fmt.Errorf("toAccess: unrecognized file access '%s' for %s", input, kind)
			}
		}

	default:
		return toValues(kind, "access", input)
	}

	slices.SortFunc(res, cmpFileAccess)
	return slices.Compact(res), nil
}
