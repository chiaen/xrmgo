package xrmgo

import (
	"github.com/beevik/etree"
)

const (
	emptyID = "00000000-0000-0000-0000-000000000000"
)

type Attributes map[string]interface{}

func (a *Attributes) ToXML() *etree.Document {
	if len(*a) == 0 {
		return nil
	}
	doc := etree.NewDocument()
	doc.ReadFromString(attributeRequest)
	root := doc.Root()
	for k, v := range *a {
		kvPair := root.CreateElement("a:KeyValuePairOfstringanyType")
		key := kvPair.CreateElement("c:key")
		key.SetText(k)
		value := kvPair.CreateElement("c:value")
		value.CreateAttr("i:type", "s:string")
		value.CreateAttr("xmlns:s", "http://www.w3.org/2001/XMLSchema")
		value.SetText(v.(string))
	}
	return doc
}

var attributeRequest = `
       <a:Attributes xmlns:c='http://schemas.datacontract.org/2004/07/System.Collections.Generic'>
       </a:Attributes>
`

type Entity struct {
	Attr        Attributes
	ID          string
	logicalName string
}

func (e *Entity) ToXML() *etree.Document {
	doc := etree.NewDocument()
	doc.ReadFromString(entityRequest)
	root := doc.Root()
	entityState := root.FindElement("//EntityState")
	root.InsertChild(entityState, e.Attr.ToXML().Root())

	if e.ID == "" {
		e.ID = emptyID
	}
	id := root.FindElement("//Id")
	id.SetText(e.ID)
	logicalName := root.FindElement("//LogicalName")
	logicalName.SetText(e.logicalName)
	return doc
}

var entityRequest = `
        <entity xmlns:a="http://schemas.microsoft.com/xrm/2011/Contracts">
          <a:EntityState i:nil="true" />
          <a:FormattedValues xmlns:b="http://schemas.datacontract.org/2004/07/System.Collections.Generic" />
          <a:Id></a:Id>
          <a:LogicalName></a:LogicalName>
          <a:RelatedEntities xmlns:b="http://schemas.datacontract.org/2004/07/System.Collections.Generic" />
	</entity>
`
