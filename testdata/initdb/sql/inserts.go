package sql

const Inserts = `
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
