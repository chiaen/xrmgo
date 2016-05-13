package xrmgo

import (
	"errors"
)

type clientImpl struct {
}

func GetClient() Client {
	return &clientImpl{}
}

func (c *clientImpl) Auth(username, password string) (bool, error) {
	return false, errors.New("No impl")
}

func (c *clientImpl) User() (*User, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Create(entity, guid string, criteria map[string]interface{}, columns ...string) ([]*Entity, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Update(entity, guid string, attrs map[string]interface{}) (*UpdateResponse, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Delete(entity, guid string) (*DeleteResponse, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Execute(action string, params map[string]interface{}) error {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Fetch() {
}

func (c *clientImpl) Associate(entity, guid, relation, relatedEntity, relatedId string) error {
	return errors.New("No impl")
}

func (c *clientImpl) Describe(entity ...string) (map[string]*MetaData, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) DescribeAttr(entity, field string) (*MetaData, error) {
	return nil, errors.New("No impl")
}
