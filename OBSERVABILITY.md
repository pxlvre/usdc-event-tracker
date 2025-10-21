# USDC Event Tracker - Observability with ELK Stack

This document describes the comprehensive observability setup for the USDC Event Tracker using the ELK Stack (Elasticsearch, Logstash, Kibana).

## 🎯 Overview

The observability setup provides:

- **Structured JSON Logging** - All application logs in JSON format for easy parsing
- **Elasticsearch Storage** - Direct event storage and log aggregation
- **Logstash Processing** - Log enrichment and routing
- **Kibana Visualizations** - Real-time dashboards and analytics
- **Performance Monitoring** - Sink performance, gas usage, block processing metrics

## 🚀 Quick Start

### 1. Start the ELK Stack

```bash
# Start all ELK services
./deploy/observability/start-elk.sh

# This will:
# - Create docker containers for Elasticsearch, Logstash, Kibana
# - Wait for services to be ready
# - Provide access URLs
```

### 2. Configure the USDC Tracker

```bash
# Copy ELK configuration
cp .env.elk.example .env.elk

# Edit with your settings
vi .env.elk

# Key settings:
# WEBHOOK_URL=wss://sepolia.infura.io/ws/v3/YOUR_PROJECT_ID
# SINKS=console,elasticsearch
# ELASTICSEARCH_URLS=http://localhost:9200
```

### 3. Setup Kibana Dashboards

```bash
# Import pre-built dashboards and index patterns
./deploy/observability/setup-kibana.sh
```

### 4. Start the USDC Tracker

```bash
# Run with ELK configuration
cp .env.elk .env
go run main.go
```

## 📊 Data Flow Architecture

```
Go Application → JSON Logs → Logstash → Elasticsearch → Kibana
              ↘ Direct Events → Elasticsearch Sink
```

### Logging Flow
1. **Application** generates structured JSON logs
2. **Logstash** processes and enriches logs
3. **Elasticsearch** stores processed data
4. **Kibana** provides visualization and analytics

### Event Flow
1. **USDC Events** detected by tracker
2. **Elasticsearch Sink** stores events directly
3. **Index Templates** ensure proper mapping
4. **Kibana** displays real-time data

## 🔍 Available Data & Indexes

### Index Patterns

| Index Pattern | Description | Key Fields |
|---------------|-------------|------------|
| `usdc-transactions-*` | USDC Transfer and Approval events | `tx_hash`, `from_address`, `to_address`, `value`, `gas_used` |
| `usdc-blocks-*` | Block processing metrics | `block_number`, `transaction_count`, `processing_time` |
| `usdc-sinks-*` | Sink performance data | `sink_name`, `duration_ms`, `event_count`, `success` |
| `usdc-logs-*` | General application logs | `level`, `component`, `message`, `error` |

### Key Metrics

#### Transaction Metrics
- **Volume**: Transactions per minute/hour
- **Gas Usage**: Average, min, max gas consumption
- **Event Types**: Transfer vs Approval distribution
- **Addresses**: Most active from/to addresses
- **Value Transfers**: USDC amounts being moved

#### System Metrics
- **Block Processing**: Time per block, transactions per block
- **Sink Performance**: Write latency, success rates
- **Error Rates**: Failed transactions, system errors
- **Network Health**: Connection status, RPC response times

## 📈 Pre-built Dashboards

### USDC Event Tracker - Overview Dashboard

**Panels:**
1. **Transaction Timeline** - Volume over time
2. **Gas Usage Analysis** - Average, min, max gas consumption patterns
3. **Block Processing Performance** - Transactions per block over time
4. **Transaction Types** - Transfer vs Approval pie chart
5. **Sink Performance** - Processing duration by sink type

**Use Cases:**
- Monitor real-time USDC activity
- Identify gas price trends
- Track system performance
- Alert on anomalies

## 🔧 Configuration

### Environment Variables

#### Core Settings
```bash
# Application
WEBHOOK_URL=wss://sepolia.infura.io/ws/v3/YOUR_PROJECT_ID
NETWORK=sepolia
SINKS=console,elasticsearch

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

#### Elasticsearch Settings
```bash
ELASTICSEARCH_URLS=http://localhost:9200
ELASTICSEARCH_USERNAME=
ELASTICSEARCH_PASSWORD=
ELASTICSEARCH_INDEX_PREFIX=usdc-events
ELASTICSEARCH_BATCH_SIZE=100
ELASTICSEARCH_USE_TIMESTAMP_SUFFIX=true
```

### Logstash Configuration

**Input Sources:**
- TCP port 5000 for JSON logs
- File input for log files
- Container logs via Filebeat

**Processing:**
- Event type detection and tagging
- Numeric field conversion
- Time-based field addition
- Error enrichment

**Output Routing:**
- `usdc-transactions-*` - USDC events
- `usdc-blocks-*` - Block processing
- `usdc-sinks-*` - Sink operations
- `usdc-logs-*` - General logs

## 📋 Monitoring & Alerting

### Key Metrics to Monitor

#### Business Metrics
- **Transaction Volume** - Unusual spikes or drops
- **Gas Price Patterns** - Significant increases
- **Large Transfers** - High-value USDC movements
- **Failed Transactions** - Error rate increases

#### System Metrics
- **Processing Latency** - Block processing delays
- **Sink Failures** - Data loss prevention
- **Memory Usage** - Resource monitoring
- **Connection Health** - RPC endpoint status

### Setting Up Alerts

1. **In Kibana:**
   - Go to Stack Management > Rules and Connectors
   - Create rules based on log patterns
   - Set up email/Slack notifications

2. **Example Alert Rules:**
   ```
   # High gas usage alert
   fields.gas_used > 200000 AND @timestamp > now-5m
   
   # Sink failure alert
   level: ERROR AND component: *sink* AND @timestamp > now-1m
   
   # Transaction volume drop
   Count of usdc-transactions-* < 10 in last 10 minutes
   ```

## 🛠 Advanced Features

### Custom Visualizations

Create custom Kibana visualizations for:
- **Heatmaps** - Transaction activity by hour/day
- **Network Graphs** - Address relationship mapping
- **Geographic Maps** - Transaction origin analysis
- **Anomaly Detection** - ML-based pattern recognition

### Data Retention

Configure index lifecycle management:
```json
{
  "policy": {
    "phases": {
      "hot": {"actions": {"rollover": {"max_size": "5GB"}}},
      "warm": {"min_age": "7d"},
      "cold": {"min_age": "30d"},
      "delete": {"min_age": "90d"}
    }
  }
}
```

### Performance Optimization

**Elasticsearch:**
- Adjust JVM heap size based on data volume
- Configure index sharding strategy
- Use index templates for field mapping

**Logstash:**
- Tune batch size and worker threads
- Use persistent queues for reliability
- Configure memory-based buffering

## 🔍 Troubleshooting

### Common Issues

#### Elasticsearch Connection Failures
```bash
# Check Elasticsearch health
curl http://localhost:9200/_cluster/health

# Check logs
docker logs usdc-elasticsearch
```

#### Logstash Processing Issues
```bash
# Check Logstash logs
docker logs usdc-logstash

# Test configuration
docker exec usdc-logstash logstash --config.test_and_exit
```

#### Kibana Visualization Problems
```bash
# Refresh index patterns
curl -X POST "localhost:5601/api/saved_objects/index-pattern/usdc-transactions/_refresh"

# Check Kibana logs
docker logs usdc-kibana
```

### Performance Issues

#### Slow Queries
- Check index patterns and field mappings
- Use filters instead of full-text search
- Optimize aggregation queries

#### High Memory Usage
- Adjust Elasticsearch heap size
- Configure field data circuit breaker
- Use doc values for aggregations

## 📚 Additional Resources

### Elasticsearch
- [Index Lifecycle Management](https://www.elastic.co/guide/en/elasticsearch/reference/current/index-lifecycle-management.html)
- [Mapping and Field Types](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html)

### Logstash
- [Configuration Syntax](https://www.elastic.co/guide/en/logstash/current/configuration.html)
- [Filter Plugins](https://www.elastic.co/guide/en/logstash/current/filter-plugins.html)

### Kibana
- [Dashboard Creation](https://www.elastic.co/guide/en/kibana/current/dashboard.html)
- [Alerting and Actions](https://www.elastic.co/guide/en/kibana/current/alerting-getting-started.html)

### Blockchain Analytics
- [DeFi Data Analysis Patterns](https://ethereum.org/en/developers/tutorials/)
- [Gas Price Analysis](https://ethereum.org/en/developers/docs/gas/)

---

**Need Help?** Check the logs, review configurations, and consult the official Elastic Stack documentation for advanced configurations.