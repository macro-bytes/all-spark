package cloud

import (
	"allspark/daemon"
	"allspark/logger"
	b64 "encoding/base64"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type imageFilter struct {
	Name   string
	Values []string
}

// AwsEnvironment interface
type AwsEnvironment struct {
	ClusterID        string
	Image            []imageFilter
	InstanceType     string
	EBSVolumeSize    int64
	SubnetID         string
	SecurityGroupIds []string
	WorkerNodes      int64
	Region           string
	IAMRole          string
	KeyName          string
	EnvParams        []string
	AssumeArn        string
	ExternalID       string
}

func (e *AwsEnvironment) getEc2Client() *ec2.EC2 {
	if len(e.AssumeArn) > 0 {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(e.Region)},
		))

		if len(e.ExternalID) > 0 {
			creds := stscreds.NewCredentials(sess, e.AssumeArn, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = aws.String(e.ExternalID)
			})
			return ec2.New(sess, &aws.Config{Credentials: creds})
		}
		creds := stscreds.NewCredentials(sess, e.AssumeArn)
		return ec2.New(sess, &aws.Config{Credentials: creds})
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.Region)},
	)

	if err != nil {
		logger.GetError().Println(err)
	}

	return ec2.New(sess)
}

func (e *AwsEnvironment) resolveAMI() (string, error) {
	cli := e.getEc2Client()

	imageFilters := make([]*ec2.Filter, len(e.Image))
	for idx, el := range e.Image {
		imageFilters[idx] = &ec2.Filter{
			Name:   aws.String(el.Name),
			Values: aws.StringSlice(el.Values),
		}
	}

	resp, err := cli.DescribeImages(
		&ec2.DescribeImagesInput{
			Filters: imageFilters,
		},
	)

	if err != nil {
		logger.GetError().Println(e)
	}

	if len(resp.Images) > 1 {
		return "", errors.New("image filters returned " +
			"more than one image; unable to resolve AMI")
	} else if len(resp.Images) == 0 {
		return "", errors.New("image filters returned " +
			"no images; unable to resolve AMI")
	}

	return *resp.Images[0].ImageId, err
}

func (e *AwsEnvironment) launchInstances(identifier string,
	instanceCount int64, userData string) (*ec2.Reservation, error) {

	cli := e.getEc2Client()
	encodedUserData := b64.StdEncoding.EncodeToString([]byte(userData))

	imageID, err := e.resolveAMI()
	if err != nil {
		return nil, err
	}

	input := &ec2.RunInstancesInput{

		ImageId:          aws.String(imageID),
		InstanceType:     aws.String(e.InstanceType),
		MinCount:         aws.Int64(instanceCount),
		MaxCount:         aws.Int64(instanceCount),
		SecurityGroupIds: aws.StringSlice(e.SecurityGroupIds),
		SubnetId:         aws.String(e.SubnetID),
		UserData:         aws.String(encodedUserData),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(e.IAMRole),
		},

		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(identifier),
					},
				},
			},
		},

		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					Encrypted:  aws.Bool(true),
					VolumeSize: aws.Int64(e.EBSVolumeSize),
					VolumeType: aws.String("gp2"),
				},
			},
		},
	}

	if len(e.KeyName) > 0 {
		input.KeyName = aws.String(e.KeyName)
	}

	resp, err := cli.RunInstances(input)

	if err != nil {
		return nil, err
	}

	for _, el := range resp.Instances {
		logger.GetInfo().Printf("launched ec2 instance %s, with identifier %s",
			*el.InstanceId, identifier)

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
	userData := "EXPECTED_WORKERS=" + workers +
		"\nSPARK_WORKER_PORT=" + strconv.FormatInt(sparkWorkerPort, 10) +
		"\nCLUSTER_ID=" + e.ClusterID +
		"\nALLSPARK_CALLBACK=" + daemon.GetAllSparkConfig().CallbackURL

	for _, el := range e.EnvParams {
		userData += "\n" + el
	}

	res, err := e.launchInstances(e.ClusterID+masterIdentifier, 1, userData)
	if err != nil {
		return "", "", err
	}

	privateIP := *res.Instances[0].PrivateIpAddress

	return *res.Instances[0].InstanceId, privateIP, err
}

func (e *AwsEnvironment) launchWorkers(masterIP string) (*ec2.Reservation, error) {

	userData := "MASTER_IP=" + masterIP +
		"\nSPARK_WORKER_PORT=" + strconv.FormatInt(sparkWorkerPort, 10)

	for _, el := range e.EnvParams {
		userData += "\n" + el
	}

	return e.launchInstances(e.ClusterID+workerIdentifier,
		e.WorkerNodes, userData)
}

// CreateCluster - creates a spark cluster in AWS
func (e *AwsEnvironment) CreateCluster() (string, error) {
	_, privateIP, err := e.launchMaster()
	if err != nil {
		return "", err
	}

	if e.WorkerNodes > 0 {
		_, err = e.launchWorkers(privateIP)
	}

	return "", err
}

// DestroyCluster - destroys a spark cluster in AWS
func (e *AwsEnvironment) DestroyCluster() error {
	cli := e.getEc2Client()
	instances, err := e.getClusterNodes()
	if err != nil {
		return err
	}
	if len(instances) > 0 {
		logger.GetInfo().Printf("destroying cluster %v with instance-ids %v", e.ClusterID, instances)

		_, err = cli.TerminateInstances(
			&ec2.TerminateInstancesInput{
				InstanceIds: aws.StringSlice(instances),
			},
		)
	} else {
		logger.GetInfo().Printf("cluster %v nas no instances and may have been terminated", e.ClusterID)
	}

	return err
}

// DestructionConfirmed - returns true if the cluster has been terminated; false otherwise
func (e *AwsEnvironment) DestructionConfirmed() bool {
	instances, err := e.getClusterNodes()
	if err != nil {
		logger.GetError().Println(err)
		logger.GetError().Printf("unable to confirm destruction of cluster %v", e.ClusterID)
		return false
	}

	return len(instances) == 0
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
					Name: aws.String("instance-state-name"),
					Values: aws.StringSlice([]string{"running", "pending",
						"shutting-down", "stopping", "stopped"}),
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
