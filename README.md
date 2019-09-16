**Building all-spark**

* `git clone https://github.com/macrobytes/allspark-orchestrator.git
* `cd allspark-orchestrator
* `make install_prerequisites`
* `make all`


**Example Usage**

The example below will create and destroy a spark cluster in docker, based on the configuration specified in `dist/sample_templates/docker.json`


create-cluster:

`./allspark_cli create-cluster --cloud-environment docker --template dist/sample_templates/docker.json`

destroy-cluster:

`./allspark_cli destroy-cluster --cloud-environment docker --template dist/sample_templates/docker.json `
