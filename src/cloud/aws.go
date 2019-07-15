package cloud

import (
	"container/list"
	"log"
	"template"
	"util/template_reader"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsEnvironment struct {
	instanceIDs *list.List
}

func (e *AwsEnvironment) getEc2Client() *ec2.EC2 {
	session, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	return ec2.New(session)
}

func (e *AwsEnvironment) launchInstances(template template.AwsTemplate,
	identifier string, instanceCount int64, tags []*ec2.Tag, userData *string) error {

	cli := e.getEc2Client()

	resp, err := cli.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(template.ImageId),
		InstanceType:     aws.String(template.InstanceType),
		MinCount:         aws.Int64(instanceCount),
		MaxCount:         aws.Int64(instanceCount),
		SecurityGroupIds: aws.StringSlice(template.SecurityGroupIds),
		UserData:         userData,
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					Encrypted:  aws.Bool(true),
					VolumeSize: aws.Int64(template.EBSVolumeSize),
					VolumeType: aws.String("gp2"),
				},
			},
		},
	})

	if err != nil {
		return err
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
			return err
		}
	}

	return nil
}

func (e *AwsEnvironment) launchMaster(template template.AwsTemplate,
	baseIdentifier string) (string, error) {

	tags := []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(baseIdentifier + MASTER_IDENTIFIER),
		},
	}

	return "", e.launchInstances(template, baseIdentifier+MASTER_IDENTIFIER,
		1, tags, nil)
}

func (e *AwsEnvironment) launchWorkers(template template.AwsTemplate,
	baseIdentifier string, masterUrl string) error {

	tags := []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(baseIdentifier + WORKER_IDENTIFIER),
		},
	}

	return e.launchInstances(template,
		baseIdentifier+MASTER_IDENTIFIER,
		template.WorkerNodes,
		tags,
		aws.String(masterUrl))
}

func (e *AwsEnvironment) CreateCluster(templatePath string) error {
	var awsTemplate template.AwsTemplate
	err := template_reader.Deserialize(templatePath, &awsTemplate)
	if err != nil {
		log.Fatal(err)
	}

	baseIdentifier := buildBaseIdentifier(awsTemplate.ClusterID)
	masterUrl, err := e.launchMaster(awsTemplate, baseIdentifier)
	if err != nil {
		return err
	}
	err = e.launchWorkers(awsTemplate, baseIdentifier, masterUrl)

	return err
}

func (e *AwsEnvironment) DestroyCluster(identifier string) error {
	return nil
}

func (e *AwsEnvironment) getClusterNodes(identifier string) ([]string, error) {
	return []string{}, nil
}
