<script setup>
import { computed } from 'vue'

const props = defineProps({
  loading: Boolean,
  metrics: Array,
  selectedInstances: Array,
  regionResults: Array,
  regionOptions: Array,
})

const emit = defineEmits(['query', 'toggle-instance', 'select-all', 'clear-selection'])

const allInstances = computed(() => {
  const instances = []
  if (props.regionResults) {
    for (const region of props.regionResults) {
      if (region.instances) {
        for (const inst of region.instances) {
          instances.push({
            instanceId: inst.instanceId,
            instanceName: inst.instanceName,
            regionId: inst.regionId,
            status: inst.status
          })
        }
      }
    }
  }
  return instances
})

const isInstanceSelected = (instanceId, region) => {
  return props.selectedInstances?.some(
    i => i.instanceId === instanceId && i.region === region
  )
}

const getRegionLabel = (regionValue) => {
  const opt = props.regionOptions?.find(r => r.value === regionValue)
  return opt ? opt.label : regionValue
}

const formatMetricValue = (value, unit = '') => {
  if (value === null || value === undefined) return '-'
  if (typeof value === 'number') {
    if (value < 0.01) return '< 0.01' + unit
    return value.toFixed(2) + unit
  }
  return value + unit
}

const getMetricStatus = (value, threshold = 80) => {
  if (value === null || value === undefined) return 'unknown'
  if (value >= threshold) return 'warning'
  if (value >= 95) return 'critical'
  return 'normal'
}
</script>

<template>
  <section class="cms-monitor">
    <h2>云监控 (CloudMonitor)</h2>

    <!-- Instance Selection -->
    <div class="instance-selection">
      <div class="selection-header">
        <h3>选择监控实例</h3>
        <div class="selection-actions">
          <button class="btn btn-sm" @click="emit('select-all', allInstances)">全选</button>
          <button class="btn btn-sm" @click="emit('clear-selection')">清空选择</button>
          <span class="selected-count">已选 {{ selectedInstances?.length || 0 }} 个实例</span>
        </div>
      </div>

      <div class="instance-grid" v-if="allInstances.length > 0">
        <div
          v-for="inst in allInstances"
          :key="inst.instanceId + inst.regionId"
          class="instance-card"
          :class="{ selected: isInstanceSelected(inst.instanceId, inst.regionId) }"
          @click="emit('toggle-instance', inst.instanceId, inst.regionId)"
        >
          <div class="instance-check">
            <input
              type="checkbox"
              :checked="isInstanceSelected(inst.instanceId, inst.regionId)"
              @click.stop
              @change="emit('toggle-instance', inst.instanceId, inst.regionId)"
            />
          </div>
          <div class="instance-info">
            <span class="instance-id">{{ inst.instanceId }}</span>
            <span class="instance-name">{{ inst.instanceName || '-' }}</span>
            <span class="instance-region">{{ getRegionLabel(inst.regionId) }}</span>
          </div>
          <span class="instance-status" :class="inst.status?.toLowerCase()">
            {{ inst.status }}
          </span>
        </div>
      </div>
      <div v-else class="empty-hint">
        请先在 ECS 查询页面查询实例列表
      </div>

      <div class="query-actions">
        <button
          class="btn btn-primary"
          :disabled="loading || !selectedInstances?.length"
          @click="emit('query', selectedInstances)"
        >
          {{ loading ? '查询监控中...' : '查询监控指标' }}
        </button>
      </div>
    </div>

    <!-- Metrics Display -->
    <div v-if="metrics && metrics.length > 0" class="metrics-display">
      <h3>监控指标</h3>
      <div class="metrics-grid">
        <div
          v-for="metric in metrics"
          :key="metric.instanceId"
          class="metric-card"
        >
          <div class="metric-header">
            <span class="metric-instance-id">{{ metric.instanceId }}</span>
            <span class="metric-time" v-if="metric.updateTime">
              {{ new Date(metric.updateTime).toLocaleTimeString() }}
            </span>
          </div>

          <div class="metric-items">
            <div class="metric-item">
              <span class="metric-label">CPU 使用率</span>
              <div class="metric-bar-wrapper">
                <div
                  class="metric-bar"
                  :class="getMetricStatus(metric.cpuUtilization)"
                  :style="{ width: (metric.cpuUtilization || 0) + '%' }"
                ></div>
              </div>
              <span class="metric-value" :class="getMetricStatus(metric.cpuUtilization)">
                {{ formatMetricValue(metric.cpuUtilization, '%') }}
              </span>
            </div>

            <div class="metric-item">
              <span class="metric-label">内存使用率</span>
              <div class="metric-bar-wrapper">
                <div
                  class="metric-bar"
                  :class="getMetricStatus(metric.memoryUtilization)"
                  :style="{ width: (metric.memoryUtilization || 0) + '%' }"
                ></div>
              </div>
              <span class="metric-value" :class="getMetricStatus(metric.memoryUtilization)">
                {{ formatMetricValue(metric.memoryUtilization, '%') }}
              </span>
            </div>

            <div class="metric-item compact">
              <span class="metric-label">磁盘读取</span>
              <span class="metric-value">{{ formatMetricValue(metric.diskReadBps, ' B/s') }}</span>
            </div>

            <div class="metric-item compact">
              <span class="metric-label">磁盘写入</span>
              <span class="metric-value">{{ formatMetricValue(metric.diskWriteBps, ' B/s') }}</span>
            </div>

            <div class="metric-item compact">
              <span class="metric-label">网络入流量</span>
              <span class="metric-value">{{ formatMetricValue(metric.internetRx, ' bps') }}</span>
            </div>

            <div class="metric-item compact">
              <span class="metric-label">网络出流量</span>
              <span class="metric-value">{{ formatMetricValue(metric.internetTx, ' bps') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="!loading && metrics && metrics.length === 0" class="empty-metrics">
      暂无监控数据，请选择实例后查询
    </div>
  </section>
</template>

<style scoped>
.cms-monitor {
  padding: 20px;
}

.instance-selection {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.selection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.selection-header h3 {
  margin: 0;
  font-size: 14px;
}

.selection-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.selected-count {
  font-size: 12px;
  color: var(--text-secondary);
}

.instance-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 8px;
  margin-bottom: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.instance-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.instance-card:hover {
  border-color: var(--primary-color);
}

.instance-card.selected {
  border-color: var(--primary-color);
  background: var(--primary-bg);
}

.instance-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.instance-id {
  font-family: monospace;
  font-size: 12px;
  color: var(--primary-color);
}

.instance-name {
  font-size: 11px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.instance-region {
  font-size: 10px;
  color: var(--text-tertiary);
}

.instance-status {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.instance-status.running {
  background: var(--success-bg);
  color: var(--success-color);
}

.instance-status.stopped {
  background: var(--error-bg);
  color: var(--error-color);
}

.query-actions {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.metrics-display {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 16px;
}

.metrics-display h3 {
  margin: 0 0 16px 0;
  font-size: 14px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.metric-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.metric-instance-id {
  font-family: monospace;
  font-size: 13px;
  font-weight: 600;
  color: var(--primary-color);
}

.metric-time {
  font-size: 11px;
  color: var(--text-tertiary);
}

.metric-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.metric-item {
  display: grid;
  grid-template-columns: 80px 1fr 60px;
  align-items: center;
  gap: 8px;
}

.metric-item.compact {
  grid-template-columns: 80px 1fr;
}

.metric-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.metric-bar-wrapper {
  height: 8px;
  background: var(--bg-tertiary);
  border-radius: 4px;
  overflow: hidden;
}

.metric-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.metric-bar.normal {
  background: var(--success-color);
}

.metric-bar.warning {
  background: var(--warning-color);
}

.metric-bar.critical {
  background: var(--error-color);
}

.metric-value {
  font-size: 12px;
  font-weight: 600;
  text-align: right;
}

.metric-value.normal {
  color: var(--success-color);
}

.metric-value.warning {
  color: var(--warning-color);
}

.metric-value.critical {
  color: var(--error-color);
}

.empty-hint,
.empty-metrics {
  text-align: center;
  padding: 40px;
  color: var(--text-tertiary);
  font-size: 14px;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}
</style>
