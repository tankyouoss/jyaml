# jyaml
JSON to Yaml and vice versa conversion

## Usage
```
jyaml examples/test.yml .tmp/test.json
jyaml examples/test.json .tmp/test.yaml
jyaml examples/test.yml -j .tmp/test.noextension
```
It detects input type (JSON or Yaml)
Base on the output extension, it detects expected output type (json or yaml).
Flags:
 * -j: force json output
 * -y: force yaml output
 * -p: prettify json output
