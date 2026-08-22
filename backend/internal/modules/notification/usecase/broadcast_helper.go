package usecase

import "reflect"

func toBroadcastMap(payload any) (map[string]string, bool) {
	v := reflect.ValueOf(payload)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, false
	}
	out := map[string]string{}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() != reflect.String {
			continue
		}
		out[stringsLower(f.Name)] = v.Field(i).String()
	}
	if out["title"] == "" && out["Title"] == "" {
		if title, ok := out["Title"]; ok {
			out["title"] = title
		}
	}
	title := out["title"]
	if title == "" {
		title = fieldString(v, "Title")
	}
	body := fieldString(v, "Body")
	if title == "" && body == "" {
		return nil, false
	}
	return map[string]string{"title": title, "body": body}, true
}

func fieldString(v reflect.Value, name string) string {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

func stringsLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
