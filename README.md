# cnab-arm-driver

Tool for generating ARM template from a CNAB bundle.json, and for invoking the bundle in ACI using the cnab-azure-driver.


## Usage

Generating the ARM template

```shell
Usage:
  cnabarmdriver generate [flags]

Flags:
  -b, --bundle string   name of bundle file to generate template for , default is bundle.json (default "bundle.json")
  -f, --file string     file name for generated template,default is azuredeploy.json (default "azuredeploy.json")
  -h, --help            help for cnabarmdriver
  -i, --indent          specifies if the json output should be indented
  -o, --overwrite       specifies if to overwrite the output file if it already exists, default is false
```

Invoking bundle  in ACI using the cnab-azure-driver

```shell
Usage:
  cnabarmdriver
```