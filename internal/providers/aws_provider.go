package providers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sullivtr/k8s_platform/internal/modules"
	"github.com/sullivtr/k8s_platform/internal/types"
)

// AWSProvider is a port for the aws client.
// It provides an api for interacting with aws services for the needs of khub.
type AWSProvider struct {
	Session AWSSession
}

// Compile time proof of implementation
var _ IAWSProvider = (*AWSProvider)(nil)

// AWS session represents the AWS session, and the account it belongs to
type AWSSession struct {
	SDK *modules.AWSSDK
}

// InitAWSProvider will initialize the AWSProvider implementation.
func (p *ModuleProviders) InitAWSProvider() {
	sess := session.Must(session.NewSession())

	awsSDK := modules.NewAWSSDK(sess, p.Config.AWSRegion, p.Config.ReportsBucket)

	p.AWSProvider = &AWSProvider{
		Session: AWSSession{
			SDK: awsSDK,
		},
	}
}

// GetS3Reports will fetch the list of reports from the configured S3 bucket
func (p *AWSProvider) GetS3Reports() ([]types.Report, error) {
	return p.Session.SDK.GetReports()
}

// GetReportDownloadURL will generate a presigned URL for the given report
func (p *AWSProvider) GetReportDownloadURL(reportName string) (string, error) {
	return p.Session.SDK.GetReportDownloadURL(reportName)
}
