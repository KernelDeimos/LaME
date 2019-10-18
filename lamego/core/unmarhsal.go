package core

// Everything in this file and relayed to the YAML library is temporary.
// Eventually, a new YAML parser needs to be written in pure LaME, and
// the generated output of the first successful run will be the bootstrapping
// YAML parser.

// TODO: Code generation should be used to make these things due to the
//       unfortunate requirement of writing them separately.

type ArgumentForYAML struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type FieldForYAML struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type ModelForYAML struct {
	ID      string   `yaml:"id"`
	Type    string   `yaml:"type"`
	Fields  []Field  `yaml:"fields"`
	Methods []Method `yaml:"methods"`
}

// TODO: Code generation should be used to make these things due to the
//       unfortunate requirement of writing them separately.

func (out *Argument) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := struct {
		ArgumentNormal ArgumentForYAML `yaml:"args"`
		ArgumentList   []string        `yaml:"args"`
	}{}
	var err error

	err = unmarshal(&tmp.ArgumentNormal)
	if err == nil {
		*out = Argument{
			Name: tmp.ArgumentNormal.Name,
			Type: tmp.ArgumentNormal.Type,
		}
		return nil
	}

	err = unmarshal(&tmp.ArgumentList)
	if err == nil {
		*out = Argument{
			Name: tmp.ArgumentList[0],
			Type: tmp.ArgumentList[1],
		}
		return nil
	}

	return err
}

func (out *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := struct {
		FieldNormal FieldForYAML `yaml:"args"`
		FieldList   []string     `yaml:"args"`
		FieldString string       `yaml:"args"`
	}{}
	var err error

	err = unmarshal(&tmp.FieldNormal)
	if err == nil {
		*out = Field{
			Name: tmp.FieldNormal.Name,
			Type: tmp.FieldNormal.Type,
		}
		return nil
	}

	err = unmarshal(&tmp.FieldList)
	if err == nil {
		typ := ""
		if len(tmp.FieldList) > 1 {
			typ = tmp.FieldList[1]
		}
		*out = Field{
			Name: tmp.FieldList[0],
			Type: typ,
		}
		return nil
	}

	err = unmarshal(&tmp.FieldString)
	if err == nil {
		// TODO: maybe if name contains dots it infers this as a type and
		//       provides a resonable default for the name... maybe
		*out = Field{
			Name: tmp.FieldString,
		}
		return nil
	}

	return err
}

func (out *Model) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := struct {
		RealModel    ModelForYAML
		CommandModel string
	}{}
	var err error
	var errRealModel error

	err = unmarshal(&tmp.RealModel)
	if err == nil {
		*out = Model{
			ID:      tmp.RealModel.ID,
			Type:    tmp.RealModel.Type,
			Fields:  tmp.RealModel.Fields,
			Methods: tmp.RealModel.Methods,
		}
		return nil
	}

	// TODO: concatenate errors to a single object, but for now this is
	//       a bit better
	errRealModel = err

	err = unmarshal(&tmp.CommandModel)
	if err == nil {
		*out = Model{
			Type: "command",
			ID:   tmp.CommandModel,
		}
		return nil
	}

	return errRealModel
}
