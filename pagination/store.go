package pagination

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	suffixFront   = ":front"
	suffixBack    = ":back"
	defaultExpire = time.Minute * 5
)

var copierOpt = &copier.Option{Converters: []copier.TypeConverter{{
	SrcType: primitive.ObjectID{},
	DstType: copier.String,
	Fn: func(src interface{}) (interface{}, error) {
		return src.(primitive.ObjectID).Hex(), nil
	},
}}}

type Store interface {
	GetCursor() any
	LoadCursor(ctx context.Context, lastToken string, backward bool) error
	StoreCursor(ctx context.Context, lastToken *string, first, last any) (*string, error)
}

type CacheStore struct {
	cursor     any
	cursorType reflect.Type
	cache      cache.Cache
	prefix     string
}

func NewCacheStore(c cache.Cache, cursor any, prefix string) *CacheStore {
	t := reflect.TypeOf(cursor)
	for t.Kind() == reflect.Interface || t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return &CacheStore{
		cursor:     cursor,
		cursorType: t,
		cache:      c,
		prefix:     prefix,
	}
}
func (s *CacheStore) GetCursor() any {
	return s.cursor
}

func (s *CacheStore) LoadCursor(ctx context.Context, lastToken string, backward bool) error {
	var key string
	if backward {
		key = s.prefix + lastToken + suffixFront
	} else {
		key = s.prefix + lastToken + suffixBack
	}
	s.cursor = reflect.New(s.cursorType).Interface()
	err := s.cache.GetCtx(ctx, key, s.cursor)
	if err != nil {
		return err
	}
	return nil
}

func (s *CacheStore) StoreCursor(ctx context.Context, lastToken *string, first, last any) (*string, error) {
	if lastToken == nil {
		lastToken = new(string)
		*lastToken = uuid.New().String()
	}
	front := reflect.New(s.cursorType).Interface()
	err := copier.CopyWithOption(front, first, *copierOpt)
	if err != nil {
		return nil, err
	}
	//TODO 假如第一次成功，第二次失败会发生什么
	err = s.cache.SetWithExpireCtx(ctx, s.prefix+*lastToken+suffixFront, front, defaultExpire)
	if err != nil {
		return nil, err
	}

	back := reflect.New(s.cursorType).Interface()
	err = copier.CopyWithOption(back, last, *copierOpt)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetWithExpireCtx(ctx, s.prefix+*lastToken+suffixBack, back, defaultExpire)
	if err != nil {
		return nil, err
	}
	return lastToken, nil
}

type RawStore struct {
	cursor     any
	cursorType reflect.Type
}

func NewRawStore(cursor any) *RawStore {
	t := reflect.TypeOf(cursor)
	for t.Kind() == reflect.Interface || t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return &RawStore{
		cursor:     cursor,
		cursorType: t,
	}
}

func (s *RawStore) GetCursor() any {
	return s.cursor
}

func (s *RawStore) LoadCursor(_ context.Context, lastToken string, backward bool) error {
	cursors := reflect.New(reflect.ArrayOf(2, reflect.PointerTo(s.cursorType)))
	err := json.Unmarshal([]byte(lastToken), cursors.Interface())
	if backward {
		s.cursor = cursors.Elem().Index(0).Interface()
	} else {
		s.cursor = cursors.Elem().Index(1).Interface()
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *RawStore) StoreCursor(_ context.Context, lastToken *string, first, last any) (*string, error) {
	front := reflect.New(s.cursorType).Interface()
	err := copier.CopyWithOption(front, first, *copierOpt)
	if err != nil {
		return nil, err
	}
	back := reflect.New(s.cursorType).Interface()
	err = copier.CopyWithOption(back, last, *copierOpt)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal([2]any{front, back})
	if err != nil {
		return nil, err
	}
	if lastToken == nil {
		lastToken = new(string)
	}
	*lastToken = string(bytes)
	return lastToken, nil
}
