
# Carta
Dead simple SQL data mapper for complex Go structs. 

Loads SQL data onto Go structs while keeping track of has-one, has-many relationships


## Examples 

Assume you have a schema containing has-one and has-many relationships:



![schema](https://i.ibb.co/SPH3zhQ/Schema.png)

And we wish to map the the following query on a Go struct:
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
Carta maps retrieved row while keeping track of those relationships. 

All you need to do is: 
```
blogs := []Blog
carta.Map(rows, &blogs)
```
Resulting Mapping: 
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
Sqlx does not track has-many relationships when marshaling SQL data onto a struct. This works fine when all your relationships are at most has-one (Blog has one Author) ie, each SQL row corresponds to one struct. Has-many relationships (Blog has man Posts), are not mapped correctly since a new instance of a struct is always instantiated even if data representing the same instance has already been mapped.
  
## Installation 
```
go get -u github.com/jackskj/carta
```


## Important Notes 

 - Carta automatically removes any duplicate rows returned by your query. If this is not a desired outcome, you should include a uniquely identifiable columns in your query and the corresponding fields in your structs.

## Approach
Carta adopts the "database mapping" approach (described in Martin Fowler's [book](https://books.google.com/books?id=FyWZt5DdvFkC&lpg=PA1&dq=Patterns%20of%20Enterprise%20Application%20Architecture%20by%20Martin%20Fowler&pg=PT187#v=onepage&q=active%20record&f=false)) which is useful among organizations with strict code review processes.

Carta is not an object-relational mapper(ORM). With large and complex datasets, using ORMs becomes restrictive and reduces performance when working with complex queries. 

### License
Apache License
