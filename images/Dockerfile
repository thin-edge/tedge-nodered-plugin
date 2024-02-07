FROM ghcr.io/thin-edge/tedge-demo-main-systemd:latest

ARG TEST_USER=iotadmin
RUN useradd -ms /bin/bash "${TEST_USER}" && echo "${TEST_USER}:${TEST_USER}" | chpasswd && adduser "${TEST_USER}" sudo

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
        curl

# Add custom sm-plugins
COPY src/sm-plugin/nodered /etc/tedge/sm-plugins/
