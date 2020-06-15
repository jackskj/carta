package carta_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jackskj/carta"
	td "github.com/jackskj/carta/testdata"
	"github.com/jackskj/carta/testdata/initdb"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	// "github.com/golang/protobuf/proto"
)

const (
	bufSize = 1024 * 1024
	pg      = "postgres"
	mysql   = "mysql"
)

var (
	update = false
	initDB = true
)

var (
	conn        *grpc.ClientConn
	ctx         context.Context
	dbs         map[string]*sql.DB
	grpcServer  *grpc.Server
	lis         *bufconn.Listener
	requests    *td.Requests
	testResults map[string]interface{}
)

// Generate test data before running tests
// Start local server with bufconn
func setup() {
	testResults = make(map[string]interface{})
	ctx = context.Background()
	lis = bufconn.Listen(bufSize)
	grpcServer = grpc.NewServer()
	dbs = map[string]*sql.DB{
		pg:    td.GetPG(),
		mysql: td.GetMySql(),
	}
	initdb.RegisterInitServiceServer(grpcServer, &initdb.InitServiceMapServer{DBs: dbs})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	if connection, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure()); err != nil {
		log.Fatalf("bufnet dial fail: %v", err)
	} else {
		conn = connection
	}
	if initDB {
		createDatabase(dbs)
	}
}

func TestMain(m *testing.M) {
	updatePtr := flag.Bool("update", false, "update the golden file, results are always considered correct")
	initdbPtr := flag.Bool("initdb", false, "initialize and populate testing database")
	flag.Parse()
	update = *updatePtr
	initDB = *initdbPtr
	setup()
	code := m.Run()
	goldenFile := "testdata/mapper.golden"
	if update {
		// update golden file
		updateGoldenFile(goldenFile)
	} else {
		// compare existing results
		compareResults(goldenFile)
	}
	teardown()
	os.Exit(code)
}

func createDatabase(dbs map[string]*sql.DB) {
	requests = td.GenerateRequests()
	for dbName, _ := range dbs {
		meta := &initdb.Meta{Db: dbName}
		initService := initdb.NewInitServiceClient(conn)
		initService.InitDB(ctx, &initdb.InitRequest{Meta: &initdb.Meta{Db: dbName}})
		for i := 0; i < len(requests.InsertAuthorRequests); i++ {
			requests.InsertAuthorRequests[i].Meta = meta
			if _, err := initService.InsertAuthor(ctx, requests.InsertAuthorRequests[i]); err != nil {
				log.Fatalf("InsertAuthor: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertBlogRequests); i++ {
			requests.InsertBlogRequests[i].Meta = meta
			if _, err := initService.InsertBlog(ctx, requests.InsertBlogRequests[i]); err != nil {
				log.Fatalf("InsertBlog: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertCommentRequests); i++ {
			requests.InsertCommentRequests[i].Meta = meta
			if _, err := initService.InsertComment(ctx, requests.InsertCommentRequests[i]); err != nil {
				log.Fatalf("InsertComment: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertPostRequests); i++ {
			requests.InsertPostRequests[i].Meta = meta
			if _, err := initService.InsertPost(ctx, requests.InsertPostRequests[i]); err != nil {
				log.Fatalf("InsertPost: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertPostTagRequests); i++ {
			requests.InsertPostTagRequests[i].Meta = meta
			if _, err := initService.InsertPostTag(ctx, requests.InsertPostTagRequests[i]); err != nil {
				log.Fatalf("InsertPostTag: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertTagRequests); i++ {
			requests.InsertTagRequests[i].Meta = meta
			if _, err := initService.InsertTag(ctx, requests.InsertTagRequests[i]); err != nil {
				log.Fatalf("InsertTag: %s", err)
			}
		}
	}
}

func updateGoldenFile(goldenFile string) {
	jsonResult := generateResultBytes()
	if err := ioutil.WriteFile(goldenFile, jsonResult, 0644); err != nil {
		log.Fatalln(err)
	}
}

func compareResults(goldenFile string) {
	goldenFileJson, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		log.Fatalln(err)
	}

	jsonResult := generateResultBytes()

	resultDiff := diff.New()
	d, err := resultDiff.Compare(goldenFileJson, jsonResult)
	if err != nil {
		log.Fatalln(err)
	}
	formatter := formatter.NewDeltaFormatter()
	diffString, err := formatter.Format(d)
	if diffString != "{}\n" {
		log.Println("Results Do Not Match Golden File, " +
			"if this is expecred result with go test with --update")
		log.Fatalln(diffString)
	}
}

func generateResultBytes() []byte {
	var jsonResult []byte
	if r, err := json.MarshalIndent(testResults, "", "    "); err != nil {
		log.Fatalln(err)
	} else {
		jsonResult = r
	}
	return jsonResult
}

func teardown() {
	defer conn.Close()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func query(rawSql string) map[string]*sql.Rows {
	resp := map[string]*sql.Rows{}
	for dbName, db := range dbs {
		stmt, err := db.Prepare(rawSql)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if rows, err := stmt.Query(); err != nil {
			log.Fatal(err)
		} else {
			resp[dbName] = rows
		}
	}
	return resp
}

func queryPG(rawSql string) (rows *sql.Rows) {
	var err error
	if rows, err = dbs[pg].Query(rawSql); err != nil {
		log.Fatal(err)
	}
	return
}

func queryMysql(rawSql string) (rows *sql.Rows) {
	var err error
	stmt, err := dbs[mysql].Prepare(rawSql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if rows, err = stmt.Query(); err != nil {
		log.Fatal(err)
	}
	return
}

func TestBlog(m *testing.T) {
	ans := []byte{}
	for _, rows := range query(td.BlogQuery) {
		resp := []td.Blog{}
		if err := carta.Map(rows, &resp); err != nil {
			log.Fatal(err.Error())
		}
		e, _ := json.Marshal(resp)
		if len(ans) == 0 {
			ans = e
		} else if string(ans) != string(e) {
			log.Fatal(errors.New("Test Blog Produced Inconsistent Results"))
		}
		testResults["TestBlog"] = resp
	}
}

func TestNull(m *testing.T) {
	respPG := []td.NullTest{}
	if err := carta.Map(queryPG(td.NullQueryPG), &respPG); err != nil {
		log.Fatal(err.Error())
	}
	respMySQL := []td.NullTest{}
	if err := carta.Map(queryMysql(td.NullQueryMySql), &respMySQL); err != nil {
		log.Fatal(err.Error())
	}
	ansPG, _ := json.Marshal(respPG)
	ansMySQL, _ := json.Marshal(respMySQL)
	if string(ansPG) != string(ansMySQL) {
		log.Fatal(errors.New("Test Null Produced Inconsistent Results"))
	}
	testResults["TestNull"] = respPG
}

func TestNotNull(m *testing.T) {
	respPG := []td.NullTest{}
	if err := carta.Map(queryPG(td.NotNullQueryPG), &respPG); err != nil {
		log.Fatal(err.Error())
	}
	respMySQL := []td.NullTest{}
	if err := carta.Map(queryMysql(td.NotNullQueryMySQL), &respMySQL); err != nil {
		log.Fatal(err.Error())
	}
	ansPG, _ := json.Marshal(respPG)
	ansMySQL, _ := json.Marshal(respMySQL)
	if string(ansPG) != string(ansMySQL) {
		log.Println(string(ansMySQL))
		log.Fatal(errors.New("Test Not Null Produced Inconsistent Results"))
	}
	testResults["TestNotNull"] = respPG
}

func TestPGTypes(m *testing.T) {
	resp := []td.PGDTypes{}
	if err := carta.Map(queryPG(td.PGDTypesQuery), &resp); err != nil {
		log.Fatal(err.Error())
	}
	testResults["TestPGTypes"] = resp
}

func TestRelation(m *testing.T) {
	ans := []byte{}
	for _, rows := range query(td.RelationTestQuery) {
		resp := []td.RelationTest{}
		if err := carta.Map(rows, &resp); err != nil {
			log.Fatal(err.Error())
		}
		e, _ := json.Marshal(resp)
		if len(ans) == 0 {
			ans = e
		} else if string(ans) != string(e) {
			log.Fatal(errors.New("Test Blog Produced Inconsistent Results"))
		}
		testResults["TestRelation"] = resp
	}
}
