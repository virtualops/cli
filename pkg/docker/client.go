package docker

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/jhoonb/archivex"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Client struct {
	*client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	return &Client{cli}, err
}

type PathFilter func(path string) bool

// Create the docker build context and return the archive's location
func CreateBuildContext(filter PathFilter) (string, error) {
	unixTimestamp := fmt.Sprintf("%d", time.Now().Unix())
	filename := fmt.Sprintf("%x.tar", md5.Sum([]byte(unixTimestamp)))
	filename = fmt.Sprintf("/tmp/%s.tar", filename)

	tar := new(archivex.TarFile)
	err := tar.Create(filename)
	if err != nil {
		return "", err
	}

	items, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if !filter(strings.TrimPrefix(item.Name(), "/")) {
			continue
		}

		if item.IsDir() {
			err = tar.AddAll(item.Name(), true)
			if err != nil {
				return "", err
			}

			continue
		}

		fileContent, err := ioutil.ReadFile(item.Name())

		if err != nil {
			return "", err
		}
		err = tar.Add(item.Name(), bytes.NewBuffer(fileContent), item)
		if err != nil {
			return "", err
		}
	}

	return filename, tar.Close()
}

func (c *Client) BuildImage(ctx context.Context, dockerfile string, imageName string) (string, error) {
	// the tag of the image is the MD5 hash of the current timestamp
	unixTimestamp := fmt.Sprintf("%d", time.Now().Unix())
	autoTag := fmt.Sprintf("%s:%x", imageName, md5.Sum([]byte(unixTimestamp)))
	buildContextFile, err := CreateBuildContext(func(path string) bool {
		return !strings.HasPrefix(path, "vendor") &&
			!strings.HasPrefix(path, "node_modules")
	})

	if err != nil {
		return "", err
	}

	buildContext, err := os.Open(buildContextFile)

	if err != nil {
		return "", err
	}

	defer buildContext.Close()

	r, err := c.Client.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Tags:       []string{autoTag},
		Dockerfile: dockerfile,
	})

	if err != nil {
		return "", err
	}

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	defer r.Body.Close()
	//b, _ := ioutil.ReadAll(r.Body)
	//fmt.Println("build", string(b))

	return autoTag, jsonmessage.DisplayJSONMessagesStream(r.Body, os.Stderr, termFd, isTerm, nil)
}
