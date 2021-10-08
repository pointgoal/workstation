<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [workstation](#workstation)
  - [Quick start](#quick-start)
  - [Backend repository](#backend-repository)
    - [MySql](#mysql)
    - [local](#local)
    - [memory](#memory)
  - [API](#api)
    - [Organizations](#organizations)
      - [List organizations](#list-organizations)
      - [Create organization](#create-organization)
      - [Get organization](#get-organization)
      - [Update organization](#update-organization)
      - [Delete organization](#delete-organization)
    - [Projects](#projects)
      - [List projects](#list-projects)
      - [Create project](#create-project)
      - [Get project](#get-project)
      - [Update project](#update-project)
      - [Delete project](#delete-project)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# workstation
[![build](https://github.com/pointgoal/workstation/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/pointgoal/workstation/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pointgoal/workstation)](https://goreportcard.com/report/github.com/pointgoal/workstation)
[![codecov](https://codecov.io/gh/pointgoal/workstation/branch/main/graph/badge.svg?token=4L3ZS1E16P)](https://codecov.io/gh/pointgoal/workstation)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Workstation backend.

## Quick start
Start workstation by running bellow command or right click main.go file on your IDE.
```shell script
$ go run main.go
```

## Backend repository
Currently, we support three types of repository which are localFs, memory and MySql.

### MySql
Configure workstation to use mysql as backend repository

- boot.yaml
```yaml
---
...
repository:
  enabled: true
  provider: mySql
  mySql:
    user: root
    pass: pass
    protocol: tcp
    addr: localhost:3306
    params:
      - "charset=utf8mb4"
      - "parseTime=True"
      - "loc=Local"
```

### local
Configure workstation to use local file system as backend repository

- boot.yaml
```yaml
---
...
repository:
  enabled: true
  provider: local
  local:
    dataDir: /Users/dongxuny/workspace/workstation/workstation/tmp
```

### memory
Configure workstation to use local memory as backend repository

- boot.yaml
```yaml
---
...
repository:
  enabled: true
  provider: memory
```

## API
### Organizations
| API | Description |
| --- | --- |
| GET /v1/org | List organizations |
| PUT /v1/org | Create organization |
| GET /v1/org/{orgId} | Get organization |
| POST /v1/org/{orgId} | Update organization |
| DELETE /v1/org/{orgId} | Delete organization |

#### List organizations
```shell script
$ curl -X GET "http://localhost:8080/v1/org"
{
  "orgList": [
    {
      "meta": {
        "id": 1,
        "createdAt": "2021-10-08T00:48:12.523+08:00",
        "updatedAt": "2021-10-08T00:48:12.523+08:00",
        "name": "org-1"
      },
      "projIds": [
        1,
        2
      ]
    },
    {
      "meta": {
        "id": 2,
        "createdAt": "2021-10-08T00:48:19.742+08:00",
        "updatedAt": "2021-10-08T00:48:19.742+08:00",
        "name": "org-2"
      },
      "projIds": []
    }
  ]
}
```

#### Create organization
```shell script
$ curl -X PUT "http://localhost:8080/v1/org?orgName=my-org-5"
{
  "orgId": 3
}
```

#### Get organization
```shell script
$ curl -X GET "http://localhost:8080/v1/org/1"
{
  "org": {
    "meta": {
      "id": 1,
      "createdAt": "2021-10-08T00:48:12.523+08:00",
      "updatedAt": "2021-10-08T00:48:12.523+08:00",
      "name": "org-1"
    },
    "projIds": [
      1,
      2
    ]
  }
}
```

#### Update organization
```shell script
$ curl -X POST "http://localhost:8080/v1/org/4" -d "{  \"name\": \"my-new-org-4\"}"
{
  "status": true
}
```

#### Delete organization
```shell script
$ curl -X DELETE "http://localhost:8080/v1/org/4"
{
  "status": true
}
```

### Projects
| API | Description |
| --- | --- |
| GET /v1/org/{orgId}/proj | List projects |
| PUT /v1/org/{orgId}/proj | Create project |
| GET /v1/org/{orgId}/proj/{projId} | Get project |
| POST /v1/org/{orgId}/proj/{projId} | Update project |
| DELETE /v1/org/{orgId}/proj/{projId} | Delete project |

#### List projects
```shell script
$ curl -X GET "http://localhost:8080/v1/org/1/proj"
{
  "projList": [
    {
      "meta": {
        "id": 1,
        "createdAt": "2021-10-08T00:49:07.928+08:00",
        "updatedAt": "2021-10-08T00:49:07.928+08:00",
        "orgId": 1,
        "name": "proj-1"
      }
    },
    {
      "meta": {
        "id": 2,
        "createdAt": "2021-10-08T00:50:09.859+08:00",
        "updatedAt": "2021-10-08T00:50:09.859+08:00",
        "orgId": 1,
        "name": "proj-2"
      }
    }
  ]
}
```

#### Create project
```shell script
$ curl -X PUT "http://localhost:8080/v1/org/3/proj" -d "{  \"name\": \"my-proj-4\"}"
{
  "orgId": 3,
  "projId": 3
}
```

#### Get project
```shell script
$ curl -X GET "http://localhost:8080/v1/org/3/proj/3"
{
  "proj": {
    "meta": {
      "id": 3,
      "createdAt": "2021-10-08T16:39:08.794+08:00",
      "updatedAt": "2021-10-08T16:39:08.794+08:00",
      "orgId": 3,
      "name": "my-proj-4"
    }
  }
}
```

#### Update project
```shell script
$ curl -X POST "http://localhost:8080/v1/org/3/proj/3" -d "{  \"name\": \"my-new-proj\"}"
{
  "status": true
}
```

#### Delete project
```shell script
$ curl -X DELETE "http://localhost:8080/v1/org/3/proj/3"
{
  "status": true
}
```