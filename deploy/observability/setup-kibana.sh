#!/bin/bash

# USDC Event Tracker - Kibana Setup Script

set -e

KIBANA_URL="http://localhost:5601"

echo "🎨 Setting up Kibana dashboards for USDC Event Tracker..."

# Wait for Kibana to be ready
echo "⏳ Waiting for Kibana to be ready..."
while ! curl -s "$KIBANA_URL/api/status" >/dev/null 2>&1; do
    echo "   Waiting for Kibana..."
    sleep 5
done

echo "✅ Kibana is ready!"

# Import index patterns
echo "📊 Importing index patterns..."
curl -X POST "$KIBANA_URL/api/saved_objects/_import" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  --form file=@./kibana/dashboards/index-patterns.ndjson

# Import dashboards
echo "📈 Importing dashboards..."
curl -X POST "$KIBANA_URL/api/saved_objects/_import" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  --form file=@./kibana/dashboards/usdc-overview-dashboard.ndjson

echo ""
echo "🎉 Kibana setup complete!"
echo ""
echo "📊 Available dashboards:"
echo "   • USDC Event Tracker - Overview Dashboard"
echo ""
echo "📈 Index patterns created:"
echo "   • usdc-transactions-*"
echo "   • usdc-blocks-*"
echo "   • usdc-sinks-*"
echo ""
echo "🔍 To explore your data:"
echo "   1. Open Kibana: $KIBANA_URL"
echo "   2. Go to 'Analytics' > 'Dashboard'"
echo "   3. Open 'USDC Event Tracker - Overview Dashboard'"
echo ""
echo "💡 Useful Kibana features for blockchain analysis:"
echo "   • Discover: Search and filter raw events"
echo "   • Visualize: Create custom charts and graphs"
echo "   • Machine Learning: Detect anomalies in transaction patterns"
echo "   • Alerts: Set up notifications for unusual activity"
echo ""