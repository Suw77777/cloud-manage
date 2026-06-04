<script setup>
import { computed } from 'vue'
import PaginationBar from './PaginationBar.vue'

const props = defineProps({
  region: String,
  project: String,
  logstore: String,
  query: String,
  timeRange: String,
  maxLines: Number,
  logStores: Array,
  logs: Array,
  logCount: Number,
  hasMore: Boolean,
  loading: Boolean,
  loadingStores: Boolean,
  streaming: Boolean,
  streamProgress: Number,
  regionOptions: Array,
  // 分页
  pageSize: Number,
  currentPage: Number,
  totalLogs: Number,
})

const emit = defineEmits([
  'update:region',
  'update:project',
  'update:logstore',
  'update:query',
  'update:timeRange',
  'update:maxLines',
  'fetch-stores',
  'query-logs',
  'query-logs-stream',
  'page-change',
  'page-size-change',
  'export-logs',
  'clear'
])

const timeRangeOptions = [
  { label: '最近 15 分钟', value: '15m' },
  { label: '最近 30 分钟', value: '30m' },
  { label: '最近 1 小时', value: '1h' },
  { label: '最近 6 小时', value: '6h' },
  { label: '最近 12 小时', value: '12h' },
  { label: '最近 24 小时', value: '24h' },
  { label: '最近 7 天', value: '7d' },
]

const maxLinesOptions = [50, 100, 200, 500, 1000]

const formatTimestamp = (ts) => {
  if (!ts) return '-'
  const date = new Date(ts)
  return date.toLocaleString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    fractionalSecondDigits: 3
  })
}

// 分页后的日志数据
const pagedLogs = computed(() => {
  if (!props.logs) return []
  const start = ((props.currentPage || 1) - 1) * (props.pageSize || 50)
  const end = start + (props.pageSize || 50)
  return props.logs.slice(start, end)
})

const logEntries = computed(() => {
  return pagedLogs.value.map((log, index) => ({
    id: index,
    timestamp: formatTimestamp(log.timestamp),
    content: log.content || {},
    raw: Object.entries(log.content || {}).map(([k, v]) => `${k}: ${v}`).join(' | ')
  }))
})

const totalPages = computed(() => {
  if (!props.logs || !props.pageSize) return 1
  return Math.ceil(props.logs.length / props.pageSize)
})
</script>

<template>
  <section class="sls-query">
    <h2>SLS 日志查询</h2>

    <!-- Query Configuration -->
    <div class="query-config">
      <div class="form-grid">
        <div class="form-group">
          <label>Region</label>
          <select :value="region" @change="emit('update:region', $event.target.value)">
            <option v-for="r in regionOptions" :key="r.value" :value="r.value">
              {{ r.label }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>Project</label>
          <div class="input-with-btn">
            <input
              :value="project"
              @input="emit('update:project', $event.target.value)"
              placeholder="输入 SLS Project 名称"
              @keyup.enter="emit('fetch-stores')"
            />
            <button
              class="btn btn-sm"
              :disabled="loadingStores || !project"
              @click="emit('fetch-stores')"
            >
              {{ loadingStores ? '加载中...' : '获取 Logstore' }}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>Logstore</label>
          <select
            :value="logstore"
            @change="emit('update:logstore', $event.target.value)"
            :disabled="!logStores || logStores.length === 0"
          >
            <option value="">请选择 Logstore</option>
            <option v-for="ls in logStores" :key="ls.logstoreName" :value="ls.logstoreName">
              {{ ls.logstoreName }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>时间范围</label>
          <select :value="timeRange" @change="emit('update:timeRange', $event.target.value)">
            <option v-for="opt in timeRangeOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>最大行数</label>
          <select :value="maxLines" @change="emit('update:maxLines', Number($event.target.value))">
            <option v-for="n in maxLinesOptions" :key="n" :value="n">{{ n }} 行</option>
          </select>
        </div>
      </div>

      <div class="form-group full-width">
        <label>查询语句</label>
        <div class="query-input-row">
          <input
            :value="query"
            @input="emit('update:query', $event.target.value)"
            placeholder="输入查询语句，如: level: ERROR OR status: 500 (留空查询全部)"
            @keyup.enter="emit('query-logs')"
          />
          <button
            class="btn btn-primary"
            :disabled="loading || streaming || !logstore"
            @click="emit('query-logs')"
          >
            {{ loading ? '查询中...' : '查询日志' }}
          </button>
          <button
            class="btn btn-info"
            :disabled="loading || streaming || !logstore"
            @click="emit('query-logs-stream')"
          >
            {{ streaming ? '流式查询中...' : '流式查询' }}
          </button>
          <button
            class="btn btn-success"
            :disabled="loading || streaming || !logs || logs.length === 0"
            @click="emit('export-logs', 'csv')"
          >
            导出 CSV
          </button>
          <button
            class="btn btn-warning"
            :disabled="loading || streaming || !logs || logs.length === 0"
            @click="emit('export-logs', 'json')"
          >
            导出 JSON
          </button>
          <button class="btn btn-secondary" @click="emit('clear')">清空</button>
        </div>
      </div>

      <!-- Stream Progress -->
      <div v-if="streaming" class="stream-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: streamProgress + '%' }"></div>
        </div>
        <span class="progress-text">已接收 {{ logs?.length || 0 }} 条日志</span>
      </div>
    </div>

    <!-- Results Summary -->
    <div v-if="logs && logs.length > 0" class="results-summary">
      <span>内存中共 {{ logs.length }} 条日志</span>
      <span v-if="hasMore" class="has-more">（结果已截断，请缩小查询范围）</span>
    </div>

    <!-- Log Results with Virtual Scrolling -->
    <div v-if="logEntries.length > 0" class="log-results">
      <div class="log-table-wrapper">
        <table class="log-table">
          <thead>
            <tr>
              <th class="col-time">时间</th>
              <th class="col-content">日志内容</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logEntries" :key="log.id">
              <td class="col-time">{{ log.timestamp }}</td>
              <td class="col-content">
                <div class="log-content" :title="log.raw">
                  <span v-for="(value, key) in log.content" :key="key" class="log-field">
                    <span class="field-key">{{ key }}</span>
                    <span class="field-sep">=</span>
                    <span class="field-value">{{ value }}</span>
                  </span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <PaginationBar
        :current-page="currentPage || 1"
        :total-pages="totalPages"
        :total-items="logs?.length || 0"
        :page-size="pageSize || 50"
        :has-more="currentPage < totalPages"
        :has-prev="currentPage > 1"
        @page-change="emit('page-change', $event)"
        @page-size-change="emit('page-size-change', $event)"
      />
    </div>

    <div v-else-if="!loading && !streaming && logstore" class="empty-state">
      暂无日志数据，请输入查询条件后查询
    </div>
  </section>
</template>

<style scoped>
.sls-query {
  padding: 0;
}

.query-config {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
}

.form-group input,
.form-group select {
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.input-with-btn {
  display: flex;
  gap: 8px;
}

.input-with-btn input {
  flex: 1;
}

.query-input-row {
  display: flex;
  gap: 8px;
}

.query-input-row input {
  flex: 1;
}

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
  white-space: nowrap;
}

.stream-progress {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #e4e7ed;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
}

.results-summary {
  font-size: 13px;
  color: #606266;
  margin-bottom: 12px;
  padding: 8px 0;
  display: flex;
  gap: 8px;
}

.has-more {
  color: #e6a23c;
  font-weight: 500;
}

.log-results {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.log-table-wrapper {
  overflow-x: auto;
  max-height: 500px;
  overflow-y: auto;
}

.log-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
}

.log-table th {
  background: #fafafa;
  padding: 10px 12px;
  text-align: left;
  font-weight: 600;
  color: #606266;
  border-bottom: 2px solid #ebeef5;
  position: sticky;
  top: 0;
  z-index: 1;
}

.log-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #ebeef5;
  color: #333;
  vertical-align: top;
}

.log-table tr:hover td {
  background: #f5f7fa;
}

.col-time {
  width: 120px;
  white-space: nowrap;
  color: #909399;
}

.col-content {
  min-width: 400px;
}

.log-content {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 12px;
  line-height: 1.6;
}

.log-field {
  display: inline;
}

.field-key {
  color: #667eea;
  font-weight: 500;
}

.field-sep {
  color: #909399;
  margin: 0 2px;
}

.field-value {
  color: #333;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
  font-size: 14px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
</style>
