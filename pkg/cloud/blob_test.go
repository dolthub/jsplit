package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCloudURITrue(t *testing.T) {
	var uri = "gs://bucket/path/to/file.json"
	var expected = true
	var actual = IsCloudURI(uri)

	require.Equal(t, expected, actual)
}

func TestIsCloudURIFalse(t *testing.T) {
	var uri = "/path/to/file.json"
	var expected = false
	var actual = IsCloudURI(uri)

	require.Equal(t, expected, actual)
}

func TestParseBlobUriCloud(t *testing.T) {
	var uri = "gs://bucket/path/to/file.json"
	var exoectedScheme = "gs"
	var expectedBucket = "bucket"
	var expectedKey = "path/to/file.json"
	var actualScheme, actualBucket, actualKey, err = ParseBlobURI(uri)

	require.Nil(t, err)
	require.Equal(t, exoectedScheme, actualScheme)
	require.Equal(t, expectedBucket, actualBucket)
	require.Equal(t, expectedKey, actualKey)
}

func TestParseBlobUriFS(t *testing.T) {
	var uri = "/path/to/file.json"
	var exoectedScheme = ""
	var expectedBucket = ""
	var expectedKey = "path/to/file.json"
	var actualScheme, actualBucket, actualKey, err = ParseBlobURI(uri)

	require.Nil(t, err)
	require.Equal(t, exoectedScheme, actualScheme)
	require.Equal(t, expectedBucket, actualBucket)
	require.Equal(t, expectedKey, actualKey)
}
