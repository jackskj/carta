# Carta
Dead simple SQL data mapper for complex Go structs. 

Load SQL data onto Go structs while keeping track of has-one, has-many relationships

## Examples 
Using carta is very simple. All you need to do is: 
```
// 1) Run your query
if rows, err = sqlDB.Query(blogQuery); err != nil {
	// err
}

// 2) Instantiate a slice(or struct) which you want to populate 
blogs := []Blog

// 3) Map the SQL rows to your slice
carta.Map(rows, &blogs)
```

Assume that in above exmple, we are using a schema containing has-one and has-many relationships:

![schema](https://i.ibb.co/SPH3zhQ/Schema.png)

And here is our SQL query along with the corresponging Go struct:
```
select
       id          as  blog_id,
       title       as  blog_title,
       P.id        as  posts_id,         
       P.name      as  posts_name,
       A.id        as  author_id,      
       A.username  as  author_username
from blog
       left outer join author A    on  blog.author_id = A.id
       left outer join post P      on  blog.id = P.blog_id
```

```
type Blog struct {
        Id     int    `db:"blog_id"`
        Title  string `db:"blog_title"`
        Posts  []Post
        Author Author
}
type Post struct {
        Id   int    `db:"posts_id"`
        Name string `db:"posts_name"`
}
type Author struct {
        Id       int    `db:"author_id"`
        Username string `db:"author_username"`
}
```
Carta will map the SQL rows while keeping track of those relationships. 

Results: 
```
rows:
blog_id | blog_title | posts_id | posts_name | author_id | author_username
1       | Foo        | 1        | Bar        | 1         | John
1       | Foo        | 2        | Baz        | 1         | John
2       | Egg        | 3        | Beacon     | 2         | Ed

blogs:
[{
	"blog_id": 1,
	"blog_title": "Foo",
	"author": {
		"author_id": 1,
		"author_username": "John"
	},
	"posts": [{
			"post_id": 1,
			"posts_name": "Bar"
		}, {
			"post_id": 2,
			"posts_name": "Baz"
		}]
}, {
	"blog_id": 2,
	"blog_title": "Egg",
	"author": {
		"author_id": 2,
		"author_username": "Ed"
	},
	"posts": [{
			"post_id": 3,
			"posts_name": "Beacon"
		}]
}]
```


## Comparison to Related Projects
#### GORM
Carta is NOT an an object-relational mapper(ORM). Read more in [Approach](#Approach)

#### sqlx
Sqlx does not track has-many relationships when mapping SQL data. This works fine when all your relationships are at most has-one (Blog has one Author) ie, each SQL row corresponds to one struct. However, handling has-many relationships (Blog has many Posts), requires  running many queries or running manual post-processing of the result. Carta handles these complexities automatically.

## Guide

### Column and Field Names

Carta will match your SQL columns with corresponding fields. You can use a "db" tag to represent a specific column name.
Example:

```
type Blog struct {
	// When tag is not used, the snake case of the fiels is used
	BlogId int // expected column name : "blog_id"

	// When tag is specified, it takes priority
	Abc string `db:"blog_title"` // expected column name : "blog_title"

	// If you define multiple fiels with the same struct,
	// you can use a tag to identify a column prefix (with underscore concatination)

	Writer Author `db: "writer"` // possible column names: "author_id", "writer_author_id",

	Reviewer Author `db: "reviewer"` // possible column names: "author_id", "rewiewer_author_id",
}

type Author struct {
	AuthorId int `db:"author_id"`
}
```
### Data Types and Relationships

Any primative types, time.Time, protobuf Timestamp, and sql.NullXXX can be loaded with Carta.
These types are one-to-one mapped with your SQL columns

To define more complex SQL relationships use slices and structs as in example below:

```
type Blog struct {
	BlogId int  // Will map directly with "blog_id" column 

	// If your SQL data can potentially be "null", make sure to use pointers or sql.NullXXX
	AuthorId  *int
	CreatedOn *timestamp.Timestamp // protobuf timestamp
	UpdatedOn *time.Time
	SonsorId  sql.NullInt64

	// To define has-one relationship, use nested struct or pointer to a struct
	Author *Author

	// To define has-many relationship, use slices
	Posts []*Post // valid options include: *[]*Post, []*Post , *[]Post, and []Post

	// when your has-many relationship only has one corresponding column, you can use a slice of a settable type
	TagIds     []int
	CommentIds []sql.NullInt64
}
```

### Drivers 

Recommended driver for Postgres is [lib/pg](https://github.com/lib/pq), for MySql use [go-sql-driver/mysql](https://github.com/go-sql-driver/mysqlq).

When using MySql, carta expects TIME and DATETIME to arrive in time.Time format. Therefore, make sure to add "parseTime=true" in your connection string. 

Time data conversion from plain text will be added in the future versions of Carta.
## Installation 
```
go get -u github.com/jackskj/carta
```


## Important Notes 

Carta removes any duplicate rows. This is a side effect of the data mapping as it is unclear which object to instantiate if the same data arrives more than once.
If this is not a desired outcome, you should include a uniquely identifiable columns in your query and the corresponding fields in your structs.
 
To prevent relatively expensive reflect operations, carta caches the structure of your struct using the column mames of your query response as well as the type of your struct. 



## Approach
Carta adopts the "database mapping" approach (described in Martin Fowler's [book](https://books.google.com/books?id=FyWZt5DdvFkC&lpg=PA1&dq=Patterns%20of%20Enterprise%20Application%20Architecture%20by%20Martin%20Fowler&pg=PT187#v=onepage&q=active%20record&f=false)) which is useful among organizations with strict code review processes.

Carta is not an object-relational mapper(ORM). With large and complex datasets, using ORMs becomes restrictive and reduces performance when working with complex queries. 

### License
Apache License
