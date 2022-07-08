package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsWhitespace(t *testing.T) {
	for i := 0; i < 256; i++ {
		b := byte(i)
		if b == byte(' ') || b == byte('\t') || b == byte('\n') || b == byte('\r') {
			require.True(t, isWhitespace[b])
		} else {
			require.False(t, isWhitespace[b])
		}
	}
}

func NewTestItr(s string) *BufferedByteStreamIter {
	ctx := context.Background()
	stream := NewTestByteStream([]byte(s), 32)
	itr := NewBufferedStreamIter(stream, ctx)

	return itr
}

func TestSkipWhitespace(t *testing.T) {
	itr := NewTestItr("   a\t\t\tb\n\n\nc\r\r\rd \t\n\re ")

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte('a'), itr.Next())

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte('b'), itr.Next())

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte('c'), itr.Next())

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte('d'), itr.Next())

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte('e'), itr.Next())

	SkipWhitespace(itr)
	require.Equal(t, []byte{}, itr.Value())
	require.Equal(t, byte(0), itr.Next())
}

func TestIsNext(t *testing.T) {
	itr := NewTestItr("   a\t\t\tb\n\n\nc\r\r\rd \t\n\re ")
	require.NoError(t, IsNext(itr, byte('a')))
	require.NoError(t, IsNext(itr, byte('b')))
	require.NoError(t, IsNext(itr, byte('c')))
	require.NoError(t, IsNext(itr, byte('d')))
	require.NoError(t, IsNext(itr, byte('e')))

	itr = NewTestItr("   a")
	require.Error(t, IsNext(itr, byte(',')))
}

func TestParseUntil(t *testing.T) {
	itr := NewTestItr(` key":`)
	key, err := ParseUntil(itr, QM)
	require.NoError(t, err)
	require.Equal(t, []byte(` key"`), key)

	ch := itr.Next()
	require.Equal(t, byte(':'), ch)
}

func TestParseKey(t *testing.T) {
	itr := NewTestItr(`  "key": "value" `)
	key, err := ParseKey(itr)
	require.NoError(t, err)
	require.Equal(t, []byte(`"key"`), key)

	ch := itr.Next()
	require.Equal(t, byte(' '), ch)
}

func TestParseObject(t *testing.T) {
	tests := []struct {
		name     string
		objStr   string
		expected string
	}{
		{
			name:   "empty",
			objStr: `{}`,
		},
		{
			name:   "value string",
			objStr: `{"key":"value"}`,
		},
		{
			name:   "number value",
			objStr: `{"key":17}`,
		},
		{
			name:   "boolean value",
			objStr: `{"key":true}`,
		},
		{
			name:   "list value",
			objStr: `{"key":["string",0,true]}`,
		},
		{
			name: "object value with whitespace",
			objStr: `{
	"key": {
		"subkey": 0
	}
}`,
			expected: `{"key":{"subkey":0}}`,
		},
		{
			name: "complex object",
			objStr: `{
	"key1": {
		"subkey1": {
			"str1": "this, is a \"string\"\r\n with escaped characters",
			"str2": "special characters ]}[{",
		},
		"subkey2": {
			"num": 1,
			"bool": true
		}
	},
	"key2": [
		{"key": "val"}, [1,2,3,[56,78]], {"key": "val"}
	]
}`,
			expected: `{"key1":{"subkey1":{"str1":"this, is a \"string\"\r\n with escaped characters","str2":"special characters ]}[{",},"subkey2":{"num":1,"bool":true}},"key2":[{"key":"val"},[1,2,3,[56,78]],{"key":"val"}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			itr := NewTestItr(" \r\n\t" + test.objStr + ", ")
			res, err := ParseObject(itr)
			require.NoError(t, err)

			expected := test.objStr
			if len(test.expected) > 0 {
				expected = test.expected
			}

			require.Equal(t, expected, string(res))

			ch := itr.Next()
			require.Equal(t, byte(','), ch)
		})
	}
}

type ListWriteStream struct {
	data []byte
}

func NewListWriteStream() *ListWriteStream {
	data := make([]byte, 0, 128)
	data = append(data, OpenSB)
	return &ListWriteStream{
		data: data,
	}
}

func (lws *ListWriteStream) Add(item []byte) error {
	if len(lws.data) > 1 {
		lws.data = append(lws.data, COMMA)
	}

	lws.data = append(lws.data, item...)
	return nil
}

func (lws *ListWriteStream) Close() []byte {
	lws.data = append(lws.data, CloseSB)
	return lws.data
}

func TestParseVal(t *testing.T) {
	tests := []struct {
		name string
		str  string
	}{
		{
			name: "string",
			str:  `"this is a \"string\" }],"`,
		},
		{
			name: "number",
			str:  "1234",
		},
		{
			name: "object",
			str:  `{"key":{"subkey1":0,"subkey2":[1,2,3,4]}}`,
		},
		{
			name: "list",
			str:  `[1,2,3,4]`,
		},
		{
			name: "list of lists",
			str:  `[1,2,3,4,[56,78,[910]],[],{}]`,
		},
		{
			name: "empty list",
			str:  `[]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lws := NewListWriteStream()

			itr := NewTestItr(" \r\n\t" + test.str + ", ")
			isList, res, err := ParseVal(itr, lws.Add, None)
			require.NoError(t, err)

			resStr := string(res)
			if isList {
				resStr = string(lws.Close())
			}

			require.Equal(t, test.str, resStr)

			ch := itr.Next()
			require.Equal(t, byte(','), ch)
		})
	}
}

func TestSplitStream(t *testing.T) {
	var testStr = `{
	"string": "value",
	"number": 0,
    "boolean": true,
    "empty_object": {},
	"object": {
		"subkey1": {
			"str1": "this, is a \"string\"\r\n with escaped characters",
			"str2": "special characters ]}[{",
		},
		"subkey2": {
			"num": 1,
			"bool": true
		}
	},
    "empty_list": [],
	"list_of_strings": ["abc", "def", "ghi", "jkl", "mno", "qrs", "tuv", "wxy", "z"],
    "list_of_numbers": [0,1,2,3,4,5,6,7,8,9],
    "list_of_booleans": [true, false, true, true, false],
    "list_of_empty_objects": [{},{},{}],
    "list_of_objects": [{"key":"value"},
		{"key": "value"},
		{"key": 123, "list":[1,2,3]}
	],
    "list_of_empty_lists": [[],[],[]],
	"list_of_lists": [
		["abc", "def", "ghi", "jkl", "mno", "qrs", "tuv", "wxy", "z"],
		[0,1,2,3,4,5,6,7,8,9],
		[true, false, true, true, false],
		[
			{"key":"value"},
			{"key": "value"},
			{"key": [1,2,3]}
		],
	],
	"mixed_list": ["string",0,true,{"key": "value"},[1,2,3]]
}`
	expectedRoot := `{
	"string":"value",
	"number":0,
	"boolean":true,
	"empty_object":{},
	"object":{"subkey1":{"str1":"this, is a \"string\"\r\n with escaped characters","str2":"special characters ]}[{",},"subkey2":{"num":1,"bool":true}}
}`
	expectedListOfStrings := `"abc"
"def"
"ghi"
"jkl"
"mno"
"qrs"
"tuv"
"wxy"
"z"`
	expectedListOfBools := `true
false
true
true
false`
	expectedListOfEmptyObjects := `{}
{}
{}`
	expectedListOfEmptyLists := `[]
[]
[]`
	expectedListOfLists := `["abc","def","ghi","jkl","mno","qrs","tuv","wxy","z"]
[0,1,2,3,4,5,6,7,8,9]
[true,false,true,true,false]
[{"key":"value"},{"key":"value"},{"key":[1,2,3]}]`
	expectedListOfNumbers := `0
1
2
3
4
5
6
7
8
9`
	expectedListOfObjects := `{"key":"value"}
{"key":"value"}
{"key":123,"list":[1,2,3]}`
	expectedMixedList := `"string"
0
true
{"key":"value"}
[1,2,3]`

	tempDir, err := os.MkdirTemp("", "*")
	require.NoError(t, err)

	bs := NewTestByteStream([]byte(testStr), 256)
	err = SplitStream(context.Background(), bs, tempDir)
	require.NoError(t, err)

	requireContents(t, filepath.Join(tempDir, "root.json"), expectedRoot)
	requireContents(t, filepath.Join(tempDir, "list_of_booleans_00.jsonl"), expectedListOfBools)
	requireContents(t, filepath.Join(tempDir, "list_of_empty_objects_00.jsonl"), expectedListOfEmptyObjects)
	requireContents(t, filepath.Join(tempDir, "list_of_numbers_00.jsonl"), expectedListOfNumbers)
	requireContents(t, filepath.Join(tempDir, "list_of_strings_00.jsonl"), expectedListOfStrings)
	requireContents(t, filepath.Join(tempDir, "list_of_empty_lists_00.jsonl"), expectedListOfEmptyLists)
	requireContents(t, filepath.Join(tempDir, "list_of_lists_00.jsonl"), expectedListOfLists)
	requireContents(t, filepath.Join(tempDir, "list_of_objects_00.jsonl"), expectedListOfObjects)
	requireContents(t, filepath.Join(tempDir, "mixed_list_00.jsonl"), expectedMixedList)
}

func requireContents(t *testing.T, filename string, expectedContents string) {
	t.Run(filepath.Base(filename), func(t *testing.T) {
		data, err := os.ReadFile(filename)
		require.NoError(t, err)

		require.Equal(t, expectedContents, string(data))
	})
}
