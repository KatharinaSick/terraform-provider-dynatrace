package hcl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func equalLineByLine(s1, s2 string) bool {
	if len(s1) == 0 && len(s2) == 0 {
		return false
	}
	parts1 := strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s1, "\r\n", "\n"), "\r", "\n")), "\n")
	parts2 := strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s2, "\r\n", "\n"), "\r", "\n")), "\n")
	if len(parts1) != len(parts2) {
		return false
	}
	for idx := range parts1 {
		part1 := strings.TrimSpace(parts1[idx])
		part2 := strings.TrimSpace(parts2[idx])
		if part1 != part2 {
			return false
		}
	}
	return true
}

func SuppressEOT(k, old, new string, d *schema.ResourceData) bool {
	return equalLineByLine(old, new)
}

func JSONStringsEqual(s1, s2 string) bool {
	if len(s1) == 0 && len(s2) == 0 {
		return false
	}
	b1 := bytes.NewBufferString("")
	if err := json.Compact(b1, []byte(s1)); err != nil {
		return false
	}

	b2 := bytes.NewBufferString("")
	if err := json.Compact(b2, []byte(s2)); err != nil {
		return false
	}

	return JSONBytesEqual(b1.Bytes(), b2.Bytes())
}

func compact(v any) any {
	if v == nil {
		return nil
	}
	switch typedV := v.(type) {
	case map[string]any:
		m := map[string]any{}
		for mk, mv := range typedV {
			if mv != nil {
				if cv := compact(mv); cv != nil {
					m[mk] = cv
				}
			}
		}
		if len(m) == 0 {
			return nil
		}
		return m
	case []any:
		s := []any{}
		for _, sv := range typedV {
			if cv := compact(sv); cv != nil {
				s = append(s, cv)
			}
		}
		if len(s) == 0 {
			return nil
		}
		return s
	default:
		return typedV
	}
}

func JSONBytesEqual(b1, b2 []byte) bool {
	var o1 any
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}
	o1 = compact(o1)

	var o2 any
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}
	o2 = compact(o2)
	return DeepEqual(o1, o2)
}

func SuppressJSONorEOT(k, old, new string, d *schema.ResourceData) bool {
	if JSONStringsEqual(old, new) {
		return true
	}
	return equalLineByLine(old, new)
}

func DeepEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	if ta != tb {
		return false
	}

	switch aTyped := a.(type) {
	case bool:
		return aTyped == b.(bool)
	case string:
		return aTyped == b.(string)
	case int:
		return aTyped == b.(int)
	case int8:
		return aTyped == b.(int8)
	case int16:
		return aTyped == b.(int16)
	case int32:
		return aTyped == b.(int32)
	case int64:
		return aTyped == b.(int64)
	case uint:
		return aTyped == b.(uint)
	case uint8:
		return aTyped == b.(uint8)
	case uint16:
		return aTyped == b.(uint16)
	case uint32:
		return aTyped == b.(uint32)
	case uint64:
		return aTyped == b.(uint64)
	case float32:
		return aTyped == b.(float32)
	case float64:
		return aTyped == b.(float64)
	case map[string]any:
		bTyped := b.(map[string]any)
		if len(aTyped) != len(bTyped) {
			return false
		}
		for ka, va := range aTyped {
			// TODO: this is cheating. Essentially we're suppressing the diff even if metricExpressions are different
			// Reason: The REST API ALWAYS delivers that string in a different format
			// Solution: We would need to implement a metric expression parser
			if ka == "metricExpressions" {
				continue
			}
			if vb, ok := bTyped[ka]; !ok || !DeepEqual(va, vb) {
				return false
			}
		}
		return true
	case []any:
		bTyped := b.([]any)
		if len(aTyped) != len(bTyped) {
			return false
		}
		m := map[int]bool{}
		for i := 0; i < len(aTyped); i++ {
			m[i] = false
		}
		for _, aElem := range aTyped {
			found := false
			for i, bElem := range bTyped {
				if m[i] {
					continue
				}
				if DeepEqual(aElem, bElem) {
					found = true
					m[i] = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		panic(fmt.Sprintf("Type `%T` not supported yet", aTyped))
	}
}
