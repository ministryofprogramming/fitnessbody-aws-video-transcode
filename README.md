# fitnessbody-aws-video-transcode
Aws video transcoding for fitness body.

## Prerequisites

> Install [GoLang](https://golang.org/) v1.10+  
> Install [dep](https://github.com/golang/dep) - used for managing project dependencies.  

## Installation

Clone project to your pc, preferably to `$GOPATH/src/github.com/fitnessbody-aws-video-transcode/` folder

```sh
mkdir -p $GOPATH/src/github.com/fitnessbody-aws-video-transcode
cd $GOPATH/src/github.com/fitnessbody-aws-video-transcode/
git clone https://github.com/ministryofprogramming/fitnessbody-aws-video-transcode.git
```

Install the dependencies:

```sh
cd $GOPATH/src/github.com/fitnessbody-aws-video-transcode
dep ensure -v
```

## Build, deploy (manually)

Build the app (mac/linux):

```sh
env GOOS=linux GOARCH=amd64 go build -o main     
```

Zip the code
```sh
zip -j main.zip main        
```

Upload main.zip to aws lambda function


## Additional info

### Buckets

On AWS S3 create
 - input bucket
 - output bucket
 - output thumbnails bucket

Under input bucket properties, enable "events". 

In events options select "ObjectCreate(All)" action, Sendto to "Lambda" and select transcoding lambda function (this will notify lambda function every time when new item is added to S3 bucket)

### Elastic Transcoder

On AWS ET create new pipeline.

Set input bucket, output bucket, output thumbnails bucket and other options (set them to default).

### Lambda function

Lambda function must have access to elastic transcoder.

...



## Based on 
https://dev.to/mmyoji/video-processing-with-aws-lambda--elastic-transcoder-in-golang--hf2