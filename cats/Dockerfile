# copied from https://github.com/alphagov/paas-docker-cloudfoundry-tools/blob/master/cf-acceptance-tests/Dockerfile

FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    curl \
    openssh-client \
    unzip \
    git \
  && rm -rf /var/lib/apt/lists/*

 ENV GOPATH /go
 ENV PATH /go/bin:/usr/local/go/bin:$PATH
 RUN \
  wget https://storage.googleapis.com/golang/go1.10.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go1.10.linux-amd64.tar.gz -C /usr/local && \
  mkdir $GOPATH && \
  rm -rf /tmp/*

RUN go get github.com/tools/godep
RUN go get github.com/onsi/ginkgo/ginkgo

# Install the cf CLI
RUN wget -q -O cf.deb "https://cli.run.pivotal.io/stable?release=debian64&version=6.35.2&source=github-rel" && \
    dpkg -i cf.deb && \
    rm -f cf.deb

