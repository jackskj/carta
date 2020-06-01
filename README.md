# Carta
Dead simple SQL data mapper for complex Go structs. 

Loads SQL data onto Go structs while keeping track of complex has-one, has-many relationships

!work in progress! 

## Examples 

Assume you have a complex query:
```
select
       id          as  blog_id,
       title       as  blog_title,
       P.id        as  posts_id,
       P.name      as  posts_name,
       A.id        as  author_id,
       A.username  as  author_username
from blogs
       left outer join author A    on  blog.author_id = A.id
       left outer join post P      on  blog.id = P.blog_id
```
And we wish to map the response onto the following struct:
```
type Blog {
	Id int
	Title string
	Posts []Post
	Author Author
}
type Post {
	Id int
	Name string 
}
type Author {
	Id int
	Username string
}
```
Note that this query/struct pair involves a number of has-one and has-many relationships. Carta maps retrieved row while keeping track of those relationships. 

All you need to do is: 
```
blogs := []Blog
carta.Map(rows, &blogs)
```


## Installation 
```
go get -u github.com/jackskj/carta
```


## Important Notes 
Carta automatically removes any duplicate rows returned by your query. If this is not a desired outcome, you should include a uniquely identifiable columns in your query and the corresponding fields in your structs.

## Approach
carts adopts the "database mapping" approach (described in Martin Fowler's [book](https://books.google.com/books?id=FyWZt5DdvFkC&lpg=PA1&dq=Patterns%20of%20Enterprise%20Application%20Architecture%20by%20Martin%20Fowler&pg=PT187#v=onepage&q=active%20record&f=false)) which is useful among organizations with strict code review processes.

Carta is not an object-relational mapper(ORM). With large and complex datasets, using ORMs becomes restrictive and reduces performance when working with complex queries. 

### License
Apache License
