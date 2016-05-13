package xrmgo

type Authenticator interface {
	Auth(username, password string) (bool, error)

	User() (*User, error)
}

type Client interface {
	Authenticator

	Create(entity string, attrs map[string]interface{}) (string, error)

	Retrieve(entity, guid string, criteria map[string]interface{}, columns ...string) ([]*Entity, error)

	Update(entity, guid string, attrs map[string]interface{}) (*UpdateResponse, error)

	Delete(entity, guid string) (*DeleteResponse, error)

	Execute(action string, params map[string]interface{}) error

	Fetch()

	Associate(entity, guid, relation, relatedEntity, relatedId string) error

	Describe(entity ...string) (map[string]*MetaData, error)

	DescribeAttr(entity, field string) (*MetaData, error)
}

type Entity struct {
}

type User struct {
}

type UpdateResponse struct {
}

type DeleteResponse struct {
}

type MetaData struct {
}
