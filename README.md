**Building all-spark**

* `git clone https://github.com/macro-bytes/all-spark`
* `cd all-spark`
* `export GOPATH=PATH/TO/ALL-SPARK`
* `./install_prerequisites.sh`
* `go build allspark`



**Example Usage**

The example below will create and destroy a spark cluster in docker, based on the configuration specified in `sample_templates/docker.json`


create-cluster:

`./allspark create-cluster --cloud-environment docker --template sample_templates/docker.json`

destroy-cluster:

`./allspark destroy-cluster --cloud-environment docker --template sample_templates/docker.json `
