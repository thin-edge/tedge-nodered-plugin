FROM ghcr.io/thin-edge/tedge-demo-main-systemd:20241115.1635
ARG TEST_USER=iotadmin

# Install node-red
RUN curl -sL https://raw.githubusercontent.com/node-red/linux-installers/master/deb/update-nodejs-and-nodered > /tmp/install-nodered.sh \
    && chmod +x /tmp/install-nodered.sh \
    && sudo -u ${TEST_USER} /tmp/install-nodered.sh --confirm-install --nodered-user=${TEST_USER} || true \
    && systemctl enable nodered.service \
    && rm -f /tmp/install-nodered.sh

# Install additional dependencies
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y --no-install-recommends install \
        git \
        jq \
        curl \
    && tedge config unset c8y.proxy.client.host \
    && tedge config unset mqtt.client.host \
    && tedge config unset http.client.host

COPY dist/*.deb /tmp/
RUN case ${TARGETPLATFORM} in \
        "linux/386")  PKG_ARCH=linux_386  ;; \
        "linux/amd64")  PKG_ARCH=linux_amd64  ;; \
        "linux/arm64")  PKG_ARCH=linux_arm64  ;; \
        "linux/arm/v6")  PKG_ARCH=linux_armv6  ;; \
        "linux/arm/v7")  PKG_ARCH=linux_armv7  ;; \
        *) echo "Unsupported target platform: TARGETPLATFORM=$TARGETPLATFORM"; exit 1 ;; \
    esac \
    && dpkg -i --force-overwrite /tmp/*${PKG_ARCH}*.deb \
    && mkdir -p /opt/packages \
    && cp /tmp/*${PKG_ARCH}*.deb /opt/packages/ \
    && rm -f /tmp/*.deb
