package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	oss "github.com/alibabacloud-go/oss-20190517/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// OSSProvider wraps the Aliyun OSS SDK client.
type OSSProvider struct {
	client *oss.Client
	region string
}

// NewOSSProvider creates a new OSSProvider with the given access key and region.
func NewOSSProvider(accessKeyId, accessKeySecret, region string) (*OSSProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("oss-%s.aliyuncs.com", region))

	client, err := oss.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %w", err)
	}
	return &OSSProvider{client: client, region: region}, nil
}

// ListBuckets lists all OSS buckets.
func (p *OSSProvider) ListBuckets() ([]provider.OSSBucket, error) {
	request := &oss.ListBucketsRequest{
		MaxKeys: tea.Int64(100),
	}

	response, err := p.client.ListBuckets(request)
	if err != nil {
		return nil, fmt.Errorf("ListBuckets failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return []provider.OSSBucket{}, nil
	}

	buckets := make([]provider.OSSBucket, 0)
	if body.Buckets != nil && body.Buckets.Buckets != nil {
		for _, b := range body.Buckets.Buckets {
			bucket := provider.OSSBucket{
				Name:         tea.StringValue(b.Name),
				Location:     tea.StringValue(b.Location),
				CreationDate: tea.StringValue(b.CreationDate),
				StorageClass: tea.StringValue(b.StorageClass),
			}
			if b.ExtranetEndpoint != nil {
				bucket.ExtranetEndpoint = tea.StringValue(b.ExtranetEndpoint)
			}
			if b.IntranetEndpoint != nil {
				bucket.IntranetEndpoint = tea.StringValue(b.IntranetEndpoint)
			}
			buckets = append(buckets, bucket)
		}
	}

	return buckets, nil
}

// ListObjects lists objects in an OSS bucket.
func (p *OSSProvider) ListObjects(bucket, prefix string, maxKeys int32) ([]provider.OSSObject, bool, error) {
	if maxKeys > 1000 {
		maxKeys = 1000
	}
	if maxKeys < 1 {
		maxKeys = 100
	}
	request := &oss.ListObjectsV2Request{
		MaxKeys: tea.Int64(int64(maxKeys)),
	}
	if prefix != "" {
		request.Prefix = tea.String(prefix)
	}

	response, err := p.client.ListObjectsV2(tea.String(bucket), request)
	if err != nil {
		return nil, false, fmt.Errorf("ListObjects failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return []provider.OSSObject{}, false, nil
	}

	objects := make([]provider.OSSObject, 0)
	if body.Contents != nil {
		for _, obj := range body.Contents {
			object := provider.OSSObject{
				Key:          tea.StringValue(obj.Key),
				Size:         tea.Int64Value(obj.Size),
				LastModified: tea.StringValue(obj.LastModified),
				ETag:         tea.StringValue(obj.ETag),
				Type:         tea.StringValue(obj.Type),
				StorageClass: tea.StringValue(obj.StorageClass),
			}
			objects = append(objects, object)
		}
	}

	isTruncated := tea.BoolValue(body.IsTruncated)
	return objects, isTruncated, nil
}

// GetObjectURL gets the download URL for an object.
func (p *OSSProvider) GetObjectURL(bucket, objectKey string) string {
	return fmt.Sprintf("https://%s.%s.aliyuncs.com/%s", bucket, p.region, objectKey)
}
