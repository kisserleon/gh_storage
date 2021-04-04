# Use Github as your storage

Use the GitHub repository as your storage by leveraging the [Github content APIs](https://docs.github.com/en/rest/reference/repos#get-repository-content).

## Prerequisites

* A Github repository
* [A Github token](https://github.com/settings/tokens/new)
* Go lang

## Usages

The typical CRUD usages.

* Help

`go run main.go -h`


```go
  -delete
        delete a file in the repository
  -dst string
        the dest file path in the repository
  -message string
        git commit message
  -retrieve
        retrieve a file or files from github
  -src string
        the file path to be uploaded/retrieved/deleted
  -update
        update a file in the repository
  -upload
        upload a file to github
```

* Upload a file

`go run main.go -upload -src=./hello.go -dst=/hello.go`

* Retrieve a file raw content

`go run main.go -retrieve -dst=/hello.go`

You can also download the file by using git.

* Retrieve file list

`go run main.go -retrieve`

* Update a file

`go run main.go -update -src=./hello.go -dst=/hello.go`

* Delete a file

`go run main.go -delete -dst=/hello.go`

## Limit

[What is my disk quota?](https://docs.github.com/en/github/managing-large-files/what-is-my-disk-quota)

```text
We recommend repositories remain small, ideally less than 1 GB, and less than 5 GB is strongly recommended. Smaller repositories are faster to clone and easier to work with and maintain. Individual files in a repository are strictly limited to a 100 MB maximum size limit. 
```

## Recommendation

It's best used to store small files.
