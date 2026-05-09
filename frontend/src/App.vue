<script setup>
import { ref } from 'vue'
import { QueryECS } from '../wailsjs/go/main/App'

const accessKeyId = ref('')
const accessKeySecret = ref('')
const region = ref('cn-hangzhou')
const env = ref('dev')
const instances = ref([])
const errorMessage = ref('')
const successMessage = ref('')
const loading = ref(false)

const envOptions = [
  { label: '开发环境 (dev)', value: 'dev' },
  { label: '预发布环境 (pre)', value: 'pre' },
  { label: '生产环境 (prod)', value: 'prod' },
]

async function queryECS() {
  errorMessage.value = ''
  successMessage.value = ''
  instances.value = []

  if (!accessKeyId.value.trim()) {
    errorMessage.value = '请输入 AccessKey ID'
    return
  }
  if (!accessKeySecret.value.trim()) {
    errorMessage.value = '请输入 AccessKey Secret'
    return
  }
  if (!region.value.trim()) {
    errorMessage.value = '请输入 Region'
    return
  }

  loading.value = true
  try {
    const result = await QueryECS(
      accessKeyId.value.trim(),
      accessKeySecret.value.trim(),
      region.value.trim(),
      env.value
    )
    if (result.success) {
      instances.value = result.instances || []
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
  region.value = 'cn-hangzhou'
  env.value = 'dev'
  errorMessage.value = ''
  successMessage.value = ''
}

function clearResults() {
  instances.value = []
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
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>Cloud 管理小助手</h1>
      <span class="version">v0.0.1</span>
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
            <label for="region">Region</label>
            <input
              id="region"
              v-model="region"
              type="text"
              placeholder="例如: cn-hangzhou"
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
        <div v-if="instances.length === 0 && !loading" class="empty-state">
          暂无数据，请输入凭据后点击"查询 ECS"
        </div>
        <div v-if="instances.length > 0" class="table-wrapper">
          <table class="ecs-table">
            <thead>
              <tr>
                <th>InstanceId</th>
                <th>InstanceName</th>
                <th>Status</th>
                <th>RegionId</th>
                <th>ZoneId</th>
                <th>PublicIp</th>
                <th>PrivateIp</th>
                <th>CreationTime</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inst in instances" :key="inst.instanceId">
                <td>{{ inst.instanceId }}</td>
                <td>{{ inst.instanceName }}</td>
                <td>
                  <span class="status-badge" :class="getStatusClass(inst.status)">
                    {{ inst.status }}
                  </span>
                </td>
                <td>{{ inst.regionId }}</td>
                <td>{{ inst.zoneId }}</td>
                <td>{{ inst.publicIp || '-' }}</td>
                <td>{{ inst.privateIp || '-' }}</td>
                <td>{{ inst.creationTime }}</td>
              </tr>
            </tbody>
          </table>
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

.table-wrapper {
  overflow-x: auto;
}

.ecs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.ecs-table th {
  background: #f5f7fa;
  padding: 10px 12px;
  text-align: left;
  font-weight: 600;
  color: #606266;
  border-bottom: 2px solid #ebeef5;
  white-space: nowrap;
}

.ecs-table td {
  padding: 10px 12px;
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
