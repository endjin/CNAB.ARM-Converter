FROM alpine:latest

ARG CNAB_ARM_DRIVER_VERSION
ARG CNAB_ARM_DRIVER_REPO

ENV PORTER_URL=https://cdn.deislabs.io/porter
ENV PORTER_VERSION=v0.17.1-beta.1

ENV CNAB_AZURE_DRIVER_VERSION=v.0.0.3
ENV CNAB_AZURE_DRIVER_REPO=simongdavies/cnab-azure-driver

ENV CNAB_ARM_DRIVER_VERSION=$CNAB_ARM_DRIVER_VERSION
ENV CNAB_ARM_DRIVER_REPO=$CNAB_ARM_DRIVER_REPO

ENV TOOLHOME="/usr/local/bin"
ENV PORTER_HOME="/usr/local/bin/.porter"

RUN apk add --update curl coreutils && rm -rf /var/cache/apk/*

# Install porter 

RUN curl "${PORTER_URL}/${PORTER_VERSION}/porter-linux-amd64" -fLo "${TOOLHOME}/porter"
RUN chmod +x "${TOOLHOME}/porter"
RUN cp "${TOOLHOME}/porter" "${TOOLHOME}/porter-runtime"

# Install CNAB azure driver

RUN curl "https://github.com/${CNAB_AZURE_DRIVER_REPO}/releases/download/${CNAB_AZURE_DRIVER_VERSION}/cnab-azure-linux-amd64" -fLo "${TOOLHOME}/cnab-azure"
RUN chmod +x "${TOOLHOME}/cnab-azure"

# Install CNAB ARM driver
RUN curl "https://github.com/${CNAB_ARM_DRIVER_REPO}/releases/download/${CNAB_ARM_DRIVER_VERSION}/cnabarmdriver-linux-amd64" -fLo "${TOOLHOME}/cnabarmdriver"
RUN chmod +x "${TOOLHOME}/cnabarmdriver"

CMD ["/bin/bash", "-c", "cnabarmdriver"]