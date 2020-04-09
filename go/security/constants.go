package security

const (
	// ActionAll matches any action
	ActionAll = Action("*")

	// create, read, update, delete, list actions compatible with restful api methods

	// ActionCreate should be used for creation actions
	ActionCreate = Action("create")
	// ActionRead should be used for reading actions
	ActionRead = Action("read")
	// ActionUpdate should be used for updating actions
	ActionUpdate = Action("update")
	// ActionDelete should be used for deletion actions
	ActionDelete = Action("delete")
	// ActionList should be used for listing actions
	ActionList = Action("list")

	// ResourceAll matches any resource
	ResourceAll = Resource("*")
)
