# JSplit

JSplit is a program that can take large JSON files and split them up into a root.json files and several
[.jsonl](https://jsonlines.org) files. The program takes the list items in the root of the JSON document
and creates jsonl files containing the data from those lists.  The files representing list data take the
form [key]_%02d.jsonl where [key] is the key for the list being processed and %02d will be sequential indexes
for the files. Order of data in the lists is maintained across the files. Non-list items in the root of the JSON
document will be written to a file root.json

# Installation

To install the application you will need [Golang installed](https://go.dev/doc/install) and you will need to clone
this repository.  Once you have cloned the repository cd into the cloned jsplit directory and run:

`go install .`

# Usage

`jsplit -file <input_file> [-output <output_path>]`

  * file - (Required) Name of the json or or gz encoded json file being split into jsonl files
  * output - (Optional) Output directory. If not provided, a directory will be created based on the name of the input file.  For example, if the file myfile.json is being split and an output direce a directory named myfile\_json would be created and output would be written there.

# Example

#### example.json
```json
{
  "string": "val1",
  "number": 0,
  "bool": true,
  "object": {
    "key": "value"
  },
  "list": [
    {"idx": 0, "name":  "alex"},
    {"idx": 1, "name":  "brian"},
    {"idx": 2, "name":  "charles"},
  ]
}
```

#### Usage

`jsplit -file example.json`

### Output files

#### root.json

#### example\_json/root.json

```json
{
  "string": "val1",
  "number": 0,
  "bool": true,
  "object": {"key": "value"}
}
```

#### example\_json/list\_00.json

```json lines
{"idx": 0, "name":  "alex"}
{"idx": 1, "name":  "brian"}
{"idx": 2, "name":  "charles"}
```

In the case that a jsonl output file exceeds 4GB a new file will be created with the next sequence number. In this case
the next output file would be list\_01.jsonl
