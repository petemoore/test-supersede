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
		"JVV_XYi-QFCBbmQ8v6m0ew",
		"C-j4b1LERq-DFhd3uf-9-A",
		"ARbfJul8QC6M2i8lZRtCeg",
		"PXXbt0-9QmyeVLyPnmQ4Kg",
		"NQArEEbqTDaTxkwqShM_DA",
	}

	for _, taskID := range taskIDs {
		maxRunTime := 3600
		// primary task should fail, to trigger abort
		if taskID == "NQArEEbqTDaTxkwqShM_DA" {
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
	"supersederUrl": "https://gist.githubusercontent.com/petemoore/80f4ba8a8a47050a59e17a3c74a99432/raw/76d403f1533610911c26f88c8b27df78ae4314f1/supersede-test.txt"
  }`),
			ProvisionerID: "aws-provisioner-v1",
			WorkerType:    "gecko-1-b-linux",
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
