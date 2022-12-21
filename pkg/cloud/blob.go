package cloud

import (
	"context"
	"net/url"
	"strings"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob" // required by CDK as blob driver
	_ "gocloud.dev/blob/gcsblob"  // required by CDK as blob driver
	_ "gocloud.dev/blob/s3blob"   // required by CDK as blob driver
)

func OpenBucket(ctx context.Context, uri string) (*blob.Bucket, error) {
	var (
		err error
		b   *blob.Bucket
	)

	b, err = blob.OpenBucket(ctx, uri)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewWriter(ctx context.Context, uri string) (*blob.Writer, error) {
	var (
		err error
		bkt string
		k   string
	)

	_, bkt, k, err = ParseBlobURI(uri)
	if err != nil {
		return nil, err
	}

	b, err := OpenBucket(ctx, bkt)
	if err != nil {
		return nil, err
	}

	w, err := b.NewWriter(ctx, k, nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func NewReader(ctx context.Context, uri string) (*blob.Reader, error) {
	var (
		err error
		bkt string
		k   string
	)

	_, bkt, k, err = ParseBlobURI(uri)
	if err != nil {
		return nil, err
	}

	b, err := OpenBucket(ctx, bkt)
	if err != nil {
		return nil, err
	}

	r, err := b.NewReader(ctx, k, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ParseBlobURI(uri string) (scheme, bucket, key string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", "", err
	}

	return u.Scheme, u.Host, strings.TrimLeft(u.Path, "/"), nil
}

func IsCloudURI(uri string) bool {
	s, _, _, err := ParseBlobURI(uri)
	if err != nil {
		return false
	}

	return s != ""
}
