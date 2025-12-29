package main

func TransformType(typeStr, fieldName string) string {
	switch fieldName {
	case "deleted_at":
		return "soft_delete.DeletedAt"
	}

	// TODO year & spatial data types

	switch typeStr {
	case "bit", "tinyint", "bool", "boolean", "smallint", "mediumint", "int", "integer":
		return "int64"
	case "bigint":
		return "int64"
	case "decimal", "dec", "float", "double":
		return "float64"
	case "date", "time", "datetime", "timestamp":
		return "time.Time"
	case "char", "varchar",
		"binary", "varbinary",
		"tinyblob", "tinytext",
		"blob", "text",
		"mediumblob", "mediumtext",
		"longblob", "longtext",
		"enum", "set", "json":
		return "string"
	default:
		return "interface{}"
	}
}
