FROM ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y git
COPY build/current/linux/amd64/prog1 /prog1
CMD /prog1