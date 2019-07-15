package cloud

import (
	"context"
	"log"
	"strconv"
	"template"
	"time"
	"util/netutil"
	"util/template_reader"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func (e *DockerEnvironment) getDockerClient() *client.Client {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	return cli
}

type DockerEnvironment struct{}

func (e *DockerEnvironment) CreateCluster(templatePath string) error {
	var dockerTemplate template.DockerTemplate
	err := template_reader.Deserialize(templatePath, &dockerTemplate)
	if err != nil {
		log.Fatal(err)
	}

	baseIdentifier := buildBaseIdentifier(dockerTemplate.ClusterID)

	containerID, err := e.createSparkNode(dockerTemplate,
		baseIdentifier+MASTER_IDENTIFIER, []string{})

	masterIP, err := e.getIpAddress(containerID, dockerTemplate.Network)
	if err != nil {
		log.Fatal(err)
	}

	masterURL := "MASTER_URL=spark://" + masterIP + ":" + SPARK_PORT
	log.Println("spark master URL is: " + masterURL)

	if netutil.IsListeningOnPort(masterIP, SPARK_PORT, 30*time.Second, 120) {
		env := []string{"MASTER_URL=spark://" + masterIP + ":" + SPARK_PORT}
		for i := 1; i <= dockerTemplate.WorkerNodes; i++ {
			identifier := baseIdentifier + WORKER_IDENTIFIER + strconv.Itoa(i)
			log.Println("createing spark worker " + identifier)
			e.createSparkNode(dockerTemplate, identifier, env)
		}
	} else {
		log.Fatal("master node has failed to come online")
	}

	log.Println("spark cluster is online, and can be accessed via http://" +
		masterIP + ":8080")
	return nil
}

func (e *DockerEnvironment) DestroyCluster(identifier string) error {
	cli := e.getDockerClient()
	defer cli.Close()

	clusterNodes, err := e.getClusterNodes(identifier)
	if err != nil {
		return err
	}

	for _, el := range clusterNodes {
		err = cli.ContainerRemove(context.Background(), el,
			types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *DockerEnvironment) getClusterNodes(identifier string) ([]string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("name", identifier)

	resp, err := cli.ContainerList(context.Background(),
		types.ContainerListOptions{Filters: filters})

	if err != nil {
		return nil, err
	}

	var result []string
	for _, el := range resp {
		result = append(result, el.Names[0])
	}
	return result, nil
}

func (e *DockerEnvironment) getIpAddress(id string, network string) (string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", nil
	}

	return resp.NetworkSettings.Networks[network].IPAddress, nil
}

func (e *DockerEnvironment) createSparkNode(dockerTemplate template.DockerTemplate,
	identifier string,
	envParams []string) (string, error) {

	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: SPARK_BASE_IMAGE,
			Env:   envParams,
		},
		&container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: dockerTemplate.NanoCpus,
				Memory:   dockerTemplate.MemBytes,
			},
			NetworkMode: "all-spark-bridge",
		},
		&network.NetworkingConfig{},
		identifier)
	if err != nil {
		return "", err
	}

	if err = cli.ContainerStart(context.Background(),
		resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
