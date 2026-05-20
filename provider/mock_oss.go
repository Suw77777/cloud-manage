package provider

// MockOSSProvider implements OSSProvider for testing.
type MockOSSProvider struct {
	Buckets     []OSSBucket
	Objects     []OSSObject
	IsTruncated bool
	Err         error
}

func (m *MockOSSProvider) ListBuckets() ([]OSSBucket, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Buckets, nil
}

func (m *MockOSSProvider) ListObjects(bucket, prefix string, maxKeys int32) ([]OSSObject, bool, error) {
	if m.Err != nil {
		return nil, false, m.Err
	}
	return m.Objects, m.IsTruncated, nil
}
