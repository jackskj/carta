package initdb

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"text/template"

	"github.com/Masterminds/sprig"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	sqlTpl "github.com/jackskj/carta/testdata/initdb/sql"
	tpl "github.com/jackskj/protoc-gen-map/templates"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InitServiceMapServer struct {
	DBs map[string]*sql.DB
}

var InitTemplate, _ = template.New("InitTemplate").Funcs(sprig.TxtFuncMap()).Funcs(tpl.Funcs()).Parse(sqlTpl.Init)

func (m *InitServiceMapServer) InitDB(ctx context.Context, r *InitRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InitDB", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil
}

func (m *InitServiceMapServer) InsertAuthor(ctx context.Context, r *InsertAuthorRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertAuthor", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil
}

func (m *InitServiceMapServer) InsertBlog(ctx context.Context, r *InsertBlogRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertBlog", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil
}

func (m *InitServiceMapServer) InsertComment(ctx context.Context, r *InsertCommentRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertComment", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil
}

func (m *InitServiceMapServer) InsertPost(ctx context.Context, r *InsertPostRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertPost", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil

}

func (m *InitServiceMapServer) InsertPostTag(ctx context.Context, r *InsertPostTagRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertPostTag", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		log.Fatal(err.Error())
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil

}

func (m *InitServiceMapServer) InsertTag(ctx context.Context, r *InsertTagRequest) (*EmptyResponse, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := InitTemplate.ExecuteTemplate(sqlBuffer, "InsertTag", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	if _, err := m.DBs[r.GetMeta().Db].Exec(rawSql); err != nil {
		return nil, status.Error(codes.InvalidArgument, "error: executing query")
	}
	return &EmptyResponse{}, nil
}
