package testdata

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackskj/carta/testdata/initdb"

	// default deivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	chars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	iLength = 10
)

var (
	seed        = rand.New(rand.NewSource(1))
	tsSample, _ = ptypes.TimestampProto(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	blogs       = rand.Perm(iLength)
	authors     = rand.Perm(iLength)
	posts       = rand.Perm(40)
	comments    = rand.Perm(200)
	tags        = rand.Perm(iLength)
	blogAuthor  map[int]int
	blogPost    map[int][]int
	postAuthor  map[int]int
	postBlog    map[int]int
	postComment map[int][]int
	commentPost map[int]int
	postTag     [][2]int
)

type Requests struct {
	InsertAuthorRequests  []*initdb.InsertAuthorRequest
	InsertBlogRequests    []*initdb.InsertBlogRequest
	InsertCommentRequests []*initdb.InsertCommentRequest
	InsertPostRequests    []*initdb.InsertPostRequest
	InsertPostTagRequests []*initdb.InsertPostTagRequest
	InsertTagRequests     []*initdb.InsertTagRequest
}

func lorem() string {
	length := 10
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = chars[seed.Intn(len(chars))]
	}
	return string(bytes)
}

func GenerateRequests(dbname string) *Requests {
	meta := &initdb.Meta{Db: dbname}
	blogAuthor = make(map[int]int)
	for i, b := range blogs {
		blogAuthor[i] = b
	}
	blogPost = make(map[int][]int)
	postAuthor = make(map[int]int)
	postBlog = make(map[int]int)
	for _, post := range posts {
		randN := rand.Intn(iLength)
		blogPost[randN] = append(blogPost[randN], post)
		postBlog[post] = randN
		postAuthor[post] = blogAuthor[randN]
	}
	postComment = make(map[int][]int)
	commentPost = make(map[int]int)
	for _, comment := range comments {
		randN := rand.Intn(len(posts))
		postComment[randN] = append(postComment[randN], comment)
		commentPost[comment] = randN
	}
	for i := 0; i < len(posts); i++ {
		t := rand.Perm(len(tags))
		for j := 0; j < 3; j++ {
			postTag = append(postTag, [2]int{i, t[j]})
		}
	}

	requests := Requests{
		InsertAuthorRequests:  []*initdb.InsertAuthorRequest{},
		InsertBlogRequests:    []*initdb.InsertBlogRequest{},
		InsertCommentRequests: []*initdb.InsertCommentRequest{},
		InsertPostRequests:    []*initdb.InsertPostRequest{},
		InsertPostTagRequests: []*initdb.InsertPostTagRequest{},
		InsertTagRequests:     []*initdb.InsertTagRequest{},
	}
	for _, i := range tags {
		requests.InsertTagRequests = append(
			requests.InsertTagRequests,
			&initdb.InsertTagRequest{
				Meta: meta,
				Id:   uint32(i),
				Name: lorem(),
			},
		)
	}
	for _, i := range authors {
		requests.InsertAuthorRequests = append(
			requests.InsertAuthorRequests,
			&initdb.InsertAuthorRequest{
				Meta:             meta,
				Id:               uint32(i),
				Username:         lorem(),
				Password:         lorem(),
				Email:            lorem(),
				Bio:              lorem(),
				FavouriteSection: getSection(),
			},
		)
	}
	for _, i := range postTag {
		requests.InsertPostTagRequests = append(
			requests.InsertPostTagRequests,
			&initdb.InsertPostTagRequest{
				Meta:   meta,
				PostId: uint32(i[0]),
				TagId:  uint32(i[1]),
			},
		)
	}
	for _, i := range comments {
		requests.InsertCommentRequests = append(
			requests.InsertCommentRequests,
			&initdb.InsertCommentRequest{
				Meta:    meta,
				Id:      uint32(i),
				PostId:  uint32(commentPost[i]),
				Name:    lorem(),
				Comment: lorem(),
			},
		)
	}

	for _, i := range posts {
		requests.InsertPostRequests = append(
			requests.InsertPostRequests,
			&initdb.InsertPostRequest{
				Meta:      meta,
				Id:        uint32(i),
				AuthorId:  uint32(postAuthor[i]),
				BlogId:    uint32(postBlog[i]),
				CreatedOn: tsSample,
				Section:   getSection(),
				Subject:   lorem(),
				Draft:     lorem(),
				Body:      lorem(),
			},
		)
	}

	for _, i := range blogs {
		requests.InsertBlogRequests = append(
			requests.InsertBlogRequests,
			&initdb.InsertBlogRequest{
				Meta:     meta,
				Id:       uint32(i),
				Title:    lorem(),
				AuthorId: uint32(blogAuthor[i]),
			},
		)
	}
	return &requests
}

func getSection() string {
	sections := []string{
		"cooking",
		"painting",
		"woodworking",
		"snowboarding",
	}
	return sections[rand.Intn(len(sections))]
}

func GetPG() *sql.DB {
	connStr := "postgres://postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil || db.Ping() != nil {
		log.Println("Cannot connect to testing database, \n" +
			"to run local postgres testing DB, run \"docker run --env POSTGRES_HOST_AUTH_METHOD=trust -d  -p 5432:5432 postgres\"")
		log.Fatal(err)
	}
	return db
}

func GetMySql() *sql.DB {
	connStr := "root@/mysql"
	db, err := sql.Open("mysql", connStr)
	if err != nil || db.Ping() != nil {
		log.Println("Cannot connect to testing database, \n" +
			"to run local postgres testing DB, run \"docker run --name carta-mysql-test -d  --env MYSQL_ALLOW_EMPTY_PASSWORD=yes --env MYSQL_DATABASE=mysql -p 3306:3306 mysql\"")
		log.Fatal(err)
	}
	return db
}

func GetSampleTS() timestamppb.Timestamp {
	return *tsSample
}
