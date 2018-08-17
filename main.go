package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	transcoder "github.com/aws/aws-sdk-go/service/elastictranscoder"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) {
	debug("Running lambda")
	pipelineID := os.Getenv("PIPELINE_ID")
	region := os.Getenv("REGION")
	presetID := os.Getenv("PRESET_ID")
	recordsCount := len(s3Event.Records)

	debug(fmt.Sprintf("PIPELINE_ID: %s", pipelineID))
	debug(fmt.Sprintf("REGION: %s", region))
	debug(fmt.Sprintf("PRESET_ID: %s", presetID))
	debug(fmt.Sprintf("RECORDS COUNT: %d", recordsCount))

	for i, record := range s3Event.Records {
		debug(fmt.Sprintf("Processing records: %d/%d", i+1, recordsCount))
		s3 := record.S3
		bucket := s3.Bucket.Name
		objKey := s3.Object.Key
		filenameWithoutExtension := objKey[0 : len(objKey)-len(filepath.Ext(objKey))]
		filenameWithExtension := filenameWithoutExtension + ".mp4"

		debug(fmt.Sprintf("Event info :: [%s - %s] BUCKET = %s, KEY = %s \n", record.EventSource, record.EventTime, bucket, objKey))

		debug("Creating session...")
		sess := session.Must(
			session.NewSession(
				&aws.Config{
					// Set supported region
					Region: aws.String(region),
				},
			),
		)

		debug("Creating transcoder...")
		svc := transcoder.New(sess)

		debug("Creating job...")
		resp, err := svc.CreateJob(
			&transcoder.CreateJobInput{
				Input: &transcoder.JobInput{
					Key: aws.String(objKey),
				},
				Outputs: []*transcoder.CreateJobOutput{
					&transcoder.CreateJobOutput{
						PresetId: aws.String(presetID),
						Key:      aws.String(filenameWithExtension),
						// If the original filename is "test-file.mp4",
						// this generates "test-file.mp4-00001.png"
						ThumbnailPattern: aws.String(fmt.Sprintf("%s-{count}", filenameWithExtension)),
					},
				},
				PipelineId: aws.String(pipelineID),
			},
		)

		if err != nil {
			debug("Job created with error")
			debug(err.Error())
			debug("Job error written to log")
		} else {
			debug("Job created without error")
		}

		debug("Record processed: %v\n", resp.Job)
	}

	debug("...Lambda COMPLETED")
}

func debug(msg string, args ...interface{}) {
	if os.Getenv("VERBOSE") == "true" {
		prefix := "TRANSCODER: "
		if args == nil {
			println(prefix + msg)
		} else {
			println(prefix+msg, args)
		}
	}
}
