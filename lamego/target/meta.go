package target

type SerializeMeta struct {
	JSON bool
}

type ExtraMeta map[string]string

type ClassMeta struct {
	Serialize   SerializeMeta
	SourceModel string
	Extra       ExtraMeta
}
