# Fork of dolthub/jsplit
Reworked somewhat for use in data pipelines.
- Input and Output paths can be URIs to AWS S3 or Google Cloud Storage objects thanks to the [Google CDK](https://gocloud.dev/howto/blob/).
- A `Split` function that wraps what was previously in `main()`, more readily allowing `split` to be used as a module in other apps.
- Dockerfile to generate a lightweight container
- Makefile with build, test, container build, and container deploy targets.

## TODO:
- Investigate heavy GC activity and mitigation. 
- When calling `Split`, potentially throwing `SplitStream` into a goroutine and returning `ctx`, allowing the calling code to cancel if necessary.

## Performance
`jsplit` currently makes heavy usage of the GC. If you notice low core utilization, trading off memory usage against the GC can be done by setting a `GOGC` value far higher than the default of 200. 
  
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

`jsplit -file <input_file> -output <output_path>`

  * file - (Required) Name of the json or or gz encoded json file being split into jsonl files
  * output - (Required) Output directory. 

Input and output can be either local filesystem paths or AWS S3 or Google Cloud Storage URIs. 

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

#### Example Usage

`jsplit -file example.json` -output s3://mybucket/example_json/

### Example Output files

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
