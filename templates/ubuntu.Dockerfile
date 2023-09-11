FROM {{ .Image }}

USER root

RUN apt-get update -y && \
  DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
  linux-image-virtual \
  initramfs-tools \
  systemd-sysv \
  systemd \
{{- if .Grub }}
  grub2 \
{{- end }}
  dbus \
  isc-dhcp-client \
  iproute2 \
  iputils-ping && \
  find /boot -type l -exec rm {} \;

{{- if not .Grub }}
RUN mv $(find /boot -name 'vmlinuz-*') /boot/vmlinuz && \
      mv $(find /boot -name 'initrd.img-*') /boot/initrd.img
{{- end }}

RUN systemctl preset-all

{{ if .Password }}RUN echo "root:{{ .Password }}" | chpasswd {{ end }}

{{ if eq .NetworkManager "netplan" }}
RUN apt install -y netplan.io
RUN mkdir -p /etc/netplan && printf '\
network:\n\
  version: 2\n\
  renderer: networkd\n\
  ethernets:\n\
    eth0:\n\
      dhcp4: true\n\
      dhcp-identifier: mac\n\
      nameservers:\n\
        addresses:\n\
        - 8.8.8.8\n\
        - 8.8.4.4\n\
' > /etc/netplan/00-netcfg.yaml
{{ else if eq .NetworkManager "ifupdown"}}
RUN if [ -z "$(apt-cache madison ifupdown-ng 2> /dev/nul)" ]; then apt install -y ifupdown; else apt install -y ifupdown-ng; fi
RUN mkdir -p /etc/network && printf '\
auto eth0\n\
allow-hotplug eth0\n\
iface eth0 inet dhcp\n\
' > /etc/network/interfaces
{{ end }}

{{- if .Luks }}
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends cryptsetup-initramfs && \
    update-initramfs -u -v
{{- end }}
