package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tcclient "github.com/taskcluster/taskcluster-client-go"
	"github.com/taskcluster/taskcluster-client-go/queue"
)

type Supersedes struct {
	Supersedes []string `json:"supersedes"`
}

func main() {

	sourceURL := os.Args[1]
	supersederURL := os.Args[2]
	owner := os.Args[3]

	creds := &tcclient.Credentials{
		ClientID:    os.Getenv("TASKCLUSTER_CLIENT_ID"),
		AccessToken: os.Getenv("TASKCLUSTER_ACCESS_TOKEN"),
		Certificate: os.Getenv("TASKCLUSTER_CERTIFICATE"),
	}

	resp, err := http.Get(supersederURL)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(resp.Body)
	var supersedes Supersedes
	err = decoder.Decode(&supersedes)
	if err != nil {
		panic(err)
	}
	taskIDs := supersedes.Supersedes

	for _, taskID := range taskIDs {
		maxRunTime := 3600
		// primary task should fail, to trigger abort
		if taskID == taskIDs[len(taskIDs)-1] {
			maxRunTime = 10
		}
		created := time.Now()
		deadline := created.Add(time.Hour * 1)
		expires := created.AddDate(0, 0, 1)
		tdr := &queue.TaskDefinitionRequest{
			Created:  tcclient.Time(created),
			Deadline: tcclient.Time(deadline),
			Expires:  tcclient.Time(expires),
			Metadata: struct {
				Description string `json:"description"`
				Name        string `json:"name"`
				Owner       string `json:"owner"`
				Source      string `json:"source"`
			}{
				Description: "supersedes test",
				Name:        "supersedes test",
				Owner:       owner,
				Source:      sourceURL,
			},
			Payload: json.RawMessage(`{
    "maxRunTime": ` + strconv.Itoa(maxRunTime) + `,
    "image": {
      "path": "public/image.tar.zst",
      "type": "task-image",
      "taskId": "Pr9OcxSqQlOjbytRDpHd2g"
    },
    "command": [
      "sleep",
      "60"
    ],
	"supersederUrl": "` + supersederURL + `"
  }`),
			ProvisionerID: "aws-provisioner-v1",
			WorkerType:    "tutorial",
		}
		q := queue.New(creds)
		log.Printf("Creating task %v...", taskID)
		_, err := q.CreateTask(taskID, tdr)
		if err != nil {
			panic(err)
		}
	}
	log.Print("Done!")
}
