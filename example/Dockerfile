# Pull loggo to use it in the main container
FROM ghcr.io/renbou/loggo:latest as loggo

# Basic installation of a python poetry app
FROM python:3.11.7-slim-bookworm

ENV \
  PYTHONUNBUFFERED=1 \
  PIP_DISABLE_PIP_VERSION_CHECK=on \
  PIP_DEFAULT_TIMEOUT=100 \
  POETRY_VIRTUALENVS_CREATE=false
WORKDIR /app

RUN \
  pip install -U pip setuptools && \
  pip install poetry==1.7.1

COPY pyproject.toml poetry.lock app.py .
RUN poetry install --no-interaction --no-ansi

# Pipe the output of the main app into the loggo client, pigeon,
# which will send logs to the HQ located at the "loggo" service from compose.yml
USER nobody
COPY --from=loggo /bin/loggo /bin/loggo
CMD poetry run app | loggo pigeon --hq.addr loggo:20081
