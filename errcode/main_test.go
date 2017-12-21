package errcode

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

/*

<Array>
	<data>
	<CustType>1</CustType>
	<CustStatus>1</CustStatus>
	<ThirdCustId>19871116</ThirdCustId>
	<CustName>tiny</CustName>
	<TotalBalance>10.00</TotalBalance>
	<TotalFreezeAmount>0.00</TotalFreezeAmount>
	<OrderAmount>0.00</OrderAmount>
	<TotalAmount>10.00</TotalAmount>
	<WayAmount>0.00</WayAmount>
	<TranDate>20170524</TranDate>
	</data>
</Array>

*/

type Data struct {
	XMLName     xml.Name `xml:"data"`
	CustType    string   `xml:"CustType"`
	CustStatus  string   `xml:"CustStatus"`
	ThirdCustId string   `xml:"ThirdCustId"`
}

type TestXml struct {
	XMLName  xml.Name `xml:"Array"`
	TestData []Data   `xml:"data"`
}

func TestMain(t *testing.T) {

	var x interface{}

	x = `
<Array>
	<data>
		<CustType>1</CustType>
		<CustStatus>1</CustStatus>
		<ThirdCustId>19871116</ThirdCustId>
	</data>
</Array>		`

	str, _ := json.Marshal(x)

	fmt.Println("...", string(str))

	var tx TestXml
	if err := xml.Unmarshal([]byte(`
<Array>
	<data>
		<CustType>1</CustType>
		<CustStatus>1</CustStatus>
		<ThirdCustId>19871116</ThirdCustId>
	</data>
</Array>`), &tx); err != nil {
		fmt.Println("...", err)
		t.Log("111", err)
	}
	t.Log("2222", tx)
	fmt.Println("...1", tx)
	fmt.Println(tx.TestData[0].CustStatus)

}
