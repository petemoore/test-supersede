package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	tcclient "github.com/taskcluster/taskcluster-client-go"
	"github.com/taskcluster/taskcluster-client-go/queue"
)

func main() {
	creds := &tcclient.Credentials{
		ClientID:    os.Getenv("TASKCLUSTER_CLIENT_ID"),
		AccessToken: os.Getenv("TASKCLUSTER_ACCESS_TOKEN"),
		Certificate: os.Getenv("TASKCLUSTER_CERTIFICATE"),
	}

	taskIDs := []string{
		"Wb9WW8nUTTKhd2lkH8yeLw",
		"BYpgT70gSdmQSlahptyqPg",
		"DNNF4jGZRDmYYasNWnS_Og",
		"QznWlI1fTJml5M4yMqWcsw",
		"EtVkQJlqT2ipsNYJ3lRbHQ",
	}

	for _, taskID := range taskIDs {
		maxRunTime := 3600
		// primary task should fail, to trigger abort
		if taskID == "EtVkQJlqT2ipsNYJ3lRbHQ" {
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
				Description: "supercedes test",
				Name:        "supercedes test",
				Owner:       "pmoore@mozilla.com",
				Source:      "https://github.com/",
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
	"supersederUrl": "https://gist.githubusercontent.com/petemoore/80f4ba8a8a47050a59e17a3c74a99432/raw/d3c15c9ec382b035c62bd8f7ed802c5f8d06c527/supersede-test.txt"
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
