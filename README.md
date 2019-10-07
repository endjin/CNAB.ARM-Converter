# CNAB.ARM-Converter - Azure ARM Template from CNAB Bundle

Generates an ARM template from a CNAB bundle.json. Uses [duffle-aci-docker](https://github.com/simongdavies/silver-garbanzo/tree/master/client/duffle-aci-docker) container image to run the bundle in Azure ACI.

## Usage

```shell
Usage:
  atfcnab [flags]

Flags:
  -b, --bundle string   name of bundle file to generate template for , default is bundle.json (default "bundle.json")
  -f, --file string     file name for generated template,default is azuredeploy.json (default "azuredeploy.json")
  -h, --help            help for atfcnab
  -i, --indent          specifies if the json output should be indented
  -o, --overwrite       specifies if to overwrite the output file if it already exists, default is false
  ```
