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

type XrmgoTestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func (s *XrmgoTestSuite) SetupSuite() {
	s.ts = httptest.NewServer(new(testHandler))
	fmt.Println("hostname: ", s.ts.URL)
}

func (s *XrmgoTestSuite) SetupTest() {
	InitParams("hc340000.crm5.dynamics.com")
}

func (s *XrmgoTestSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *XrmgoTestSuite) TestAuthSuccess() {
	u := ocpLoginURL
	ocpLoginURL = s.ts.URL
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

func TestXrmgoTestSuite(t *testing.T) {
	suite.Run(t, new(XrmgoTestSuite))
}
