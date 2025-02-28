FROM golang:1.20-alpine AS dev
WORKDIR /cli
RUN apk update \
    && apk --no-cache --update add git build-base
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o debricked ./cmd/debricked
ENTRYPOINT ["debricked"]

FROM alpine:latest AS cli
ENV DEBRICKED_TOKEN=""
RUN apk add git
WORKDIR /root/
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
ENTRYPOINT [ "debricked",  "scan" ]

FROM cli AS resolution
RUN apk --no-cache --update add \
    openjdk8-jre \
    python3 \
    py3-scipy \
    py3-pip \
    go~=1.20

ENV MAVEN_VERSION 3.9.0
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

ENV GRADLE_VERSION 8.0.2
ENV GRADLE_HOME /usr/lib/gradle
ENV PATH $GRADLE_HOME/gradle-$GRADLE_VERSION/bin:$PATH
RUN wget https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-bin.zip && \
  unzip gradle-$GRADLE_VERSION-bin.zip -d $GRADLE_HOME && \
  rm gradle-$GRADLE_VERSION-bin.zip \
