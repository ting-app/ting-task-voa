# ting-task-voa [![Build](https://github.com/ting-app/ting-task-voa/actions/workflows/Build.yml/badge.svg?branch=main)](https://github.com/ting-app/ting-task-voa/actions/workflows/Build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/ting-app/ting-task-voa)](https://goreportcard.com/report/github.com/ting-app/ting-task-voa) [![codecov](https://codecov.io/gh/ting-app/ting-task-voa/branch/main/graph/badge.svg?token=SLBRCMSYZS)](https://codecov.io/gh/ting-app/ting-task-voa)
A scheduled job that saves [VOA](https://learningenglish.voanews.com/) as ting.

## Getting started
Run with docker:

```sh
docker run -e DB_USER_NAME=user name of MySQL database \
  -e DB_PASSWORD=password of MySQL user \
  -e DB_HOST=host of MySQL database \
  -e DB_PORT=port of MySQL database \
  -e ENABLE_SENTRY=true \
  -e SENTRY_DSN=your sentry dsn \
  -d xiaodanmao/ting-task-voa:latest
```

## License
[MIT](LICENSE)
