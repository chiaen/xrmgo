package xrmgo

import (
	"errors"
	"net/http"
	"strings"

	"github.com/beevik/etree"
)

var (
	ocpLoginURL = "https://login.microsoftonline.com/RST2.srf"

	hostName   string
	regionName string
)

func InitParams(hostname string) error {
	if hostname == "" {
		return errors.New("empty hostname")
	}
	r := region[re.FindString(hostname)]
	if r == "" {
		return errors.New("invalid hostname")
	}
	hostName = hostname
	regionName = r
	return nil
}

type clientImpl struct {
	securityToken0 string
	securityToken1 string
	keyIdentifier  string
}

func GetClient() Client {
	return &clientImpl{}
}

func (c *clientImpl) Auth(username, password string) (bool, error) {
	body, err := c.buildOCPRequest(username, password, regionName, ocpLoginURL)
	if err != nil {
		return false, err
	}
	resp, err := http.Post(ocpLoginURL, "application/soap+xml; charset=utf-8", strings.NewReader(body))
	if err != nil {
		return false, err
	}
	doc := etree.NewDocument()
	_, err = doc.ReadFrom(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	root := doc.Root()
	faults := root.FindElements("//Fault")
	if len(faults) > 0 {
		reason := root.FindElement("//Reason//Text")
		return false, errors.New(reason.Text())
	}
	securityTokens := root.FindElements("//CipherValue")
	keyIdentifier := root.FindElement("//KeyInfo//SecurityTokenReference//KeyIdentifier")
	if len(securityTokens) != 2 || keyIdentifier == nil {
		return false, errors.New("Invalid Response")
	}
	c.securityToken0 = securityTokens[1].Text()
	c.securityToken1 = securityTokens[0].Text()
	c.keyIdentifier = keyIdentifier.Text()
	return true, nil
}

func (c *clientImpl) User() (*User, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Create(entity string, attrs map[string]interface{}) (string, error) {
	return "", errors.New("No impl")
}

func (c *clientImpl) Retrieve(entity, guid string, criteria map[string]interface{}, columns ...string) ([]*Entity, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Update(entity, guid string, attrs map[string]interface{}) (*UpdateResponse, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Delete(entity, guid string) (*DeleteResponse, error) {
	return nil, errors.New("No impl")
}

func (c *clientImpl) Execute(action string, params map[string]interface{}) error {
	return errors.New("No impl")
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
