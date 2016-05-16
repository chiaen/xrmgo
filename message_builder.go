package xrmgo

import (
	"fmt"
	"time"

	"github.com/beevik/etree"
	"github.com/satori/go.uuid"
)

func (c *clientImpl) buildOCPRequest(username, password, region, loginUrl string) (string, error) {
	doc := etree.NewDocument()
	doc.ReadFromString(requestTemplate)
	root := doc.Root()

	fmt.Println("root: ", root.Tag)

	id := root.FindElement("//MessageID")
	id.SetText(fmt.Sprintf("urn:uuid:%s", uuid.NewV4()))
	to := root.FindElement("//To")
	to.SetText(loginUrl)

	now := time.Now()
	created := root.FindElement("//Created")
	created.SetText(toCurrentTime(now))
	expired := root.FindElement("//Expires")
	expired.SetText(toTomorrowTime(now))

	name := root.FindElement("//Username")
	name.SetText(username)
	pass := root.FindElement("//Password")
	pass.SetText(password)

	addr := root.FindElement("//EndpointReference//Address")
	addr.SetText(region)

	return doc.WriteToString()
}

var requestTemplate = `
	  <s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
            xmlns:a="http://www.w3.org/2005/08/addressing"
            xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
            <s:Header>
              <a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/02/trust/RST/Issue</a:Action>
              <a:MessageID></a:MessageID>
              <a:ReplyTo>
                <a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
              </a:ReplyTo>
              <a:To s:mustUnderstand="1">%s</a:To>
              <o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
                <u:Timestamp u:Id="_0">
                  <u:Created></u:Created>
                  <u:Expires></u:Expires>
                </u:Timestamp>
                <o:UsernameToken u:Id="uuid-cdb639e6-f9b0-4c01-b454-0fe244de73af-1">
                  <o:Username></o:Username>
                  <o:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">
                  </o:Password>
                </o:UsernameToken>
              </o:Security>
            </s:Header>
            <s:Body>
              <t:RequestSecurityToken xmlns:t="http://schemas.xmlsoap.org/ws/2005/02/trust">
                <wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
                  <a:EndpointReference>
                    <a:Address></a:Address>
                  </a:EndpointReference>
                </wsp:AppliesTo>
                <t:RequestType>http://schemas.xmlsoap.org/ws/2005/02/trust/Issue</t:RequestType>
              </t:RequestSecurityToken>
            </s:Body>
          </s:Envelope>
`
