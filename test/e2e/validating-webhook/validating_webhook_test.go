package validating_webhook_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	imageTag = "validating-webhook"
)

var _ = Describe("ValidatingWebhook", func() {

	BeforeSuite(func() {
		// delete the if it is already running
		fmt.Println("deleting existing cluster")
		cmd := exec.Command("kind", "delete", "cluster")
		err := cmd.Run()
		if err != nil {
			message := fmt.Sprintf("error while deleting cluster. err: %v", err)
			Fail(message)
		}

		// create a new cluster
		fmt.Println("creating new cluster")
		cmd = exec.Command("kind", "create", "cluster")
		err = cmd.Run()
		if err != nil {
			message := fmt.Sprintf("error while creating cluster. err: %v", err)
			Fail(message)
		}

		repoName := os.Getenv("REPO_NAME")
		if repoName == "" {
			Fail("env REPO_NAME not set")
		}

		// load terrascan image to the cluster node
		imageWithTag := fmt.Sprintf("%s:%s", repoName, imageTag)
		fmt.Println("loading docker image", imageWithTag)
		cmd = exec.Command("kind", "load", "docker-image", imageWithTag)
		err = cmd.Run()
		if err != nil {
			message := fmt.Sprintf("error while loading docker image %s to cluster. err: %v", imageWithTag, err)
			Fail(message)
		}
	})

	AfterSuite(func() {
		// delete the cluster
		cmd := exec.Command("kind", "delete", "cluster")
		err := cmd.Run()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("terrascan server is run in a pod", func() {
		It("should create pod successfully", func() {
			yamlPath, err := filepath.Abs(filepath.Join("test-data", "yamls", "terrascan-pod.yaml"))
			Expect(err).NotTo(HaveOccurred())

			fmt.Println("creating terrascan pod", yamlPath)
			cmd := exec.Command("./scripts/create-pod.sh", yamlPath)
			err = cmd.Run()
			if err != nil {
				if e, ok := err.(*exec.ExitError); ok {
					fmt.Println(e.String())
					Fail("error while creating pod")
				}
			}

			fmt.Println("wait for pod to be in running state")
			cmd = exec.Command("./scripts/wait-pod.sh")
			err = cmd.Run()
			if err != nil {
				if e, ok := err.(*exec.ExitError); ok {
					fmt.Println(e.String())
					Fail("error while creating pod")
				}
			}

			fmt.Println("getting pod logs")
			cmd = exec.Command("./scripts/get-logs.sh")
			out, err := cmd.Output()
			if err != nil {
				if e, ok := err.(*exec.ExitError); ok {
					fmt.Println(e.String())
					Fail("error while get logs from pod")
				}
			}
			fmt.Println(string(out))
		})
	})
})
