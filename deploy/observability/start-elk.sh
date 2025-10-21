#!/bin/bash

# USDC Event Tracker - ELK Stack Startup Script

set -e

echo "🚀 Starting ELK Stack for USDC Event Tracker..."

# Create logs directory if it doesn't exist
mkdir -p ../../logs

# Check if .env.elk exists, if not copy from example
if [ ! -f ../../.env.elk ]; then
    echo "📝 Creating .env.elk from example..."
    cp ../../.env.elk.example ../../.env.elk
    echo "⚠️  Please edit .env.elk with your configuration before running the tracker"
fi

# Start ELK stack
echo "📊 Starting Elasticsearch, Logstash, and Kibana..."
docker-compose -f docker-compose.elk.yml up -d

# Wait for Elasticsearch to be ready
echo "⏳ Waiting for Elasticsearch to be ready..."
while ! curl -s http://localhost:9200/_cluster/health >/dev/null 2>&1; do
    echo "   Waiting for Elasticsearch..."
    sleep 5
done

echo "✅ Elasticsearch is ready!"

# Wait for Kibana to be ready
echo "⏳ Waiting for Kibana to be ready..."
while ! curl -s http://localhost:5601/api/status >/dev/null 2>&1; do
    echo "   Waiting for Kibana..."
    sleep 5
done

echo "✅ Kibana is ready!"

echo ""
echo "🎉 ELK Stack is running!"
echo ""
echo "📊 Services:"
echo "   Elasticsearch: http://localhost:9200"
echo "   Kibana:        http://localhost:5601" 
echo "   Logstash:      http://localhost:9600"
echo ""
echo "📈 To view logs and analytics:"
echo "   1. Open Kibana: http://localhost:5601"
echo "   2. Go to 'Stack Management' > 'Index Patterns'"
echo "   3. Create index patterns for:"
echo "      - usdc-logs-*"
echo "      - usdc-transactions-*"
echo "      - usdc-blocks-*"
echo "      - usdc-sinks-*"
echo "   4. Go to 'Discover' to explore your USDC data!"
echo ""
echo "🔧 To start the USDC Event Tracker:"
echo "   cp .env.elk .env"
echo "   go run main.go"
echo ""
echo "🛑 To stop the ELK stack:"
echo "   docker-compose -f docker-compose.elk.yml down"
echo ""