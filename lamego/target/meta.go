package target

type SerializeMeta struct {
	JSON bool
}

type ExtraMeta map[string]string

type ClassMeta struct {
	Serialize SerializeMeta
	Extra     ExtraMeta
}
