package docker

import (
	"context"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
)

type dockerRangeFetcher struct {
	*dockerBase
	start int
	end   int
}

func (r dockerRangeFetcher) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {

}
