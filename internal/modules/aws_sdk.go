package modules

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type AWSSDK struct {
	Session       *session.Session
	Region        string
	S3Client      *s3.S3
	ReportsBucket string
}

func NewAWSSDK(session *session.Session, region string, reportsBucket string) *AWSSDK {
	return &AWSSDK{
		Session:       session,
		Region:        region,
		S3Client:      s3.New(session, aws.NewConfig().WithRegion(region)),
		ReportsBucket: reportsBucket,
	}
}

func (sdk *AWSSDK) GetReports() ([]types.Report, error) {
	bucketObjs, err := sdk.S3Client.ListObjects(&s3.ListObjectsInput{Bucket: &sdk.ReportsBucket})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch s3 reports data : %v", err)
	}

	reports, err := sdk.toReportContract(bucketObjs.Contents)
	if err != nil {
		return nil, fmt.Errorf("unable to process s3 reports data : %v", err)
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].LastModified.After(reports[j].LastModified)
	})

	return reports, nil
}

func (sdk *AWSSDK) GetReportDownloadURL(reportName string) (string, error) {
	req, _ := sdk.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(sdk.ReportsBucket),
		Key:    aws.String(reportName),
	})
	urlStr, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("unable to generate presigned url for %s : %v", reportName, err)
	}

	return urlStr, nil
}

// toReportContract transforms the list of s3.Objects into the user friendly Report model
func (sdk *AWSSDK) toReportContract(objects []*s3.Object) ([]types.Report, error) {
	dumpFileObjects := make([]types.Report, len(objects))
	for _, obj := range objects {
		if obj.Key != nil && *obj.Key != "" {
			expiresDate := obj.LastModified.Add(time.Hour * 24 * 90) // Objects expire in 90 days from create date
			var objSize int64
			var sizeUnits string

			if *obj.Size > 1000000 {
				objSize = *obj.Size / 1000000
				sizeUnits = "MB"
			} else {
				objSize = *obj.Size / 1000
				sizeUnits = "KB"
			}

			objTagOutput, err := sdk.S3Client.GetObjectTagging(&s3.GetObjectTaggingInput{
				Bucket: aws.String(sdk.ReportsBucket),
				Key:    obj.Key,
			})
			if err != nil {
				return dumpFileObjects, fmt.Errorf("unable to fetch object tags for %s, error: %v", *obj.Key, err)
			}

			userTag := getBucketTagByKey("User", objTagOutput)
			reasonTag := getBucketTagByKey("Reason", objTagOutput)
			typeTag := getBucketTagByKey("Type", objTagOutput)

			dumpFileObjects = append(dumpFileObjects, types.Report{
				Name:         *obj.Key,
				LastModified: *obj.LastModified,
				Created:      obj.LastModified.Format(time.UnixDate),
				Expires:      expiresDate.Format(time.UnixDate),
				Bucket:       sdk.ReportsBucket,
				Size:         objSize,
				SizeUnits:    sizeUnits,
				User:         strings.ReplaceAll(userTag, "@gmail.com", ""),
				Reason:       reasonTag,
				Type:         typeTag,
			})
		}
	}
	return dumpFileObjects, nil
}

func getBucketTagByKey(key string, objTagOutput *s3.GetObjectTaggingOutput) string {
	for _, t := range objTagOutput.TagSet {
		if strings.EqualFold(*t.Key, key) {
			return *t.Value
		}
	}
	return "unspecified"
}
