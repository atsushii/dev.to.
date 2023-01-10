# dev.to
# Introduction

This repo can manage your dev.to articles. Write down your article on your favorite editor and post your article to dev.to via GitHub Actions.

# Setup Env

### Github Actions

Create a secret for GHA. Set your API key which was generated on [dev.to](https://dev.to/)
This setup is necessary to post your article automatically via GHA.


```
DEV_TO_GIT_TOKEN=yourApiKey
```

### Local

Set API key as env. It is used for building your draft article by golang. Check the below Building section.
This setup is necessary to build your draft article.

```
# copy env
cp .envrc .env
```


```
# set your API key in .env 
DEV_TO_GIT_TOKEN=yourApiKey
```


# Building

You can just run below to create draft article and dir.

```
# Create default draft artilcle
go run main.go

```

