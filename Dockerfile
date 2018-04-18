FROM python:3-alpine
ENV DATABASE_URI "sqlite:////data/sslcp.db"
ENV PORT 80
EXPOSE 80
VOLUME /data
WORKDIR /usr/src/app
COPY requirements.txt ./
RUN apk update &&\
    apk add python3-dev postgresql-dev libffi-dev gcc musl-dev &&\
    pip install --no-cache-dir -r requirements.txt &&\
    apk del python3-dev postgresql-dev libffi-dev gcc musl-dev
COPY main.py .
COPY app/ ./app
CMD [ "python", "./main.py"]