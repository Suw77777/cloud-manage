<script setup>
import { ref, computed } from 'vue'
import { QueryECSMultiRegion } from '../wailsjs/go/main/App'

const accessKeyId = ref('')
const accessKeySecret = ref('')
const selectedRegions = ref(['cn-hangzhou'])
const env = ref('dev')
const regionResults = ref([])
const errorMessage = ref('')
const successMessage = ref('')
const loading = ref(false)

const envOptions = [
  { label: '开发环境 (dev)', value: 'dev' },
  { label: '预发布环境 (pre)', value: 'pre' },
  { label: '生产环境 (prod)', value: 'prod' },
]

const regionOptions = [
  { label: '华东1 杭州', value: 'cn-hangzhou' },
  { label: '华东2 上海', value: 'cn-shanghai' },
  { label: '华北1 青岛', value: 'cn-qingdao' },
  { label: '华北2 北京', value: 'cn-beijing' },
  { label: '华北3 张家口', value: 'cn-zhangjiakou' },
  { label: '华北5 呼和浩特', value: 'cn-huhehaote' },
  { label: '华北6 乌兰察布', value: 'cn-wulanchabu' },
  { label: '华南1 深圳', value: 'cn-shenzhen' },
  { label: '华南2 河源', value: 'cn-heyuan' },
  { label: '华南3 广州', value: 'cn-guangzhou' },
  { label: '西南1 成都', value: 'cn-chengdu' },
  { label: '中国香港', value: 'cn-hongkong' },
  { label: '新加坡', value: 'ap-southeast-1' },
  { label: '东京', value: 'ap-northeast-1' },
  { label: '弗吉尼亚', value: 'us-east-1' },
  { label: '硅谷', value: 'us-west-1' },
  { label: '法兰克福', value: 'eu-central-1' },
  { label: '伦敦', value: 'eu-west-1' },
]

const totalInstances = computed(() => {
  return regionResults.value.reduce((sum, r) => sum + (r.instances ? r.instances.length : 0), 0)
})

function toggleRegion(regionValue) {
  const idx = selectedRegions.value.indexOf(regionValue)
  if (idx === -1) {
    selectedRegions.value.push(regionValue)
  } else {
    selectedRegions.value.splice(idx, 1)
  }
}

async function queryECS() {
  errorMessage.value = ''
  successMessage.value = ''
  regionResults.value = []

  if (!accessKeyId.value.trim()) {
    errorMessage.value = '请输入 AccessKey ID'
    return
  }
  if (!accessKeySecret.value.trim()) {
    errorMessage.value = '请输入 AccessKey Secret'
    return
  }
  if (selectedRegions.value.length === 0) {
    errorMessage.value = '请至少选择一个 Region'
    return
  }

  loading.value = true
  try {
    const result = await QueryECSMultiRegion(
      accessKeyId.value.trim(),
      accessKeySecret.value.trim(),
      selectedRegions.value,
      env.value
    )
    if (result.success) {
      regionResults.value = result.regions || []
      successMessage.value = result.message
    } else {
      errorMessage.value = result.message
    }
  } catch (err) {
    errorMessage.value = '调用失败: ' + (err.message || err)
  } finally {
    loading.value = false
  }
}

function clearInputs() {
  accessKeyId.value = ''
  accessKeySecret.value = ''
  selectedRegions.value = ['cn-hangzhou']
  env.value = 'dev'
  errorMessage.value = ''
  successMessage.value = ''
}

function clearResults() {
  regionResults.value = []
  errorMessage.value = ''
  successMessage.value = ''
}

function getStatusClass(status) {
  const s = (status || '').toLowerCase()
  if (s === 'running') return 'status-running'
  if (s === 'stopped') return 'status-stopped'
  if (s === 'starting') return 'status-starting'
  if (s === 'stopping') return 'status-stopping'
  return 'status-default'
}

function getRegionLabel(value) {
  const found = regionOptions.find(r => r.value === value)
  return found ? found.label : value
}
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>Cloud 管理小助手</h1>
      <span class="version">v0.0.2</span>
    </header>

    <main class="app-main">
      <!-- Input Section -->
      <section class="input-section">
        <h2>连接配置</h2>
        <div class="form-grid">
          <div class="form-group">
            <label for="accessKeyId">AccessKey ID</label>
            <input
              id="accessKeyId"
              v-model="accessKeyId"
              type="text"
              placeholder="请输入 AccessKey ID"
              autocomplete="off"
            />
          </div>
          <div class="form-group">
            <label for="accessKeySecret">AccessKey Secret</label>
            <input
              id="accessKeySecret"
              v-model="accessKeySecret"
              type="password"
              placeholder="请输入 AccessKey Secret"
              autocomplete="off"
            />
          </div>
          <div class="form-group">
            <label for="env">环境</label>
            <select id="env" v-model="env">
              <option v-for="opt in envOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>

        <!-- Region Multi-Select -->
        <div class="region-section">
          <label>Region（可多选）</label>
          <div class="region-grid">
            <label
              v-for="r in regionOptions"
              :key="r.value"
              class="region-checkbox"
              :class="{ checked: selectedRegions.includes(r.value) }"
            >
              <input
                type="checkbox"
                :value="r.value"
                :checked="selectedRegions.includes(r.value)"
                @change="toggleRegion(r.value)"
              />
              <span class="region-label">{{ r.label }}</span>
              <span class="region-value">{{ r.value }}</span>
            </label>
          </div>
          <div class="selected-summary">
            已选择 {{ selectedRegions.length }} 个 Region
          </div>
        </div>

        <div class="button-group">
          <button class="btn btn-primary" :disabled="loading" @click="queryECS">
            {{ loading ? '查询中...' : '查询 ECS' }}
          </button>
          <button class="btn btn-secondary" @click="clearInputs">清空输入</button>
          <button class="btn btn-secondary" @click="clearResults">清空结果</button>
        </div>
      </section>

      <!-- Error / Success Messages -->
      <div v-if="errorMessage" class="message error-message">
        {{ errorMessage }}
      </div>
      <div v-if="successMessage" class="message success-message">
        {{ successMessage }}
      </div>

      <!-- Results Section -->
      <section class="results-section">
        <h2>ECS 实例列表</h2>
        <div v-if="regionResults.length === 0 && !loading" class="empty-state">
          暂无数据，请输入凭据并选择 Region 后点击"查询 ECS"
        </div>

        <div v-if="regionResults.length > 0">
          <div class="total-summary">
            共 {{ totalInstances }} 个实例，跨 {{ regionResults.length }} 个 Region
          </div>

          <div v-for="region in regionResults" :key="region.region" class="region-group">
            <div class="region-header">
              <span class="region-name">{{ getRegionLabel(region.region) }}</span>
              <span class="region-id">({{ region.region }})</span>
              <span v-if="region.error" class="region-error">{{ region.error }}</span>
              <span v-else class="region-count">{{ region.instances.length }} 个实例</span>
            </div>

            <div v-if="region.error && region.instances.length === 0" class="region-error-msg">
              查询失败: {{ region.error }}
            </div>

            <div v-if="region.instances.length > 0" class="table-wrapper">
              <table class="ecs-table">
                <thead>
                  <tr>
                    <th>InstanceId</th>
                    <th>InstanceName</th>
                    <th>Status</th>
                    <th>ZoneId</th>
                    <th>PublicIp</th>
                    <th>PrivateIp</th>
                    <th>CreationTime</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="inst in region.instances" :key="inst.instanceId">
                    <td>{{ inst.instanceId }}</td>
                    <td>{{ inst.instanceName }}</td>
                    <td>
                      <span class="status-badge" :class="getStatusClass(inst.status)">
                        {{ inst.status }}
                      </span>
                    </td>
                    <td>{{ inst.zoneId }}</td>
                    <td>{{ inst.publicIp || '-' }}</td>
                    <td>{{ inst.privateIp || '-' }}</td>
                    <td>{{ inst.creationTime }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
  color: #333;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-header h1 {
  font-size: 20px;
  font-weight: 600;
}

.version {
  font-size: 12px;
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 8px;
  border-radius: 10px;
}

.app-main {
  flex: 1;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

/* Input Section */
.input-section {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-section h2 {
  font-size: 16px;
  margin-bottom: 16px;
  color: #444;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #555;
}

.form-group input,
.form-group select {
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

/* Region Multi-Select */
.region-section {
  margin-bottom: 16px;
}

.region-section > label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #555;
  margin-bottom: 8px;
}

.region-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 8px;
  margin-bottom: 8px;
}

.region-checkbox {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
}

.region-checkbox:hover {
  border-color: #667eea;
}

.region-checkbox.checked {
  background: #f0f2ff;
  border-color: #667eea;
}

.region-checkbox input[type="checkbox"] {
  margin: 0;
  accent-color: #667eea;
}

.region-label {
  color: #333;
  font-weight: 500;
}

.region-value {
  color: #999;
  font-size: 11px;
}

.selected-summary {
  font-size: 13px;
  color: #666;
  padding: 4px 0;
}

.button-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn {
  padding: 8px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #5a6fd6;
}

.btn-secondary {
  background: #e4e7ed;
  color: #606266;
}

.btn-secondary:hover {
  background: #d4d7dd;
}

/* Messages */
.message {
  padding: 12px 16px;
  border-radius: 4px;
  margin-bottom: 16px;
  font-size: 14px;
}

.error-message {
  background: #fef0f0;
  color: #f56c6c;
  border: 1px solid #fde2e2;
}

.success-message {
  background: #f0f9eb;
  color: #67c23a;
  border: 1px solid #e1f3d8;
}

/* Results Section */
.results-section {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.results-section h2 {
  font-size: 16px;
  margin-bottom: 16px;
  color: #444;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
  font-size: 14px;
}

.total-summary {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

/* Region Group */
.region-group {
  margin-bottom: 24px;
}

.region-group:last-child {
  margin-bottom: 0;
}

.region-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: #f5f7fa;
  border-radius: 6px;
  margin-bottom: 8px;
}

.region-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.region-id {
  font-size: 13px;
  color: #999;
}

.region-count {
  margin-left: auto;
  font-size: 13px;
  color: #67c23a;
  font-weight: 500;
}

.region-error {
  margin-left: auto;
  font-size: 13px;
  color: #f56c6c;
  font-weight: 500;
}

.region-error-msg {
  padding: 12px;
  background: #fef0f0;
  color: #f56c6c;
  border-radius: 4px;
  font-size: 13px;
  margin-bottom: 8px;
}

.table-wrapper {
  overflow-x: auto;
}

.ecs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.ecs-table th {
  background: #fafafa;
  padding: 8px 12px;
  text-align: left;
  font-weight: 600;
  color: #606266;
  border-bottom: 2px solid #ebeef5;
  white-space: nowrap;
}

.ecs-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #ebeef5;
  color: #606266;
}

.ecs-table tr:hover td {
  background: #f5f7fa;
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
}

.status-running {
  background: #f0f9eb;
  color: #67c23a;
}

.status-stopped {
  background: #fef0f0;
  color: #f56c6c;
}

.status-starting {
  background: #fdf6ec;
  color: #e6a23c;
}

.status-stopping {
  background: #fdf6ec;
  color: #e6a23c;
}

.status-default {
  background: #f4f4f5;
  color: #909399;
}
</style>
