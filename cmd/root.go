package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var wOpts writeOpts

var rootCmd = &cobra.Command{
	Use:   "jyaml [input path] [output path]",
	Short: "Convert json to yaml and yaml to json. Unless specified, output is by default in yaml",
	Long: `Convert json to yaml and yaml to json.  Unless specified, output is by default in yaml`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("json") == true && fflags.Changed("yaml") == true {
			fmt.Println("You can not specify --json and --yaml flags in the same time")
			os.Exit(1)
		}

		data, err := readFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = WriteFile(args[1], data, wOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

type writeOpts struct {
	json bool
	yaml bool
	pretty bool
}

type FileType string

const(
	JSON FileType = "json"
	YAML = "yaml"
)

func getFileType(filepath string) FileType {
	ext := strings.ToLower(path.Ext(filepath))

	switch ext {
	case ".json":
		return JSON
	default:
		return YAML
	}
}

func readFile(filepath string) (*yaml.MapSlice, error) {
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	data := yaml.MapSlice{}
	err = yaml.Unmarshal(input, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func prettyJson(in []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "\t")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func getOutputType(filepath string, opts writeOpts) FileType {
	if opts.json == true {
		return JSON
	}

	if opts.yaml == true {
		return YAML
	}

	fileType := getFileType(filepath)
	if fileType == JSON {
		return JSON
	}

	return YAML
}

func WriteFile(filepath string, data *yaml.MapSlice, opts writeOpts) error {
	var marshalOptions []yaml.EncodeOption
	outputType := getOutputType(filepath, opts)

	if outputType == JSON {
		marshalOptions = append(marshalOptions, yaml.JSON())
	}

	output, err := yaml.MarshalWithOptions(data, marshalOptions...)
	if err != nil {
		return err
	}

	if outputType == JSON && opts.pretty == true {
		output, err = prettyJson(output)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filepath, output, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&wOpts.json, "json", "j", false, "force json output")
	rootCmd.PersistentFlags().BoolVarP(&wOpts.yaml, "yaml", "y", false, "force yaml output")
	rootCmd.PersistentFlags().BoolVarP(&wOpts.pretty, "pretty", "p", false, "prettify json output")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}