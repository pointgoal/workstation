<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [workstation](#workstation)
  - [API](#api)
    - [Organizations](#organizations)
    - [Projects](#projects)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# workstation
[![build](https://github.com/pointgoal/workstation/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/pointgoal/workstation/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pointgoal/workstation)](https://goreportcard.com/report/github.com/pointgoal/workstation)
[![codecov](https://codecov.io/gh/pointgoal/workstation/branch/main/graph/badge.svg?token=4L3ZS1E16P)](https://codecov.io/gh/pointgoal/workstation)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Workstation backend.

## API
### Organizations
| API | Description |
| --- | --- |
| GET /v1/org | List organizations |
| PUT /v1/org | Create organization |
| GET /v1/org/{orgId} | Get organization |
| POST /v1/org/{orgId} | Update organization |
| DELETE /v1/org/{orgId} | Delete organization |

### Projects
| API | Description |
| --- | --- |
| GET /v1/org/{orgId}/proj | List projects |
| PUT /v1/org/{orgId}/proj | Create project |
| GET /v1/org/{orgId}/proj/{projId} | Get project |
| POST /v1/org/{orgId}/proj/{projId} | Update project |
| DELETE /v1/org/{orgId}/proj/{projId} | Delete project |
