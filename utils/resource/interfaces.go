package resource

// IResource describe an API resource.
type IResource interface {
	// GetLinks returns the resource links
	GetLinks() (any, error)
	GetName() (string, error)
	GetTitle() (string, error)
	GetType() string
}
