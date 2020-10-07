FROM alpine:latest

ARG CNAB_ARM_DRIVER_VERSION
ARG CNAB_ARM_DRIVER_REPO

ENV PORTER_URL=https://cdn.porter.sh
ENV PORTER_VERSION=v0.29.1

ENV PORTER_AZURE_PLUGIN_VERSION=v0.3.0

ENV CNAB_AZURE_DRIVER_VERSION=v0.0.5.1
ENV CNAB_AZURE_DRIVER_REPO=deislabs/cnab-azure-driver

ENV CNAB_ARM_DRIVER_VERSION=$CNAB_ARM_DRIVER_VERSION
ENV CNAB_ARM_DRIVER_REPO=$CNAB_ARM_DRIVER_REPO

ENV TOOLHOME="/usr/local/bin"
ENV PORTER_HOME="${TOOLHOME}/.porter"
ENV PATH=$PATH:$PORTER_HOME

RUN apk add --update curl coreutils libc6-compat && rm -rf /var/cache/apk/*

# Install porter 
RUN mkdir -p $PORTER_HOME \
    && curl -fLo $PORTER_HOME/porter $PORTER_URL/$PORTER_VERSION/porter-linux-amd64 \
    && chmod +x $PORTER_HOME/porter \
    && cp $PORTER_HOME/porter $PORTER_HOME/porter-runtime \
    && echo Installed `porter version`

# Install CNAB azure driver
RUN curl "https://github.com/${CNAB_AZURE_DRIVER_REPO}/releases/download/${CNAB_AZURE_DRIVER_VERSION}/cnab-azure-linux-amd64" -fLo "${TOOLHOME}/cnab-azure" \
    && chmod +x "${TOOLHOME}/cnab-azure"

# Install CNAB ARM driver
RUN curl "https://github.com/${CNAB_ARM_DRIVER_REPO}/releases/download/${CNAB_ARM_DRIVER_VERSION}/cnabarmdriver-linux-amd64" -fLo "${TOOLHOME}/cnabarmdriver" \
    && chmod +x "${TOOLHOME}/cnabarmdriver"

# Install Porter Azure plugin
RUN mkdir -p $PORTER_HOME/plugins/azure \
    && curl -fLo $PORTER_HOME/plugins/azure/azure $PORTER_URL/plugins/azure/$PORTER_AZURE_PLUGIN_VERSION/azure-linux-amd64 \
    && chmod +x $PORTER_HOME/plugins/azure/azure \
    && cp $PORTER_HOME/plugins/azure/azure $PORTER_HOME/plugins/azure/azure-runtime \
    && porter plugin show azure
COPY config.toml $PORTER_HOME/

CMD cnabarmdriver