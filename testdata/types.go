package testdata

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Blog struct {
	BlogId    uint32  `db:"blog_id" json:"blog_id,omitempty"`
	BlogTitle string  `db:"blog_title" json:"blog_title,omitempty"`
	Author    *Author `db:"author" json:"author,omitempty"`
	Posts     []*Post `db:"posts" json:"posts,omitempty"`
}

type Author struct {
	AuthorId               uint32 `db:"author_id" json:"author_id,omitempty"`
	AuthorUsername         string `db:"author_username" json:"author_username,omitempty"`
	AuthorPassword         string `db:"author_password" json:"author_password,omitempty"`
	AuthorEmail            string `db:"author_email" json:"author_email,omitempty"`
	AuthorBio              string `db:"author_bio" json:"author_bio,omitempty"`
	AuthorFavouriteSection string `db:"author_favourite_section" json:"author_favourite_section,omitempty"`
}

type Post struct {
	PostId        uint32               `db:"post_id" json:"post_id,omitempty"`
	PostBlogId    uint32               `db:"post_blog_id" json:"post_blog_id,omitempty"`
	PostAuthorId  uint32               `db:"post_author_id" json:"post_author_id,omitempty"`
	PostCreatedOn *timestamp.Timestamp `db:"post_created_on" json:"post_created_on,omitempty"`
	PostSection   string               `db:"post_section" json:"post_section,omitempty"`
	PostSubject   string               `db:"post_subject" json:"post_subject,omitempty"`
	Draft         string               `db:"draft" json:"draft,omitempty"`
	PostBody      string               `db:"post_body" json:"post_body,omitempty"`
	Comments      []*Comment           `db:"comments" json:"comments,omitempty"`
	Tags          []*Tag               `db:"tags" json:"tags,omitempty"`
}

type Comment struct {
	CommentId     uint32 `db:"comment_id" json:"comment_id,omitempty"`
	CommentPostId uint32 `db:"comment_post_id" json:"comment_post_id,omitempty"`
	CommentName   string `db:"comment_name" json:"comment_name,omitempty"`
	CommentText   string `db:"comment_text" json:"comment_text,omitempty"`
}

type Tag struct {
	TagId   uint32 `db:"tag_id" json:"tag_id,omitempty"`
	TagName string `db:"tag_name" json:"tag_name,omitempty"`
}
