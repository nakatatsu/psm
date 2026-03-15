#!/bin/bash
set -euo pipefail  # Exit on error, undefined vars, and pipeline failures
IFS=$'\n\t'       # Stricter word splitting

# ==========================================================================
# Egress control strategy:
#   - iptables allows ONLY Squid proxy (domain-based filtering)
#   - Squid controls which domains are reachable (allowed-domains.txt)
#   - No more IP-based allowlisting — no dig/ipset needed
# ==========================================================================

# 1. Extract Docker DNS info BEFORE any flushing
DOCKER_DNS_RULES=$(iptables-save -t nat | grep "127\.0\.0\.11" || true)

# Flush existing rules
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X

# 2. Selectively restore ONLY internal Docker DNS resolution
if [ -n "$DOCKER_DNS_RULES" ]; then
    echo "Restoring Docker DNS rules..."
    iptables -t nat -N DOCKER_OUTPUT 2>/dev/null || true
    iptables -t nat -N DOCKER_POSTROUTING 2>/dev/null || true
    echo "$DOCKER_DNS_RULES" | xargs -L 1 iptables -t nat
else
    echo "No Docker DNS rules to restore"
fi

# Allow Docker internal DNS only (no external DNS — prevents DNS tunneling)
iptables -A OUTPUT -d 127.0.0.11 -p udp --dport 53 -j ACCEPT
iptables -A INPUT -s 127.0.0.11 -p udp --sport 53 -j ACCEPT

# Allow localhost
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow Squid proxy — all domain-based egress control is delegated to Squid
SQUID_IP=$(getent hosts outbound-filter | awk '{print $1}')
if [ -z "$SQUID_IP" ]; then
    echo "ERROR: Failed to resolve outbound-filter container IP"
    exit 1
fi
echo "Squid proxy at: $SQUID_IP:3128"
iptables -A OUTPUT -d "$SQUID_IP" -p tcp --dport 3128 -j ACCEPT

# Allow gh-token-sidecar (direct access, not proxied)
GH_SIDECAR_IP=$(getent hosts gh-token-sidecar | awk '{print $1}')
if [ -n "$GH_SIDECAR_IP" ]; then
    echo "gh-token-sidecar at: $GH_SIDECAR_IP:80"
    iptables -A OUTPUT -d "$GH_SIDECAR_IP" -p tcp --dport 80 -j ACCEPT
else
    echo "WARNING: Failed to resolve gh-token-sidecar (may not be running)"
fi

# Allow host network (Docker host communication, e.g. VS Code remote)
HOST_IP=$(ip route | grep default | cut -d" " -f3)
if [ -z "$HOST_IP" ]; then
    echo "ERROR: Failed to detect host IP"
    exit 1
fi
HOST_NETWORK=$(echo "$HOST_IP" | sed "s/\.[0-9]*$/.0\/24/")
echo "Host network: $HOST_NETWORK"
iptables -A INPUT -s "$HOST_NETWORK" -j ACCEPT
iptables -A OUTPUT -d "$HOST_NETWORK" -j ACCEPT

# Set default policies to DROP
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# Allow established connections
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Reject all other outbound traffic for immediate feedback
iptables -A OUTPUT -j REJECT --reject-with icmp-admin-prohibited

echo "Firewall configuration complete"
echo "Verifying firewall rules..."

# 1. Direct access (bypassing proxy) should be blocked
if curl --connect-timeout 5 https://example.com >/dev/null 2>&1; then
    echo "ERROR: Direct access to example.com should be blocked"
    exit 1
else
    echo "PASS: Direct access blocked"
fi

# 2. Proxy access to allowed domain should work
if ! curl --connect-timeout 5 -x http://outbound-filter:3128 https://api.github.com/zen >/dev/null 2>&1; then
    echo "ERROR: Proxy access to api.github.com failed"
    exit 1
else
    echo "PASS: Proxy access to allowed domain works"
fi

# 3. Proxy access to disallowed domain should be blocked
if curl --connect-timeout 5 -x http://outbound-filter:3128 https://example.com >/dev/null 2>&1; then
    echo "ERROR: Proxy should block example.com"
    exit 1
else
    echo "PASS: Proxy correctly blocks disallowed domain"
fi
