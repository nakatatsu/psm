FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
  squid \
  curl \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY squid.conf /etc/squid/squid.conf
COPY allowed-domains.txt /etc/squid/allowed-domains.txt

EXPOSE 3128

CMD ["squid", "-N", "-d", "1"]
