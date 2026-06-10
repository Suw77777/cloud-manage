package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
	"strings"
)

// OSSProviderFactory creates an OSSProvider given credentials and region.
type OSSProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.OSSProvider, error)

// OSSService handles OSS business logic.
type OSSService struct {
	providerFactory OSSProviderFactory
}

// NewOSSService creates a new OSSService with default provider factory (cached).
func NewOSSService() *OSSService {
	return &OSSService{
		providerFactory: CachedFactory("oss", func(accessKeyId, accessKeySecret, region string) (provider.OSSProvider, error) {
			return aliyun.NewOSSProvider(accessKeyId, accessKeySecret, region)
		}),
	}
}

// NewOSSServiceWithProvider creates a new OSSService with custom provider factory (for testing).
func NewOSSServiceWithProvider(factory OSSProviderFactory) *OSSService {
	return &OSSService{providerFactory: factory}
}

// BucketAdapter is a provider-agnostic representation of an OSS bucket.
type BucketAdapter struct {
	Name             string `json:"name"`
	Location         string `json:"location"`
	CreationDate     string `json:"creationDate"`
	StorageClass     string `json:"storageClass"`
	ExtranetEndpoint string `json:"extranetEndpoint"`
	IntranetEndpoint string `json:"intranetEndpoint"`
}

// ObjectAdapter is a provider-agnostic representation of an OSS object.
type ObjectAdapter struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	ETag         string `json:"etag"`
	Type         string `json:"type"`
	StorageClass string `json:"storageClass"`
}

// ListBucketsResult holds the result of listing buckets.
type ListBucketsResult struct {
	Buckets []BucketAdapter `json:"buckets"`
}

// ListObjectsResult holds the result of listing objects.
type ListObjectsResult struct {
	Objects     []ObjectAdapter `json:"objects"`
	IsTruncated bool            `json:"isTruncated"`
	NextMarker  string          `json:"nextMarker,omitempty"`
}

// ListBuckets lists all OSS buckets.
func (s *OSSService) ListBuckets(accessKeyId, accessKeySecret, region string) (*ListBucketsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret and region are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OSS provider: %s", security.SanitizeErrorMessage(err))
	}

	buckets, err := p.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]BucketAdapter, 0, len(buckets))
	for _, b := range buckets {
		adapters = append(adapters, BucketAdapter{
			Name:             b.Name,
			Location:         b.Location,
			CreationDate:     b.CreationDate,
			StorageClass:     b.StorageClass,
			ExtranetEndpoint: b.ExtranetEndpoint,
			IntranetEndpoint: b.IntranetEndpoint,
		})
	}

	return &ListBucketsResult{Buckets: adapters}, nil
}

// ListObjects lists objects in an OSS bucket.
func (s *OSSService) ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix string, maxKeys int32) (*ListObjectsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || bucket == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret and bucket are required")
	}

	// Auto-detect bucket region if not specified or to handle cross-region access
	actualRegion := region
	if region == "" || region == "cn-hangzhou" {
		detectedRegion, err := s.DetectBucketRegion(accessKeyId, accessKeySecret, bucket)
		if err == nil && detectedRegion != "" {
			actualRegion = detectedRegion
		}
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, actualRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OSS provider: %s", security.SanitizeErrorMessage(err))
	}

	objects, isTruncated, err := p.ListObjects(bucket, prefix, maxKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]ObjectAdapter, 0, len(objects))
	for _, obj := range objects {
		adapters = append(adapters, ObjectAdapter{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
			Type:         obj.Type,
			StorageClass: obj.StorageClass,
		})
	}

	return &ListObjectsResult{
		Objects:     adapters,
		IsTruncated: isTruncated,
	}, nil
}

// DetectBucketRegion finds the actual region of a bucket by listing all buckets.
func (s *OSSService) DetectBucketRegion(accessKeyId, accessKeySecret, bucket string) (string, error) {
	p, err := s.providerFactory(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		return "", err
	}

	buckets, err := p.ListBuckets()
	if err != nil {
		return "", err
	}

	for _, b := range buckets {
		if b.Name == bucket {
			location := b.Location
			if strings.HasPrefix(location, "oss-") {
				return location[4:], nil
			}
			return location, nil
		}
	}

	return "", fmt.Errorf("bucket %s not found", bucket)
}
