package template

type AwsTemplate struct {
	ClusterID        string
	ImageId          string
	InstanceType     string
	EBSVolumeSize    int64
	SubnetId         string
	SecurityGroupIds []string
	WorkerNodes      int64
	Region           string
}
