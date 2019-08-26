package aws_manager

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/kelseyhightower/envconfig"
	"io"
	"os"
	"time"
)

type AwsManagerInterface interface {
	Upload(context.Context, *UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	Download(context.Context, string, *DownloadInput, ...func(*s3manager.Downloader)) (int64, error)
}

type AwsManager struct {
	cfg           *Options
	awsUploader   s3manageriface.UploaderAPI
	awsDownloader s3manageriface.DownloaderAPI
}

type UploadInput struct {
	ACL                       string
	Body                      io.Reader
	Path                      string
	Bucket                    string
	CacheControl              string
	ContentDisposition        string
	ContentEncoding           string
	ContentLanguage           string
	ContentMD5                string
	ContentType               string
	Expires                   time.Time
	GrantFullControl          string
	GrantRead                 string
	GrantReadACP              string
	GrantWriteACP             string
	FileName                  string
	Metadata                  map[string]string
	ObjectLockLegalHoldStatus string
	ObjectLockMode            string
	ObjectLockRetainUntilDate time.Time
	RequestPayer              string
	SSECustomerAlgorithm      string
	SSECustomerKey            string
	SSECustomerKeyMD5         string
	SSEKMSEncryptionContext   string
	SSEKMSKeyId               string
	ServerSideEncryption      string
	StorageClass              string
	Tagging                   string
	WebsiteRedirectLocation   string
}

type DownloadInput struct {
	Bucket                     string
	IfMatch                    string
	IfModifiedSince            time.Time
	IfNoneMatch                string
	IfUnmodifiedSince          time.Time
	FileName                   string
	PartNumber                 int64
	Range                      string
	RequestPayer               string
	ResponseCacheControl       string
	ResponseContentDisposition string
	ResponseContentEncoding    string
	ResponseContentLanguage    string
	ResponseContentType        string
	ResponseExpires            time.Time
	SSECustomerAlgorithm       string
	SSECustomerKey             string
	SSECustomerKeyMD5          string
	VersionId                  string
}

type Options struct {
	AccessKeyId     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	Region          string `envconfig:"AWS_REGION" default:"eu-west-1"`
	Bucket          string `envconfig:"AWS_BUCKET" required:"true"`
	Token           string `envconfig:"AWS_TOKEN" default:""`
}

type Option func(*Options)

func AccessKeyId(accessKeyId string) Option {
	return func(opts *Options) {
		opts.AccessKeyId = accessKeyId
	}
}

func SecretAccessKey(secretAccessKey string) Option {
	return func(opts *Options) {
		opts.SecretAccessKey = secretAccessKey
	}
}

func Region(region string) Option {
	return func(opts *Options) {
		opts.Region = region
	}
}

func Bucket(bucket string) Option {
	return func(opts *Options) {
		opts.Bucket = bucket
	}
}

func Token(token string) Option {
	return func(opts *Options) {
		opts.Token = token
	}
}

func New(options ...Option) (AwsManagerInterface, error) {
	opts := Options{}
	conn := &Options{}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.HasEmptySettings() {
		err := envconfig.Process("", conn)

		if err != nil {
			return nil, err
		}
	}

	if opts.AccessKeyId != "" {
		conn.AccessKeyId = opts.AccessKeyId
	}

	if opts.SecretAccessKey != "" {
		conn.SecretAccessKey = opts.SecretAccessKey
	}

	if opts.Region != "" {
		conn.Region = opts.Region
	}

	if opts.Bucket != "" {
		conn.Bucket = opts.Bucket
	}

	if opts.Token != "" {
		conn.Token = opts.Token
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(conn.Region),
			Credentials: credentials.NewStaticCredentials(
				conn.AccessKeyId,
				conn.SecretAccessKey,
				conn.Token,
			),
		},
	)

	if err != nil {
		return nil, err
	}

	manager := &AwsManager{
		cfg:           conn,
		awsUploader:   s3manager.NewUploader(sess),
		awsDownloader: s3manager.NewDownloader(sess),
	}

	return manager, nil
}

func (m *AwsManager) Upload(
	ctx context.Context,
	in *UploadInput,
	opts ...func(*s3manager.Uploader),
) (*s3manager.UploadOutput, error) {
	if in.Body == nil && in.Path != "" {
		file, err := os.Open(in.Path)

		if err != nil {
			return nil, err
		}

		in.Body = file
		defer file.Close()
	}

	if in.Bucket == "" {
		in.Bucket = m.cfg.Bucket
	}

	s3In := in.toAwsUploadInput()
	return m.awsUploader.UploadWithContext(ctx, s3In)
}

func (m *AwsManager) Download(
	ctx context.Context,
	path string,
	in *DownloadInput,
	opts ...func(*s3manager.Downloader),
) (int64, error) {
	file, err := os.Create(path)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	if in.Bucket == "" {
		in.Bucket = m.cfg.Bucket
	}

	s3In := in.toAwsGetObjectInput()
	return m.awsDownloader.DownloadWithContext(ctx, file, s3In)
}

func (m *UploadInput) toAwsUploadInput() *s3manager.UploadInput {
	return &s3manager.UploadInput{
		ACL:                       aws.String(m.ACL),
		Body:                      m.Body,
		Bucket:                    aws.String(m.Bucket),
		CacheControl:              aws.String(m.CacheControl),
		ContentDisposition:        aws.String(m.ContentDisposition),
		ContentEncoding:           aws.String(m.ContentEncoding),
		ContentLanguage:           aws.String(m.ContentLanguage),
		ContentMD5:                aws.String(m.ContentMD5),
		ContentType:               aws.String(m.ContentType),
		Expires:                   aws.Time(m.Expires),
		GrantFullControl:          aws.String(m.GrantFullControl),
		GrantRead:                 aws.String(m.GrantRead),
		GrantReadACP:              aws.String(m.GrantReadACP),
		GrantWriteACP:             aws.String(m.GrantWriteACP),
		Key:                       aws.String(m.FileName),
		Metadata:                  aws.StringMap(m.Metadata),
		ObjectLockLegalHoldStatus: aws.String(m.ObjectLockLegalHoldStatus),
		ObjectLockMode:            aws.String(m.ObjectLockMode),
		ObjectLockRetainUntilDate: aws.Time(m.ObjectLockRetainUntilDate),
		RequestPayer:              aws.String(m.RequestPayer),
		SSECustomerAlgorithm:      aws.String(m.SSECustomerAlgorithm),
		SSECustomerKey:            aws.String(m.SSECustomerKey),
		SSECustomerKeyMD5:         aws.String(m.SSECustomerKeyMD5),
		SSEKMSEncryptionContext:   aws.String(m.SSEKMSEncryptionContext),
		SSEKMSKeyId:               aws.String(m.SSEKMSKeyId),
		ServerSideEncryption:      aws.String(m.ServerSideEncryption),
		StorageClass:              aws.String(m.StorageClass),
		Tagging:                   aws.String(m.Tagging),
		WebsiteRedirectLocation:   aws.String(m.WebsiteRedirectLocation),
	}
}

func (m *DownloadInput) toAwsGetObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket:                     aws.String(m.Bucket),
		IfMatch:                    aws.String(m.IfMatch),
		IfModifiedSince:            aws.Time(m.IfModifiedSince),
		IfNoneMatch:                aws.String(m.IfNoneMatch),
		IfUnmodifiedSince:          aws.Time(m.IfUnmodifiedSince),
		Key:                        aws.String(m.FileName),
		PartNumber:                 aws.Int64(m.PartNumber),
		Range:                      aws.String(m.Range),
		RequestPayer:               aws.String(m.RequestPayer),
		ResponseCacheControl:       aws.String(m.ResponseCacheControl),
		ResponseContentDisposition: aws.String(m.ResponseContentDisposition),
		ResponseContentEncoding:    aws.String(m.ResponseContentEncoding),
		ResponseContentLanguage:    aws.String(m.ResponseContentLanguage),
		ResponseContentType:        aws.String(m.ResponseContentType),
		ResponseExpires:            aws.Time(m.ResponseExpires),
		SSECustomerAlgorithm:       aws.String(m.SSECustomerAlgorithm),
		SSECustomerKey:             aws.String(m.SSECustomerKey),
		SSECustomerKeyMD5:          aws.String(m.SSECustomerKeyMD5),
		VersionId:                  aws.String(m.VersionId),
	}
}

func (opts *Options) HasEmptySettings() bool {
	return opts.AccessKeyId == "" || opts.SecretAccessKey == "" || opts.Region == "" || opts.Bucket == ""
}
