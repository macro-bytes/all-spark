package cloud

import (
	b64 "encoding/base64"
	"log"
	"strconv"
	"template"
	"time"
	"util/netutil"
	"util/serializer"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsEnvironment interface
type AwsEnvironment struct {
	region string
}

func (e *AwsEnvironment) getEc2Client() *ec2.EC2 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.region)},
	)

	if err != nil {
		log.Fatal(err)
	}

	return ec2.New(sess)
}

func (e *AwsEnvironment) launchInstances(template template.AwsTemplate,
	identifier string, instanceCount int64, userData string) (*ec2.Reservation, error) {

	cli := e.getEc2Client()

	encodedUserData := b64.StdEncoding.EncodeToString([]byte(userData))

	resp, err := cli.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(template.ImageId),
		InstanceType:     aws.String(template.InstanceType),
		MinCount:         aws.Int64(instanceCount),
		MaxCount:         aws.Int64(instanceCount),
		SecurityGroupIds: aws.StringSlice(template.SecurityGroupIds),
		SubnetId:         aws.String(template.SubnetId),
		UserData:         aws.String(encodedUserData),
		KeyName:          aws.String(template.KeyName),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(template.IAMRole),
		},

		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					Encrypted:  aws.Bool(true),
					VolumeSize: aws.Int64(template.EBSVolumeSize),
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

func (e *AwsEnvironment) launchMaster(template template.AwsTemplate,
	baseIdentifier string) (string, string, error) {

	workers := strconv.FormatInt(template.WorkerNodes, 10)
	userData := "export EXPECTED_WORKERS=" + workers

	res, err := e.launchInstances(template, baseIdentifier+masterIdentifier,
		1, userData)
	if err != nil {
		return "", "", err
	}

	privateIP := *res.Instances[0].PrivateIpAddress

	return *res.Instances[0].InstanceId, privateIP, err
}

func (e *AwsEnvironment) launchWorkers(template template.AwsTemplate,
	baseIdentifier string, masterIP string) (*ec2.Reservation, error) {

	userData := "export MASTER_IP=" + masterIP

	return e.launchInstances(template,
		baseIdentifier+workerIdentifier,
		template.WorkerNodes,
		userData)
}

// CreateClusterHelper - helper function for creating the spark cluster
func (e *AwsEnvironment) CreateClusterHelper(awsTemplate template.AwsTemplate) (string, error) {
	e.region = awsTemplate.Region
	baseIdentifier := buildBaseIdentifier(awsTemplate.ClusterID)
	instanceID, privateIP, err := e.launchMaster(awsTemplate, baseIdentifier)
	if err != nil {
		return "", err
	}
	_, err = e.launchWorkers(awsTemplate, baseIdentifier, privateIP)

	publicIP, err := e.getPublicIP(instanceID)
	if err != nil {
		return "", err
	}

	if netutil.IsListeningOnPort(publicIP, 8080, 1*time.Second, 60) {
		log.Println("spark master node is online")
	}

	webURL := "http://" + publicIP + ":8080"
	return webURL, err
}

// CreateCluster - deserializes the supplied template and creates a spark cluster
func (e *AwsEnvironment) CreateCluster(templatePath string) (string, error) {
	var awsTemplate template.AwsTemplate
	err := serializer.DeserializePath(templatePath, &awsTemplate)
	if err != nil {
		log.Fatal(err)
	}
	return e.CreateClusterHelper(awsTemplate)
}

// DestroyClusterHelper - helper function for destroying spark clusters
func (e *AwsEnvironment) DestroyClusterHelper(awsTemplate template.AwsTemplate) error {
	e.region = awsTemplate.Region

	identifier := awsTemplate.ClusterID

	cli := e.getEc2Client()
	instances, err := e.getClusterNodes(identifier)
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

// DestroyCluster - destroys the spark cluster
func (e *AwsEnvironment) DestroyCluster(templatePath string) error {
	var awsTemplate template.AwsTemplate
	err := serializer.DeserializePath(templatePath, &awsTemplate)
	if err != nil {
		log.Fatal(err)
	}
	return e.DestroyClusterHelper(awsTemplate)
}

func (e *AwsEnvironment) getClusterNodes(identifier string) ([]string, error) {
	var instances []string

	cli := e.getEc2Client()
	resp, err := cli.DescribeInstances(
		&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice([]string{identifier + "*"}),
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
