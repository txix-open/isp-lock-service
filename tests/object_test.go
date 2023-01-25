package tests_test

//
// import (
// 	"context"
// 	"testing"
//
// 	"isp-lock-service/assembly"
//
// 	"github.com/integration-system/isp-kit/dbx"
// 	"github.com/integration-system/isp-kit/grpc/client"
// 	"github.com/integration-system/isp-kit/test"
// 	"github.com/integration-system/isp-kit/test/dbt"
// 	"github.com/integration-system/isp-kit/test/grpct"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )
//
// type Object struct {
// 	Name string
// }
//
// func TestGetAll(t *testing.T) {
// 	t.Parallel()
// 	assert, testDb, cli := prepareTest(t)
//
// 	result := make([]Object, 0)
// 	err := cli.Invoke("isp-lock-service/object/all").
// 		ReadJsonResponse(&result).
// 		Do(context.Background())
// 	assert.NoError(err)
// 	assert.Empty(result)
//
// 	testDb.Must().Exec("insert into object (id, name) values ($1, $2)", 1, "a")
// 	testDb.Must().Exec("insert into object (id, name) values ($1, $2)", 2, "b")
//
// 	result = make([]Object, 0)
// 	err = cli.Invoke("isp-lock-service/object/all").
// 		ReadJsonResponse(&result).
// 		Do(context.Background())
// 	assert.NoError(err)
//
// 	expected := []Object{{
// 		Name: "a",
// 	}, {
// 		Name: "b",
// 	}}
// 	assert.EqualValues(expected, result)
// }
//
// func TestGetById(t *testing.T) {
// 	t.Parallel()
// 	assert, testDb, cli := prepareTest(t)
//
// 	testDb.Must().Exec("insert into object (id, name) values ($1, $2)", 1, "a")
//
// 	type reqBody struct {
// 		Id int
// 	}
//
// 	// empty req body
// 	result := Object{}
// 	err := cli.Invoke("isp-lock-service/object/get_by_id").
// 		ReadJsonResponse(&result).
// 		Do(context.Background())
// 	assert.Error(err)
// 	assert.EqualValues(codes.InvalidArgument, status.Code(err))
//
// 	// id is required
// 	result = Object{}
// 	err = cli.Invoke("isp-lock-service/object/get_by_id").
// 		ReadJsonResponse(&result).
// 		JsonRequestBody(reqBody{}).
// 		Do(context.Background())
// 	assert.Error(err)
// 	assert.EqualValues(codes.InvalidArgument, status.Code(err))
//
// 	// not found
// 	result = Object{}
// 	err = cli.Invoke("isp-lock-service/object/get_by_id").
// 		ReadJsonResponse(&result).
// 		JsonRequestBody(reqBody{Id: 2}).
// 		Do(context.Background())
// 	assert.Error(err)
// 	assert.EqualValues(codes.NotFound, status.Code(err))
//
// 	// happy path
// 	result = Object{}
// 	err = cli.Invoke("isp-lock-service/object/get_by_id").
// 		ReadJsonResponse(&result).
// 		JsonRequestBody(reqBody{Id: 1}).
// 		Do(context.Background())
// 	assert.NoError(err)
//
// 	expected := Object{Name: "a"}
// 	assert.EqualValues(expected, result)
// }
//
// func prepareTest(t *testing.T) (*require.Assertions, *dbt.TestDb, *client.Client) {
// 	t.Helper()
// 	test, assert := test.New(t)
// 	testDb := dbt.New(test, dbx.WithMigration("../migrations"))
//
// 	locator := assembly.NewLocator(test.Logger())
// 	_, cli := grpct.TestServer(test, locator.Handler())
//
// 	return assert, testDb, cli
// }