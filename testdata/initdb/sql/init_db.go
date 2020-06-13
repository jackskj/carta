package sql

const Init = `
{{define "InitDB" }}

drop table if exists author;
create table author (
  id                int
primary key,
  username          VARCHAR(255),
password          VARCHAR(255),
  email             VARCHAR(255),
  bio               VARCHAR(255),
  favourite_section VARCHAR(255)
);

drop table if exists blog;
create table blog (
  id        int
    primary key,
  title     VARCHAR(255),
  author_id int
);

drop table if exists comment;
create table comment (
  id int
    primary key,
  post_id int,
  name    VARCHAR(255),
  comment VARCHAR(255)
);

drop table if exists post;
create table post (
  id         int
    primary key,
  author_id  int,
  blog_id    int,
  created_on DATE,
  section    VARCHAR(255),
  subject    VARCHAR(255),
  draft      VARCHAR(255),
  body       VARCHAR(255)
);

drop table if exists  post_tag;
create table post_tag (
  post_id int,
  tag_id  int,
  constraint post_tag_pk
  primary key (post_id, tag_id)
);

drop table if exists tag;
create table tag (
  id   int
    primary key,
  name varchar(255)
);

drop table if exists relations;
create table relations ( 
  id                           int ,
  basic_submap                 int ,
  basic_submap_ptr             int ,
  basic_submap_ptr_ptr         int ,
  submap_sample_submap         int ,
  submap_ptr_sample_submap     int ,
  submap_ptr_ptr_sample_submap int
);

insert into relations values ( 1, 2, 3, 4, 5, 6, 7);
insert into relations values ( 1, 2, 3, 4, 5, 6, 7);
insert into relations values ( 1, 3, 4, 5, 6, 7, 8);
insert into relations values ( 1, 4, 5, 6, 7, 8, 9);

insert into relations values ( 2, 2, 3, 4, 5, 6, 7);
insert into relations values ( 2, 2, 3, 4, 5, 6, 7);
insert into relations values ( 2, 3, 4, 5, 6, 7, 8);
insert into relations values ( 2, 4, 5, 6, 7, 8, 9);

{{end}}


{{define "InsertAuthor" }}
INSERT INTO author
VALUES (
 {{ .Id }},
 {{ .Username | squote }},
 {{ .Password | squote }},
 {{ .Email | squote }},
 {{ .Bio | squote }},
 {{ .FavouriteSection | squote }}
);
{{end}}

{{define "InsertBlog" }}
INSERT INTO blog
VALUES (
  {{ .Id }},
 {{ .Title | squote }},
  {{ .AuthorId }}
);
{{end}}

{{define "InsertComment" }}
INSERT INTO comment
VALUES (
  {{ .Id }},
  {{ .PostId }},
  {{ .Name | squote }},
  {{ .Comment | squote }}
);
{{end}}

{{define "InsertPost" }}
INSERT INTO post
VALUES (
  {{ .Id }},
  {{ .AuthorId }},
  {{ .BlogId }},
 {{ .CreatedOn | timestamp | squote }},
 {{ .Section | squote }},
 {{ .Subject | squote }},
 {{ .Draft | squote }},
 {{ .Body | squote }}
);
{{end}}

{{define "InsertPostTag" }}
INSERT INTO post_tag
VALUES (
  {{ .PostId }},
  {{ .TagId }}
);
{{end}}

{{define "InsertTag" }}
INSERT INTO tag
VALUES (
  {{ .Id }},
 {{ .Name | squote }}
);
{{end}}
`
