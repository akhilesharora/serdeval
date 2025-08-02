#!/bin/bash

# SerdeVal Deployment Script Template
# Copy this to deploy.sh and fill in your specific values

set -e

echo "🚀 SerdeVal Deployment Script"
echo "=================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Please run as root"
    exit 1
fi

# Variables - CUSTOMIZE THESE FOR YOUR DEPLOYMENT
SERVICE_NAME="serdeval"
SERVICE_DIR="/opt/serdeval"
SERVICE_PORT="YOUR_PORT_HERE"  # e.g., 9001
DOMAIN="YOUR_DOMAIN_HERE"       # e.g., freedatavalidator.xyz
SERVER_IP="YOUR_SERVER_IP"      # Your server's IP address

echo "📁 Creating service directory..."
mkdir -p $SERVICE_DIR/web/static

echo "🔧 Setting up service user and permissions..."
# Create service user if doesn't exist
if ! id "serdeval" &>/dev/null; then
    useradd --system --no-create-home --shell /bin/false serdeval
fi

echo "📋 Copying application files..."
# Copy binary (you'll need to upload this first)
if [ -f "./serdeval" ]; then
    cp ./serdeval $SERVICE_DIR/
    chmod +x $SERVICE_DIR/serdeval
else
    echo "⚠️  serdeval binary not found. Please upload it first."
fi

# Copy web files
if [ -f "./web/static/index.html" ]; then
    cp -r ./web/static/* $SERVICE_DIR/web/static/
else
    echo "⚠️  Web files not found. Please upload them first."
fi

# Set ownership
chown -R serdeval:serdeval $SERVICE_DIR

echo "📄 Creating systemd service..."
cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=SerdeVal - Privacy-focused data format validator
After=network.target
Wants=network.target

[Service]
Type=simple
User=serdeval
Group=serdeval
WorkingDirectory=$SERVICE_DIR
ExecStart=$SERVICE_DIR/serdeval web --port $SERVICE_PORT
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$SERVICE_DIR

# Environment
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

echo "🌐 Creating nginx configuration..."
cat > /etc/nginx/sites-available/$DOMAIN << EOF
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;

    # Security headers for privacy
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self' 'unsafe-inline' 'unsafe-eval'; frame-ancestors 'none';" always;
    
    # No access logs for privacy
    access_log off;
    error_log /var/log/nginx/$DOMAIN.error.log error;

    location / {
        proxy_pass http://127.0.0.1:$SERVICE_PORT;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Privacy: minimize logging
        proxy_set_header X-Forwarded-For "";
    }

    # Security - deny access to hidden files
    location ~ /\. {
        deny all;
    }
}
EOF

echo "🔗 Enabling nginx site..."
ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/

echo "🧪 Testing nginx configuration..."
nginx -t

if [ $? -eq 0 ]; then
    echo "✅ Nginx configuration is valid"
    systemctl reload nginx
else
    echo "❌ Nginx configuration has errors"
    exit 1
fi

echo "🎯 Starting SerdeVal service..."
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

# Wait a moment for service to start
sleep 3

echo "📊 Service status:"
systemctl status $SERVICE_NAME --no-pager

echo ""
echo "🎉 Deployment complete!"
echo "=================================="
echo "Service: $SERVICE_NAME"
echo "Port: $SERVICE_PORT"
echo "Domain: $DOMAIN"
echo "Directory: $SERVICE_DIR"
echo ""
echo "📝 Next steps:"
echo "1. Set up DNS in Cloudflare: A record @ -> $SERVER_IP"
echo "2. Run: certbot --nginx -d $DOMAIN -d www.$DOMAIN"
echo "3. Test: curl http://localhost:$SERVICE_PORT"
echo ""
echo "🔍 Monitor logs with:"
echo "journalctl -u $SERVICE_NAME -f"
echo ""
echo "🛠️ Manage service:"
echo "systemctl start/stop/restart/status $SERVICE_NAME"