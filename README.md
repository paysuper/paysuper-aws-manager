# PaySuper AWS S3 Manager

[![License: GPL 3.0](https://img.shields.io/badge/License-GPL3.0-green.svg)](https://opensource.org/licenses/Gpl3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/paysuper/paysuper-aws-manager/issues)
[![Build Status](https://travis-ci.com/paysuper/paysuper-aws-manager.svg?branch=master)](https://travis-ci.com/paysuper/paysuper-aws-manager)
[![codecov](https://codecov.io/gh/paysuper/paysuper-aws-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-aws-manager)
[![go report](https://goreportcard.com/badge/github.com/paysuper/paysuper-aws-manager)](https://goreportcard.com/report/github.com/paysuper/paysuper-aws-manager)

PaySuper AWS S3 Manager is an AWS API wrapper.

***

## Table of Contents

- [Usage](#usage)
- [Developing](#developing)
- [Contributing](#contributing-feature-requests-and-support)
- [License](#license)

## Usage

Application handles configurations from the environment variables.

### Environment variables:

| Name                   | Required | Default   | Description                 |
|:-----------------------|:--------:|:----------|:----------------------------|
| `AWS_ACCESS_KEY_ID`      | true     | -         | AWS access key identifier   |
| `AWS_SECRET_ACCESS_KEY`  | true     | -         | AWS access secret key       |
| `AWS_BUCKET`             | true     | -         | AWS bucket name             |
| `AWS_REGION`             | -        | eu-west-1 | AWS region                  |
| `AWS_TOKEN`              | -        | ""        | AWS region                  |

### Usage example

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

## Developing

### Prerequisites

Install [Mockery](https://github.com/vektra/mockery).

### Re-Build mocks for a package

```bash
mockery -name AwsManagerInterface -output ./pkg/mocks
```

## Contributing, Feature Requests and Support

If you like this project then you can put a ⭐ on it. It means a lot to us.

If you have an idea of how to improve PaySuper (or any of the product parts) or have general feedback, you're welcome to submit a [feature request](../../issues/new?assignees=&labels=&template=feature_request.md&title=).

Chances are, you like what we have already but you may require a custom integration, a special license or something else big and specific to your needs. We're generally open to such conversations.

If you have a question and can't find the answer yourself, you can [raise an issue](../../issues/new?assignees=&labels=&template=issue--support-request.md&title=I+have+a+question+about+<this+and+that>+%5BSupport%5D) and describe what exactly you're trying to do. We'll do our best to reply in a meaningful time.

We feel that a welcoming community is important and we ask that you follow PaySuper's [Open Source Code of Conduct](https://github.com/paysuper/code-of-conduct/blob/master/README.md) in all interactions with the community.

PaySuper welcomes contributions from anyone and everyone. Please refer to [our contribution guide to learn more](CONTRIBUTING.md).

## License

The project is available as open source under the terms of the [GPL v3 License](https://www.gnu.org/licenses/gpl-3.0).
