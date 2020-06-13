package testdata

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type Blog struct {
	BlogId    int     `db:"blog_id" json:"blog_id,omitempty"`
	BlogTitle string  `db:"blog_title" json:"blog_title,omitempty"`
	Author    *Author `db:"author" json:"author,omitempty"`
	Posts     *[]Post `db:"posts" json:"posts,omitempty"`
}

type Author struct {
	AuthorId               int    `db:"author_id" json:"author_id,omitempty"`
	AuthorUsername         string `db:"author_username" json:"author_username,omitempty"`
	AuthorPassword         string `db:"author_password" json:"author_password,omitempty"`
	AuthorEmail            string `db:"author_email" json:"author_email,omitempty"`
	AuthorBio              string `db:"author_bio" json:"author_bio,omitempty"`
	AuthorFavouriteSection string `db:"author_favourite_section" json:"author_favourite_section,omitempty"`
}

type Post struct {
	PostId        int                  `db:"post_id" json:"post_id,omitempty"`
	PostBlogId    int                  `db:"post_blog_id" json:"post_blog_id,omitempty"`
	PostAuthorId  int                  `db:"post_author_id" json:"post_author_id,omitempty"`
	PostCreatedOn *timestamp.Timestamp `db:"post_created_on" json:"post_created_on,omitempty"`
	PostSection   string               `db:"post_section" json:"post_section,omitempty"`
	PostSubject   string               `db:"post_subject" json:"post_subject,omitempty"`
	Draft         string               `db:"draft" json:"draft,omitempty"`
	PostBody      string               `db:"post_body" json:"post_body,omitempty"`
	Comments      []*Comment           `db:"comments" json:"comments,omitempty"`
	Tags          []*Tag               `db:"tags" json:"tags,omitempty"`
}

type Comment struct {
	CommentId     *int           `db:"comment_id" json:"comment_id,omitempty"`
	CommentPostId *int           `db:"comment_post_id" json:"comment_post_id,omitempty"`
	CommentText   sql.NullString `db:"comment_text" json:"comment_text,omitempty"`
}

type Tag struct {
	TagId   int    `db:"tag_id" json:"tag_id,omitempty"`
	TagName string `db:"tag_name" json:"tag_name,omitempty"`
}

var BlogQuery = `
select
        B.id                as  blog_id,
        B.title             as  blog_title,
        A.id                as  author_id,
        A.username          as  author_username,
        A.password          as  author_password,
        A.email             as  author_email,
        A.bio               as  author_bio,
        A.favourite_section as  author_favourite_section,
        P.id                as  post_id,
        P.blog_id           as  post_blog_id,
        P.author_id         as  post_author_id,
        P.created_on        as  post_created_on,
        P.section           as  post_section,
        P.subject           as  post_subject,
        P.draft             as  draft,
        P.body              as  post_body,
        C.id                as  comment_id,
        C.post_id           as  comment_post_id,
        C.comment           as  comment_text,
        T.id                as  tag_id,
        T.name              as  tag_name
from blog B
        left outer join author A    on  B.author_id = A.id
        left outer join post P      on  B.id = P.blog_id
        left outer join comment C   on  P.id = C.post_id
        left outer join post_tag PT on  PT.post_id = P.id
        left outer join tag T       on  PT.tag_id = T.id
        where B.id in (1,2,3) 
order by 
        B.id, A.id, P.id, P.Id, C.id, T.id
`
