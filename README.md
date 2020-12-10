Sonar Badge Proxy
=================
[![Build Status][Build]][Travis]
[![Lines of Code][Lines]][Sonar]
[![Coverage Status][Coverage]][Sonar]

The release of _SonarQube_ 7.1 included an [API for _Project Badges_][API] for public repositories.
[Allow usage of project badges on private projects][MMF-1178] is not yet specified or possible.

The _Sonar Badge Proxy_ enables the use of _Project Badges_ with private projects.
It provides a _reverse proxy_ to authenticate the call to the _SonarQube_ instance.


Usage
-----

URL to access a specific BADGE for a PROJECT:

    localhost:4000/$BADGE/$PROJECT


### Metric mapping

The BADGE path segment does not always match the metric name used with the [API].

    status          → alert_status
    bugs            → bugs
    codesmells      → code_smells
    coverage        → coverage
    duplications    → duplicated_lines_density
    lines           → ncloc
    maintainability → sqale_rating
    reliability     → reliability_rating
    security        → security_rating
    techdept        → sqale_index
    vulnerabilities → vulnerabilities


### Environment variables

#### PORT
The port the reverse proxy server starts on

#### AUTHORIZATION
The user token passed as Basic Authorization header

#### METRIC
A comma separated list of metrics to expose as BADGE

#### REMOTE
The host of the _SonarQube_ installation

#### SECRET
A secret to create a project access token with

#### INSECURE_SKIP_VERIFY
If "true", then skip verification of the SSL certificate of the sonarqube server

### Branch badges

To access metric badges for specific branches,
a `branch` query parameter can be added to the request.


### Project Access Token

Access to the badges provided by _Sonar Badge Proxy_ can be restricted.
The `token` should be provided as a query parameter.

    assert token == md5("$PROJECT:$SECRET")



Example
-------

Assume a _SonarQube_ _project_ on `sonarcloud.io`.
To access badges for the _bugs_ and _lines_ metrics for master publicly,
start the proxy as follows:

    #!/usr/bin/env bash
    export PORT=4000
    export REMOTE=sonarcloud.io
    export SECRET=012345789abcdef
    export METRIC=bugs,lines
    # uncomment out if you want to skip the SSL certificate check
    #export INSECURE_SKIP_VERIFY=true
    ./sonar-badge-proxy

The badges can be accessed through an URL like:

    localhost:4000/coverage/project?branch=master&token=7d9ccf5d9de733c1f7aded0048739e89


Docker
------

You can build a docker image with following command line:

    docker build . 


You can also use docker-compose to build and run _Sonar Badge Proxy_ as follows:

    docker-compose --env-file docker.env up -d


You must first edit docker.env file to configure environment variable. An example of docker.env.template is available in source. Just use it and rename it to docker.env.


License
-------
    
    Copyright (C) 2019  tynn

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.


[API]: https://next.sonarqube.com/sonarqube/web_api/api/project_badges/measure
[MMF-1178]: https://jira.sonarsource.com/browse/MMF-1178
[Build]: https://img.shields.io/travis/tynn/sonar-badge-proxy.svg?logo=travis
[Travis]: https://www.travis-ci.org/tynn/sonar-badge-proxy
[Coverage]: https://sonarcloud.io/api/project_badges/measure?project=sonar-badge-proxy&metric=coverage
[Lines]: https://sonarcloud.io/api/project_badges/measure?project=sonar-badge-proxy&metric=ncloc
[Sonar]: https://sonarcloud.io/dashboard?id=sonar-badge-proxy
