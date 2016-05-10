package jsonhelper

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"math"
	"strings"
	"testing"

	. "github.com/fxding/jsonhelper/assert"
	"github.com/bmizerany/assert"
)

//  go test -bench=".*"
//  go test -run="(Util)"

var (
	jh JsonHelper
)

func init() {
	json.Unmarshal([]byte(`{
		"name":"aaron",
		"nullstring":null,
		"ints":[1,2,3,4],
		"int":1,
		"intstr":"1",
		"int64":1234567890,
		"float64":123.456,
		"float64str":"123.456",
		"float64null": null,
		"MaxSize" : 1048576,
		"strings":["string1"],
		"stringscsv":"string1,string2",
		"nested":{
			"nest":"string2",
			"strings":["string1"],
			"int":2,
			"list":["value"],
			"nest2":{
				"test":"good"
			}
		},
		"nested2":[
			{"sub":2}
		],
		"period.name":"value"
	}`), &jh)
}

func TestJsonRawWriter(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`"hello"`)
	//raw := json.RawMessage(buf.Bytes())
	//bya, _ := json.Marshal(&buf)
	//fmt.Println(string(bya))
	//bya, _ = json.Marshal(&raw)
	//fmt.Println(string(bya))

	/*
		bya, err := json.Marshal(buf)
		Assert(string(bya) == `"hello"`, t, "Should be hello but was %s", string(bya))
		Debug(string(buf.Bytes()), err)
		var jrw JsonRawWriter
		jrw.WriteString(`"hello"`)
		Debug(jrw.Raw())
		bya, err = json.Marshal(jrw.Raw())
		Assert(string(bya) == `"hello"`, t, "Should be hello but was %s", string(bya))
		Debug(string(jrw.Bytes()), err)
	*/
}

func TestJsonHelper(t *testing.T) {

	Assert(jh.String("name", "") == "aaron", t, "should get 'aaron' %s", jh.String("name", ""))
	Assert(jh.String("nullstring", "") == "", t, "should get '' %s", jh.String("nullstring", ""))

	Assert(jh.Int("int", 10) == 1, t, "get int ")
	Assert(jh.Int("ints[0]", 10) == 1, t, "get int from array %d", jh.Int("ints[0]", 10))
	Assert(jh.Int("ints[2]", 10) == 3, t, "get int from array %d", jh.Int("ints[0]", 10))
	Assert(len(jh.Ints("ints", []int{})) == 4, t, "get int array %v", jh.Ints("ints", []int{}))
	Assert(jh.Int64("int64", 10) == 1234567890, t, "get int")
	Assert(jh.Int("nested.int", 10) == 2, t, "get int")
	Assert(jh.String("nested.nest", "") == "string2", t, "should get string %s", jh.String("nested.nest", ""))
	Assert(jh.String("nested.nest2.test", "") == "good", t, "should get string %s", jh.String("nested.nest2.test", ""))
	Assert(jh.String("nested.list[0]", "") == "value", t, "get string from array")
	Assert(jh.Int("nested2[0].sub", 10) == 2, t, "get int from obj in array %d", jh.Int("nested2[0].sub", 10))

	Assert(jh.Int("MaxSize", 10) == 1048576, t, "get int, test capitalization? ")
	sl := jh.Strings("strings", []string{})
	Assert(len(sl) == 1 && sl[0] == "string1", t, "get strings ")
	sl = jh.Strings("stringscsv", []string{})
	Assert(len(sl) == 2 && sl[0] == "string1", t, "get strings ")

	l := jh.List("nested2", nil)
	Assert(len(l) == 1, t, "get list")


	fv := jh.Float64("float64null", math.NaN())
	assert.Tf(t, math.IsNaN(fv), "float64null expected Nan but got %v", fv)

	jhm := jh.Helpers("nested2")
	Assert(len(jhm) == 1, t, "get list of helpers")
	Assert(jhm[0].Int("sub", 10) == 2, t, "Should get list of helpers")
}

func TestJsonHelper_DefaultValue(t *testing.T) {
	Assert(jh.String("nostring", "nostring") == "nostring", t, "should get string %s", jh.String("nostring", "nostring"))
	strs := jh.Strings("nostrings", []string{"nostring"})
	Assert(len(strs) == 1 && strs[0] == "nostring", t, "get strings")

	Assert(jh.Bool("nobool", true) == true, t, "should get true")
	Assert(jh.Int("noint", 10) == 10, t, "should get 10")

	is := jh.Ints("nostring", []int{20, 10})
	Assert(len(is) == 2 && is[1] == 10, t, "should get []int")
	Assert(jh.Int64("noint64", 64) == 64, t, "should get int64")
	Assert(jh.Float64("nofloat64", 64) == 64, t, "should get float64")
	Assert(jh.Uint64("nouint64", 64) == 64, t, "should get uint64")

	m := jh.Map("nomap", map[string]interface{}{"name": "test"})
	Assert(len(m) == 1 && m["name"] == "test", t, "should get string test")

	li := jh.List("nolist", []interface{}{1, 2, 3})
	Assert(len(li) == 3, t, "should get list")
}

func TestJsonInterface(t *testing.T) {

	var jim map[string]JsonInterface
	err := json.Unmarshal([]byte(`{
		"nullstring":null,
		"string":"string",
		"int":22,
		"float":22.2,
		"floatstr":"22.2",
		"intstr":"22"
	}`), &jim)
	Assert(err == nil, t, "no error:%v ", err)
	Assert(jim["nullstring"].StringSh() == "", t, "nullstring: %v", jim["nullstring"])
	Assert(jim["string"].StringSh() == "string", t, "nullstring: %v", jim["string"])
	Assert(jim["int"].IntSh() == 22, t, "int: %v", jim["int"])
	Assert(jim["int"].StringSh() == "22", t, "int->string: %v", jim["int"])
	Assert(jim["int"].FloatSh() == float32(22), t, "int->float: %v", jim["int"])
	Assert(jim["float"].FloatSh() == 22.2, t, "float: %v", jim["float"])
	Assert(jim["float"].StringSh() == "22.2", t, "float->string: %v", jim["float"])
	Assert(jim["float"].IntSh() == 22, t, "float->int: %v", jim["float"])
	Assert(jim["intstr"].IntSh() == 22, t, "intstr: %v", jim["intstr"])
	Assert(jim["intstr"].FloatSh() == float32(22), t, "intstr->float: %v", jim["intstr"])
}

func TestJsonCoercion(t *testing.T) {
	assert.Tf(t, jh.Int("intstr", 10) == 1, "get string as int %s", jh.String("intstr", ""))
	assert.Tf(t, jh.String("int", "") == "1", "get int as string %s", jh.String("int", ""))
	assert.Tf(t, jh.Int("notint", 100) == 100, "get non existent int = 0??? ")
}

func TestJsonPathNotation(t *testing.T) {
	// Now lets test xpath type syntax
	assert.Tf(t, jh.Int("/MaxSize", 10) == 1048576, "get int, test capitalization? ")
	assert.Tf(t, jh.String("/nested/nest", "") == "string2", "should get string %s", jh.String("/nested/nest", ""))
	assert.Tf(t, jh.String("/nested/list[0]", "") == "value", "get string from array")
	// note this one has period in name
	assert.Tf(t, jh.String("/period.name", "") == "value", "test period in name ")
}

func TestFromReader(t *testing.T) {
	raw := `{"testing": 123}`
	reader := strings.NewReader(raw)
	jh, err := NewJsonHelperReader(reader)
	assert.Tf(t, err == nil, "Unexpected error decoding json: %s", err)
	assert.Tf(t, jh.Int("testing", 10) == 123, "Unexpected value in json: %d", jh.Int("testing", 10))
}

func TestJsonHelperGobEncoding(t *testing.T) {
	raw := `{"testing": 123,"name":"bob & more"}`
	reader := strings.NewReader(raw)
	jh, err := NewJsonHelperReader(reader)
	assert.Tf(t, err == nil, "Unexpected error decoding gob: %s", err)
	assert.Tf(t, jh.Int("testing", 1) == 123, "Unexpected value in gob: %d", jh.Int("testing", 1))
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&jh)
	assert.T(t, err == nil, err)

	var jhNew JsonHelper
	err = gob.NewDecoder(&buf).Decode(&jhNew)
	assert.T(t, err == nil, err)
	assert.Tf(t, jhNew.Int("testing", 1) == 123, "Unexpected value in gob: %d", jhNew.Int("testing", 1))
	assert.Tf(t, jhNew.String("name", "") == "bob & more", "Unexpected value in gob: %d", jhNew.String("name", ""))

	buf2 := bytes.Buffer{}
	gt := GobTest{"Hello", jh}
	err = gob.NewEncoder(&buf2).Encode(&gt)
	assert.T(t, err == nil, err)

	var gt2 GobTest
	err = gob.NewDecoder(&buf2).Decode(&gt2)
	assert.T(t, err == nil, err)
	assert.Tf(t, gt2.Name == "Hello", "Unexpected value in gob: %d", gt2.Name)
	assert.Tf(t, gt2.Data.Int("testing", 10) == 123, "Unexpected value in gob: %d", gt2.Data.Int("testing", 10))
	assert.Tf(t, gt2.Data.String("name", "") == "bob & more", "Unexpected value in gob: %d", gt2.Data.String("name", ""))
}

type GobTest struct {
	Name string
	Data JsonHelper
}
