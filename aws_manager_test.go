package aws_manager

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/paysuper/paysuper-aws-manager/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"os"
	"syscall"
	"testing"
)

var (
	filePath = "./test/test_pdf.pdf"
	fileName = "test_pdf.pdf"
)

type AwsManagerTestSuite struct {
	suite.Suite
	awsManager AwsManagerInterface
}

func Test_AwsManager(t *testing.T) {
	suite.Run(t, new(AwsManagerTestSuite))
}

func (suite *AwsManagerTestSuite) SetupTest() {
	manager, err := New()

	if err != nil {
		assert.FailNow(suite.T(), "New aws manager instance init failed", "%v", err)
	}

	assert.NotNil(suite.T(), manager)

	m, ok := manager.(*AwsManager)
	assert.True(suite.T(), ok)

	assert.NotNil(suite.T(), m.awsUploader)
	assert.NotNil(suite.T(), m.awsDownloader)
	assert.NotNil(suite.T(), m.cfg)
	assert.NotEmpty(suite.T(), m.cfg.Bucket)
	assert.NotEmpty(suite.T(), m.cfg.AccessKeyId)
	assert.NotEmpty(suite.T(), m.cfg.SecretAccessKey)
	assert.NotEmpty(suite.T(), m.cfg.Region)

	mockUploader := &test.UploaderAPI{}
	mockUploader.On("UploadWithContext", mock.Anything, mock.Anything, mock.Anything).
		Return(&s3manager.UploadOutput{}, nil)
	mockDownloader := &test.DownloaderAPI{}
	mockDownloader.On("DownloadWithContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(int64(0), nil)

	m.awsUploader = mockUploader
	m.awsDownloader = mockDownloader

	suite.awsManager = m
}

func (suite *AwsManagerTestSuite) TearDownTest() {}

func (suite *AwsManagerTestSuite) TestAwsManager_Upload_WithoutBucket_Ok() {
	file, err := os.Open(filePath)
	assert.NoError(suite.T(), err)
	defer file.Close()

	in := &UploadInput{
		Body:     file,
		FileName: fileName,
	}
	_, err = suite.awsManager.Upload(context.TODO(), in)
	assert.NoError(suite.T(), err)
}

func (suite *AwsManagerTestSuite) TestAwsManager_Upload_WithBucketAndFilePath_Ok() {
	in := &UploadInput{
		Bucket:   "bucket-name",
		Path:     filePath,
		FileName: fileName,
	}
	_, err := suite.awsManager.Upload(context.TODO(), in)
	assert.NoError(suite.T(), err)
}

func (suite *AwsManagerTestSuite) TestAwsManager_Upload_FileNotFound_Error() {
	in := &UploadInput{
		Path:     "./not_exist_file_name.pdf",
		FileName: fileName,
	}
	_, err := suite.awsManager.Upload(context.TODO(), in)
	assert.Error(suite.T(), err)
	assert.Regexp(suite.T(), syscall.ENOENT, err.Error())
}

func (suite *AwsManagerTestSuite) TestAwsManager_Download_WithoutBucket_Ok() {
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	in := &DownloadInput{
		FileName: fileName,
	}
	_, err := suite.awsManager.Download(context.TODO(), filePath, in)
	assert.NoError(suite.T(), err)
}

func (suite *AwsManagerTestSuite) TestAwsManager_Download_WithBucket_Ok() {
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	in := &DownloadInput{
		Bucket:   "bucket-name",
		FileName: fileName,
	}
	_, err := suite.awsManager.Download(context.TODO(), filePath, in)
	assert.NoError(suite.T(), err)
}

func (suite *AwsManagerTestSuite) TestAwsManager_Download_CreateFileError() {
	in := &DownloadInput{
		Bucket:   "bucket-name",
		FileName: fileName,
	}
	_, err := suite.awsManager.Download(context.TODO(), "", in)
	assert.Error(suite.T(), err)
	assert.Regexp(suite.T(), syscall.ENOENT, err.Error())
}

func (suite *AwsManagerTestSuite) TestAwsManager_NewManager_WithOptions_Ok() {
	opts := []Option{
		AccessKeyId("AccessKeyId"),
		SecretAccessKey("SecretAccessKey"),
		Region("Region"),
		Bucket("Bucket"),
		Token("Token"),
	}
	manager, err := New(opts...)
	assert.NoError(suite.T(), err)

	m, ok := manager.(*AwsManager)
	assert.True(suite.T(), ok)

	assert.NotNil(suite.T(), m.awsUploader)
	assert.NotNil(suite.T(), m.awsDownloader)
	assert.NotNil(suite.T(), m.cfg)
	assert.NotEmpty(suite.T(), m.cfg.Bucket)
	assert.NotEmpty(suite.T(), m.cfg.AccessKeyId)
	assert.NotEmpty(suite.T(), m.cfg.SecretAccessKey)
	assert.NotEmpty(suite.T(), m.cfg.Region)
	assert.NotEmpty(suite.T(), m.cfg.Token)
}

func (suite *AwsManagerTestSuite) TestAwsManager_NewManager_RequiredEnvVariableNotExist_Error() {
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	err := os.Unsetenv("AWS_ACCESS_KEY_ID")
	assert.NoError(suite.T(), err)

	manager, err := New()
	assert.Error(suite.T(), err)
	assert.Regexp(suite.T(), "AWS_ACCESS_KEY_ID", err.Error())
	assert.Nil(suite.T(), manager)

	err = os.Setenv("AWS_ACCESS_KEY_ID", accessKeyId)
	assert.NoError(suite.T(), err)
}

func (suite *AwsManagerTestSuite) TestAwsManager_NewManager_NewAwsSessionError() {

}
