package cloud

import (
	"daemon"
	b64 "encoding/base64"
	"logger"
	"strconv"
	"time"
	"util/netutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsEnvironment interface
type AwsEnvironment struct {
	ClusterID        string
	ImageID          string
	InstanceType     string
	EBSVolumeSize    int64
	SubnetID         string
	SecurityGroupIds []string
	WorkerNodes      int64
	Region           string
	IAMRole          string
	KeyName          string
	MetaData         string
}

func (e *AwsEnvironment) getEc2Client() *ec2.EC2 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.Region)},
	)

	if err != nil {
		logger.GetFatal().Fatalln(err)
	}

	return ec2.New(sess)
}

func (e *AwsEnvironment) launchInstances(identifier string,
	instanceCount int64, userData string) (*ec2.Reservation, error) {

	cli := e.getEc2Client()

	encodedUserData := b64.StdEncoding.EncodeToString([]byte(userData))

	resp, err := cli.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(e.ImageID),
		InstanceType:     aws.String(e.InstanceType),
		MinCount:         aws.Int64(instanceCount),
		MaxCount:         aws.Int64(instanceCount),
		SecurityGroupIds: aws.StringSlice(e.SecurityGroupIds),
		SubnetId:         aws.String(e.SubnetID),
		UserData:         aws.String(encodedUserData),
		KeyName:          aws.String(e.KeyName),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(e.IAMRole),
		},

		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					Encrypted:  aws.Bool(true),
					VolumeSize: aws.Int64(e.EBSVolumeSize),
					VolumeType: aws.String("gp2"),
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	for _, el := range resp.Instances {
		_, err := cli.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{el.InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(identifier),
				},
			},
		})

		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (e *AwsEnvironment) getPublicIP(instanceID string) (string, error) {
	cli := e.getEc2Client()

	cli.WaitUntilInstanceRunning(
		&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instanceID}),
		},
	)

	response, err := cli.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instanceID}),
		},
	)

	if err != nil {
		return "", err
	}

	return *response.Reservations[0].Instances[0].PublicIpAddress, nil
}

func (e *AwsEnvironment) launchMaster() (string, string, error) {

	workers := strconv.FormatInt(e.WorkerNodes, 10)
	userData := "export EXPECTED_WORKERS=" + workers +
		"\nexport CLUSTER_ID=" + e.ClusterID +
		"\nexport META_DATA=" + e.MetaData

	if len(daemon.GetAllSparkConfig().CallbackURL) > 0 {
		userData += "\nexport CALLBACK_URL=" + daemon.GetAllSparkConfig().CallbackURL
	}

	res, err := e.launchInstances(e.ClusterID+masterIdentifier, 1, userData)
	if err != nil {
		return "", "", err
	}

	privateIP := *res.Instances[0].PrivateIpAddress

	return *res.Instances[0].InstanceId, privateIP, err
}

func (e *AwsEnvironment) launchWorkers(masterIP string) (*ec2.Reservation, error) {

	userData := "export MASTER_IP=" + masterIP

	return e.launchInstances(e.ClusterID+workerIdentifier,
		e.WorkerNodes, userData)
}

// CreateCluster - creates a spark cluster in AWS
func (e *AwsEnvironment) CreateCluster() (string, error) {
	instanceID, privateIP, err := e.launchMaster()
	if err != nil {
		return "", err
	}
	_, err = e.launchWorkers(privateIP)

	publicIP, err := e.getPublicIP(instanceID)
	if err != nil {
		return "", err
	}

	if netutil.IsListeningOnPort(publicIP, 8080, 1*time.Second, 60) {
		logger.GetInfo().Println("spark master node is online")
	}

	webURL := "http://" + publicIP + ":8080"
	return webURL, err
}

// DestroyCluster - destroys a spark cluster in AWS
func (e *AwsEnvironment) DestroyCluster() error {
	cli := e.getEc2Client()
	instances, err := e.getClusterNodes()
	if err != nil {
		return err
	}

	_, err = cli.TerminateInstances(
		&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice(instances),
		},
	)
	return err
}

func (e *AwsEnvironment) getClusterNodes() ([]string, error) {
	var instances []string

	cli := e.getEc2Client()
	resp, err := cli.DescribeInstances(
		&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice([]string{e.ClusterID + "*"}),
				},
				{
					Name:   aws.String("network-interface.subnet-id"),
					Values: aws.StringSlice([]string{e.SubnetID}),
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: aws.StringSlice([]string{"running", "pending"}),
				},
			},
		},
	)

	if err != nil {
		return instances, err
	}

	for _, reservation := range resp.Reservations {
		for _, el := range reservation.Instances {
			instances = append(instances, *el.InstanceId)
		}
	}

	return instances, nil
}
