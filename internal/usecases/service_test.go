package usecases_test

import (
	"context"
	"reflect"
	"testing"

	"currency/internal/entities"
	"currency/internal/usecases"
	mock "currency/internal/usecases/mocks"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestService_GetLastPrice(t *testing.T) {
	type fields struct {
		storage *mock.MockStorage
		client  *mock.MockClient
	}
	type args struct {
		ctx    context.Context
		titles []string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr bool
		want    []entities.Coin
	}{
		{
			name: "GetLastPrice() failed - incorrect params",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetLastPrice() failed - The database does not have XRC coin",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetLastPrice() success",
			args: args{
				ctx:    context.Background(),
				titles: []string{"BTC, USDT"},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles).Return([]entities.Coin{{Title: "BTC"}, {Title: "USDT"}}, nil)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "USDT"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				storage: mock.NewMockStorage(ctrl),
				client:  nil,
			}

			tt.prepare(&f, tt.args)

			s, _ := usecases.NewService(f.storage, f.client)

			got, err := s.GetLastPrice(tt.args.ctx, tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetMaxPrice(t *testing.T) {
	type fields struct {
		storage *mock.MockStorage
		client  *mock.MockClient
	}
	type args struct {
		ctx    context.Context
		titles []string
		option usecases.Option
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr bool
		want    []entities.Coin
	}{
		{
			name: "GetMaxPrice() failed - incorrect params",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetMaxPrice() failed - The database does not have XRC coin",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetMaxPrice() success",
			args: args{
				ctx:    context.Background(),
				titles: []string{"BTC, USDT"},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return([]entities.Coin{{Title: "BTC"}, {Title: "USDT"}}, nil)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "USDT"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				storage: mock.NewMockStorage(ctrl),
				client:  nil,
			}

			tt.prepare(&f, tt.args)

			s, _ := usecases.NewService(f.storage, f.client)

			got, err := s.GetMaxPrice(tt.args.ctx, tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaxPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMaxPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetMinPrice(t *testing.T) {
	type fields struct {
		storage *mock.MockStorage
		client  *mock.MockClient
	}
	type args struct {
		ctx    context.Context
		titles []string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr bool
		want    []entities.Coin
	}{
		{
			name: "GetMinPrice() failed - incorrect params",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetMinPrice() failed - The database does not have XRC coin",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetMinPrice() success",
			args: args{
				ctx:    context.Background(),
				titles: []string{"BTC, USDT"},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return([]entities.Coin{{Title: "BTC"}, {Title: "USDT"}}, nil)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "USDT"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				storage: mock.NewMockStorage(ctrl),
				client:  nil,
			}

			tt.prepare(&f, tt.args)

			s, _ := usecases.NewService(f.storage, f.client)

			got, err := s.GetMinPrice(tt.args.ctx, tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMinPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMinPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAvgPrice(t *testing.T) {
	type fields struct {
		storage *mock.MockStorage
		client  *mock.MockClient
	}
	type args struct {
		ctx    context.Context
		titles []string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr bool
		want    []entities.Coin
	}{
		{
			name: "GetAvgPrice() failed - incorrect params",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetAvgPrice() failed - The database does not have XRC coin",
			args: args{
				ctx:    context.Background(),
				titles: []string{},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return(nil, entities.ErrInvalidParams)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetAvgPrice() success",
			args: args{
				ctx:    context.Background(),
				titles: []string{"BTC, USDT"},
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().Get(args.ctx, args.titles, gomock.Any()).Return([]entities.Coin{{Title: "BTC"}, {Title: "USDT"}}, nil)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "USDT"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				storage: mock.NewMockStorage(ctrl),
				client:  nil,
			}

			tt.prepare(&f, tt.args)

			s, _ := usecases.NewService(f.storage, f.client)

			got, err := s.GetAvgPrice(tt.args.ctx, tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvgPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAvgPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetCoinsFromAPI(t *testing.T) {
	type fields struct {
		storage *mock.MockStorage
		client  *mock.MockClient
	}
	type args struct {
		ctx    context.Context
		titles []string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr bool
		want    []entities.Coin
	}{
		{
			name: "GetCoinsFromAPI() failed - titles is nil and s.storage.GetTitles() failed",
			args: args{
				ctx:    context.Background(),
				titles: nil,
			},
			prepare: func(f *fields, args args) {
				f.storage.EXPECT().GetTitles(args.ctx).Return(nil, errors.New("s.storage.GetTitles() failed"))
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetCoinsFromAPI() failed - the client couldn't take the coin",
			args: args{
				ctx:    context.Background(),
				titles: nil,
			},
			prepare: func(f *fields, args args) {
				titles := []string{"BTC", "ETH"}
				gomock.InOrder(
					f.storage.EXPECT().GetTitles(args.ctx).Return(titles, nil),
					f.client.EXPECT().GetCoins(args.ctx, titles).Return(nil, errors.New("s.client.GetCoins() failed")),
				)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetCoinsFromAPI() failed - failed to save coins",
			args: args{
				ctx:    context.Background(),
				titles: nil,
			},
			prepare: func(f *fields, args args) {
				titles := []string{"BTC", "ETH"}
				coins := []entities.Coin{{Title: "BTC"}, {Title: "ETH"}}
				gomock.InOrder(
					f.storage.EXPECT().GetTitles(args.ctx).Return(titles, nil),
					f.client.EXPECT().GetCoins(args.ctx, titles).Return(coins, nil),
					f.storage.EXPECT().Store(args.ctx, coins).Return(errors.New("s.store failed")),
				)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "GetCoinsFromAPI() success",
			args: args{
				ctx:    context.Background(),
				titles: nil,
			},
			prepare: func(f *fields, args args) {
				titles := []string{"BTC", "ETH"}
				coins := []entities.Coin{{Title: "BTC"}, {Title: "ETH"}}
				gomock.InOrder(
					f.storage.EXPECT().GetTitles(args.ctx).Return(titles, nil),
					f.client.EXPECT().GetCoins(args.ctx, titles).Return(coins, nil),
					f.storage.EXPECT().Store(args.ctx, coins).Return(nil),
				)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "ETH"}},
		},
		{
			name: "GetCoinsFromAPI() success with titles parameter",
			args: args{
				ctx:    context.Background(),
				titles: []string{"BTC", "ETH"},
			},
			prepare: func(f *fields, args args) {
				coins := []entities.Coin{{Title: "BTC"}, {Title: "ETH"}}
				gomock.InOrder(
					f.client.EXPECT().GetCoins(args.ctx, args.titles).Return(coins, nil),
					f.storage.EXPECT().Store(args.ctx, coins).Return(nil),
				)
			},
			wantErr: false,
			want:    []entities.Coin{{Title: "BTC"}, {Title: "ETH"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				storage: mock.NewMockStorage(ctrl),
				client:  mock.NewMockClient(ctrl),
			}

			tt.prepare(&f, tt.args)

			s, _ := usecases.NewService(f.storage, f.client)

			got, err := s.GetCoinsFromAPI(tt.args.ctx, tt.args.titles...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCoinsFromAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCoinsFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
