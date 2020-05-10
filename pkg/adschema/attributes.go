package adschema

//go:generate go run gen.go
//go:generate go fmt

//var BooleanSyntax = AttributeSyntax{
//	Name:         "Boolean",
//	ConvertBytes: ConvertBool,
//}

//func (a *AttributeSyntax) UnmarshalJSON(data []byte) error {
//	var s string
//	if err := json.Unmarshal(data, &s); err != nil {
//		return err
//	}
//
//	switch s {
//	case "Boolean":
//		a.Name = "Boolean"
//		a.ConvertBytes = ConvertBool
//	case "Interval":
//		a.Name = "Interval"
//		a.ConvertBytes = func(b []byte) interface{} { return nil }
//	default:
//		a.Name = "Unknown"
//		a.ConvertBytes = func(b []byte) interface{} { return nil }
//	}
//	return nil
//}

//var EnumerationSyntax = AttributeSyntax{
//	Name:    "Enumeration",
//	Convert: ParseEnumeration,
//}
//
//var IntervalSyntax = AttributeSyntax{
//	Name:    "Interval",
//	Convert: ParseInterval,
//}
//
//var ObjectAccessPoint = AttributeSyntax{
//	Name:    "Object(Access-Point)",
//	Convert: ParseObject,
//}
