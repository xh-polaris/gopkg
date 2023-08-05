package esp

import (
	"context"
	"github.com/xh-polaris/gopkg/pagination"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type (
	EsPaginator struct {
		store pagination.Store
		opts  *pagination.PaginationOptions
	}
)

func NewEsPaginator(store pagination.Store, opts *pagination.PaginationOptions) *EsPaginator {
	opts.EnsureSafe()
	return &EsPaginator{
		store: store,
		opts:  opts,
	}
}

// MakeSortOptions 生成ID分页查询选项
func (p *EsPaginator) MakeSortOptions(ctx context.Context) ([]types.SortCombinations, []types.FieldValue, error) {
	if p.opts.LastToken != nil {
		err := p.store.LoadCursor(ctx, *p.opts.LastToken, *p.opts.Backward)
		if err != nil {
			return nil, nil, err
		}
	}

	cursor := p.store.GetCursor()
	sort, sa, err := cursor.(EsCursor).MakeSortOptions(*p.opts.Backward)
	if err != nil {
		return nil, nil, err
	}
	return sort, sa, nil
}

func (p *EsPaginator) StoreCursor(ctx context.Context, first, last any) error {
	token, err := p.store.StoreCursor(ctx, p.opts.LastToken, first, last)
	p.opts.LastToken = token
	return err
}
