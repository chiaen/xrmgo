package xrmgo

import (
	"errors"
	"fmt"
	"time"

	"github.com/beevik/etree"
)

func (c *clientImpl) buildOCPRequest(username, password, region, loginUrl string) (string, error) {
	doc := etree.NewDocument()
	doc.ReadFromString(requestTemplate)
	root := doc.Root()

	id := root.FindElement("//MessageID")
	id.SetText(guid())
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

func (c *clientImpl) buildOCPHeader(action string) (*etree.Document, error) {
	if !c.isLoggedIn() {
		return nil, errors.New("not authenticated")
	}

	doc := etree.NewDocument()
	doc.ReadFromString(requestOCPHeader)
	root := doc.Root()

	act := root.FindElement("//Action")
	act.SetText(contracted(action))
	id := root.FindElement("//MessageID")
	id.SetText(guid())
	to := root.FindElement("//To")
	to.SetText(endpoint)
	now := time.Now()
	created := root.FindElement("//Created")
	created.SetText(toCurrentTime(now))
	expired := root.FindElement("//Expires")
	expired.SetText(toTomorrowTime(now))

	identifier := root.FindElement("//KeyIdentifier")
	identifier.SetText(c.keyIdentifier)
	token0 := root.FindElement("//EncryptedKey//CipherData//CipherValue")
	token0.SetText(c.securityToken0)
	token1 := root.FindElement("//EncryptedData//CipherData/CipherValue")
	token1.SetText(c.securityToken1)

	return doc, nil
}

func (c *clientImpl) buildEnvelope(action string) (*etree.Document, error) {
	header, err := c.buildOCPHeader(action)
	if err != nil {
		return nil, err
	}
	doc := etree.NewDocument()
	doc.ReadFromString(requestEnvelope)
	root := doc.Root()
	body := root.FindElement("//Body")
	root.InsertChild(body, header.Root())
	return doc, nil
}

func (c *clientImpl) buildCreateRequest(entity Xmler) (string, error) {
	request, err := c.buildEnvelope("Create")
	if err != nil {
		return "", err
	}
	env, err := request.WriteToString()
	fmt.Println("env: ", env)

	ent, err := entity.ToXML().WriteToString()
	fmt.Println("ent: ", ent)

	doc := etree.NewDocument()
	doc.ReadFromString(createWrapper)
	root := doc.Root()
	root.AddChild(entity.ToXML().Root())

	create, err := doc.WriteToString()
	fmt.Println("create wrapper: ", create)

	body := request.Root().FindElement("//Body")
	body.AddChild(root)

	return request.WriteToString()
}

var createWrapper = `
	  <Create xmlns="http://schemas.microsoft.com/xrm/2011/Contracts/Services" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
          </Create>
`

var requestEnvelope = `
	 <s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
          <s:Body>
          </s:Body>
         </s:Envelope>
`

var requestOCPHeader = `
          <s:Header>
           <a:Action s:mustUnderstand="1">http://schemas.microsoft.com/xrm/2011/Contracts/Services/IOrganizationService/#{action}</a:Action>
           <a:MessageID>
            urn:uuid:#{uuid()}
           </a:MessageID>
           <a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo>
           <a:To s:mustUnderstand="1">
            #{@organization_endpoint}
           </a:To>
           <o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
           <u:Timestamp u:Id="_0">
            <u:Created>#{get_current_time}</u:Created>
            <u:Expires>#{get_tomorrow_time}</u:Expires>
           </u:Timestamp>
           <EncryptedData Id="Assertion0" Type="http://www.w3.org/2001/04/xmlenc#Element" xmlns="http://www.w3.org/2001/04/xmlenc#">
            <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#tripledes-cbc"></EncryptionMethod>
            <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
             <EncryptedKey>
              <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p"></EncryptionMethod>
              <ds:KeyInfo Id="keyinfo">
               <wsse:SecurityTokenReference xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
                <wsse:KeyIdentifier EncodingType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary" ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509SubjectKeyIdentifier">
                 #{@key_identifier}
                </wsse:KeyIdentifier>
               </wsse:SecurityTokenReference>
              </ds:KeyInfo>
              <CipherData>
               <CipherValue>
                #{@security_token0}
               </CipherValue>
              </CipherData>
             </EncryptedKey>
            </ds:KeyInfo>
            <CipherData>
             <CipherValue>
              #{@security_token1}
             </CipherValue>
            </CipherData>
           </EncryptedData>
           </o:Security>
          </s:Header>
`

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
