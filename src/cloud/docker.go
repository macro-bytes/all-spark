package cloud

import (
	"context"
	"daemon"
	"log"
	"strconv"
	"time"
	"util/netutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// DockerEnvironment interface
type DockerEnvironment struct {
	NanoCpus    int64
	MemBytes    int64
	Network     string
	ClusterID   string
	WorkerNodes int
	Image       string
}

func (e *DockerEnvironment) getDockerClient() *client.Client {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	cli.NegotiateAPIVersion(ctx)
	return cli
}

// CreateCluster - creates a spark cluster in docker
func (e *DockerEnvironment) CreateCluster() (string, error) {
	expectedWorkers := "EXPECTED_WORKERS=" + strconv.Itoa(e.WorkerNodes)
	containerID, err := e.createSparkNode(e.ClusterID+masterIdentifier, []string{expectedWorkers})
	if err != nil {
		log.Fatal(err)
	}

	masterIP, err := e.getIPAddress(containerID, e.Network)
	if err != nil {
		log.Fatal(err)
	}

	if netutil.IsListeningOnPort(masterIP, sparkPort, 30*time.Second, 120) {
		env := []string{"MASTER_IP=" + masterIP,
			"CLUSTER_ID=" + e.ClusterID,
			"CALLBACK_URL=" + daemon.GetAllSparkConfig().CallbackURL}
		for i := 1; i <= e.WorkerNodes; i++ {
			identifier := e.ClusterID + workerIdentifier + strconv.Itoa(i)
			e.createSparkNode(identifier, env)
		}
	} else {
		log.Fatal("master node has failed to come online")
	}

	webURL := "http://" + masterIP + ":8080"
	return webURL, nil
}

// DestroyCluster - destroys the spark cluster in docker
func (e *DockerEnvironment) DestroyCluster() error {
	cli := e.getDockerClient()
	defer cli.Close()

	clusterNodes, err := e.getClusterNodes()
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

func (e *DockerEnvironment) getClusterNodes() ([]string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("name", e.ClusterID)

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

func (e *DockerEnvironment) getIPAddress(id string, network string) (string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", nil
	}

	return resp.NetworkSettings.Networks[network].IPAddress, nil
}

func (e *DockerEnvironment) createSparkNode(identifier string,
	envParams []string) (string, error) {

	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: e.Image,
			Env:   envParams,
		},
		&container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: e.NanoCpus,
				Memory:   e.MemBytes,
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
