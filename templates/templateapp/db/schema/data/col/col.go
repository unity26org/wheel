package col

var Path = []string{"db", "schema", "data", "col", "col.go"}

var Content = `package col

// CREATE INDEX index_assets_on_address ON public.assets USING btree (address);
type Index struct {
	Name   string
	Unique bool
}

// Options are: "default", "null", "precision", "scale", "unique", "index" and "limit"
type Info struct {
	Name    string
	Type    string
	Options map[string]interface{}
}

func String(name string, options map[string]interface{}) Info {
	return setInfo(name, "VARCHAR", options)
}

func Text(name string, options map[string]interface{}) Info {
	return setInfo(name, "TEXT", options)
}

func Boolean(name string, options map[string]interface{}) Info {
	return setInfo(name, "BOOLEAN", options)
}

func Smallint(name string, options map[string]interface{}) Info {
	return setInfo(name, "SMALLINT", options)
}

func Integer(name string, options map[string]interface{}) Info {
	return setInfo(name, "INTEGER", options)
}

func Bigint(name string, options map[string]interface{}) Info {
	return setInfo(name, "BIGINT", options)
}

func Datetime(name string, options map[string]interface{}) Info {
	return setInfo(name, "TIMESTAMP", options)
}

func Numeric(name string, options map[string]interface{}) Info {
	return setInfo(name, "NUMERIC", options)
}

func References(name string, options map[string]interface{}) Info {
	return setInfo(name, "REFERENCES", options)
}

func setInfo(name string, tType string, options map[string]interface{}) Info {
	if options == nil {
		options = map[string]interface{}{}
	}

	return Info{Name: name, Type: tType, Options: options}
}`
