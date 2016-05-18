package xrmgo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func fixtures(name string) []byte {
	bytes, _ := ioutil.ReadFile(filepath.Join("fixtures", name+".xml"))
	return bytes
}

type testHandler struct{}

func (t *testHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("req: ", req)
	resp.Write(fixtures("request_security_token_response"))
}

func mockServer(code int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write(body)
		}))
}

type XrmgoTestSuite struct {
	suite.Suite
	c *clientImpl
}

func (s *XrmgoTestSuite) SetupTest() {
	InitParams("hc340000.crm5.dynamics.com")

	s.c = &clientImpl{
		securityToken0: "token0",
		securityToken1: "token1",
		keyIdentifier:  "identifier",
	}
}

func (s *XrmgoTestSuite) TestAuthSuccess() {
	server := mockServer(200, fixtures("request_security_token_response"))
	defer server.Close()
	u := ocpLoginURL
	ocpLoginURL = server.URL
	c := &clientImpl{}
	result, err := c.Auth("testing", "password")
	s.NoError(err)
	s.True(result)
	s.Contains(c.securityToken0, "tMFpDJbJHcZnRVuby5cYmRbCJo2OgOFLEOrUHj+wz")
	s.Contains(c.securityToken1, "CX7BFgRnW75tE6GiuRICjeVDV+6q4KDMKLyKmKe9A8U")
	s.Equal(c.keyIdentifier, "D3xjUG3HGaQuKyuGdTWuf6547Lo=")
	ocpLoginURL = u
}

func (s *XrmgoTestSuite) TestAuthFailed() {
	c := &clientImpl{}
	result, err := c.Auth("test@test.onmicrosoft.com", "qwerty")
	s.False(result)
	s.EqualError(err, "Authentication Failure")
}

func (s *XrmgoTestSuite) TestCreateSuccess() {
	server := mockServer(200, fixtures("create_response"))
	defer server.Close()
	u := endpoint
	endpoint = server.URL
	response, err := s.c.Create("accout", map[string]interface{}{"name": "HTC"})
	s.NoError(err)
	s.Equal("c4944f99-b5a0-e311-b64f-6c3be5a87df0", response)
	endpoint = u
}

func TestXrmgoTestSuite(t *testing.T) {
	suite.Run(t, new(XrmgoTestSuite))
}
