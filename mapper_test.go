package carta_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
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
	marsh  = jsonpb.Marshaler{}
	update = false
	initDB = true
)

var (
	conn        *grpc.ClientConn
	ctx         context.Context
	dbs         map[string]*sql.DB
	pgdb        *sql.DB
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
	pgdb = td.GetPG()
	dbs = map[string]*sql.DB{
		"postgres": pgdb,
		// "mysql":    td.GetMySql(),
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
	initdbPtr := flag.Bool("initdb", true, "initialize and populate testing database")
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
	for dbName, _ := range dbs {
		requests = td.GenerateRequests(dbName)
		initService := initdb.NewInitServiceClient(conn)
		initService.InitDB(ctx, &initdb.InitRequest{Meta: &initdb.Meta{Db: dbName}})
		for i := 0; i < len(requests.InsertAuthorRequests); i++ {
			if _, err := initService.InsertAuthor(ctx, requests.InsertAuthorRequests[i]); err != nil {
				log.Fatalf("InsertAuthor: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertBlogRequests); i++ {
			if _, err := initService.InsertBlog(ctx, requests.InsertBlogRequests[i]); err != nil {
				log.Fatalf("InsertBlog: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertCommentRequests); i++ {
			if _, err := initService.InsertComment(ctx, requests.InsertCommentRequests[i]); err != nil {
				log.Fatalf("InsertComment: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertPostRequests); i++ {
			if _, err := initService.InsertPost(ctx, requests.InsertPostRequests[i]); err != nil {
				log.Fatalf("InsertPost: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertPostTagRequests); i++ {
			if _, err := initService.InsertPostTag(ctx, requests.InsertPostTagRequests[i]); err != nil {
				log.Fatalf("InsertPostTag: %s", err)
			}
		}
		for i := 0; i < len(requests.InsertTagRequests); i++ {
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

func queryPG(rawSql string) (rows *sql.Rows) {
	var (
		err error
	)
	if rows, err = pgdb.Query(rawSql); err != nil {
		log.Fatal(err)
	}
	return
}

func TestBlog(m *testing.T) {
	resp := []td.Blog{}
	if err := carta.Map(queryPG(td.BlogQuery), &resp); err != nil {
		log.Fatal(err.Error())
	}
	testResults["TestBlog"] = resp
}

func TestNull(m *testing.T) {
	resp := []td.NullTest{}
	if err := carta.Map(queryPG(td.NullQuery), &resp); err != nil {
		log.Fatal(err.Error())
	}
	testResults["TestNull"] = resp
}

func TestNotNull(m *testing.T) {
	resp := []td.NullTest{}
	if err := carta.Map(queryPG(td.NotNullQuery), &resp); err != nil {
		log.Fatal(err.Error())
	}
	testResults["TestNotNull"] = resp
}
