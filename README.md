PaySuper AWS S3 manager wrapper
=====

[![Build Status](https://travis-ci.org/paysuper/paysuper-aws-manager.svg?branch=master)](https://travis-ci.org/paysuper/paysuper-aws-manager) 
[![codecov](https://codecov.io/gh/paysuper/paysuper-aws-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-aws-manager)
[![go report](https://goreportcard.com/badge/github.com/paysuper/paysuper-aws-manager)](https://goreportcard.com/report/github.com/paysuper/paysuper-aws-manager)

## Environment variables:

| Name                   | Required | Default   | Description                 |
|:-----------------------|:--------:|:----------|:----------------------------|
| AWS_ACCESS_KEY_ID      | true     | -         | AWS access key identifier   |
| AWS_SECRET_ACCESS_KEY  | true     | -         | AWS access secret key       |
| AWS_BUCKET             | true     | -         | AWS bucket name             |
| AWS_REGION             | -        | eu-west-1 | AWS region                  |
| AWS_TOKEN              | -        | ""        | AWS region                  |

## Usage example

```go
package main

import (
    "context"
    awsWrapper "github.com/paysuper/paysuper-aws-manager"
    "log"
    "os"
)

func main() {
    awsManager, err := awsWrapper.New()
    
    if err != nil {
        log.Fatalln(err)
    }
    
    //upload open file
    file, err := os.Open("/tmp/file.pdf")
    defer file.Close()

    out := &awsWrapper.UploadInput{
        Body:     file,
        FileName: "file.pdf",
    }
    _, err = awsManager.Upload(context.TODO(), out)

    if err != nil {
        log.Fatalln(err)    
    }
    
    log.Println("file upload successfully")

    //upload file by path
    uploadReq := &awsWrapper.UploadInput{
        Path:     "/tmp/file.pdf",
        FileName: "file.pdf",
    }
    _, err = awsManager.Upload(context.TODO(), uploadReq)

    if err != nil {
        log.Fatalln(err)    
    }
    
    log.Println("file upload successfully")

    // download file
    filePath := os.TempDir() + string(os.PathSeparator) + "file.pdf"
    downloadReq := &awsWrapper.DownloadInput{
        FileName: "file.pdf",
    }
    _, err = awsManager.Download(context.TODO(), filePath, downloadReq)
    
    if err != nil {
        log.Fatalln(err)    
    }
    
    log.Println("file download successfully")
}
```

## Re-Build mocks for package

**REQUIRED:** package [mockery](https://github.com/vektra/mockery) must be installed before run next command

```bash
mockery -name AwsManagerInterface -output ./pkg/mocks
```  

## Contributing
We feel that a welcoming community is important and we ask that you follow PaySuper's [Open Source Code of Conduct](https://github.com/paysuper/code-of-conduct/blob/master/README.md) in all interactions with the community.

PaySuper welcomes contributions from anyone and everyone. Please refer to each project's style and contribution guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

The master branch of this repository contains the latest stable release of this component.

 
