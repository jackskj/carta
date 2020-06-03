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

{{end}}
`
