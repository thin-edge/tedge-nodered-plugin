FROM ghcr.io/thin-edge/tedge-demo-main-systemd:20230517.2

# Install node-red
RUN curl -sL https://raw.githubusercontent.com/node-red/linux-installers/master/deb/update-nodejs-and-nodered > /tmp/install-nodered.sh \
    && chmod +x /tmp/install-nodered.sh \
    && sudo -u iotadmin /tmp/install-nodered.sh --confirm-install --nodered-user=iotadmin || true \
    && systemctl enable nodered.service \
    && rm -f /tmp/install-nodered.sh

# Install additional dependencies
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y --no-install-recommends install \
        git \
        jq \
        curl

# Add node-red settings (which projects feature activated)
COPY --chown=iotadmin:iotadmin images/files/settings.js /home/iotadmin/.node-red/

# Add custom sm-plugins
COPY src/sm-plugin/nodered /etc/tedge/sm-plugins/
