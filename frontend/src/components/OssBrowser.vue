<script setup>
const props = defineProps({
  region: String,
  currentBucket: String,
  currentPrefix: String,
  buckets: Array,
  objects: Array,
  isTruncated: Boolean,
  breadcrumb: Array,
  loading: Boolean,
  loadingBuckets: Boolean,
  regionOptions: Array,
})

const emit = defineEmits([
  'update:region',
  'fetch-buckets',
  'fetch-objects',
  'navigate-folder',
  'navigate-root',
  'navigate-up',
  'select-bucket',
  'back-to-buckets',
  'clear'
])

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '-'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const date = new Date(dateStr)
    return date.toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}

const getObjectIcon = (obj) => {
  if (obj.isFolder) return '📁'
  const ext = obj.key?.split('.').pop()?.toLowerCase()
  const icons = {
    'log': '📄',
    'txt': '📄',
    'json': '📋',
    'csv': '📊',
    'gz': '🗜️',
    'zip': '🗜️',
  }
  return icons[ext] || '📄'
}

const getObjectName = (key) => {
  if (!key) return ''
  const parts = key.split('/')
  return parts[parts.length - 1] || parts[parts.length - 2] || key
}
</script>

<template>
  <section class="oss-browser">
    <h2>OSS 对象存储</h2>

    <!-- Region Selection -->
    <div class="oss-config">
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
          <label>&nbsp;</label>
          <button
            class="btn btn-primary"
            :disabled="loadingBuckets"
            @click="emit('fetch-buckets')"
          >
            {{ loadingBuckets ? '加载中...' : '获取 Bucket 列表' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Bucket List -->
    <div v-if="!currentBucket && buckets && buckets.length > 0" class="bucket-list">
      <h3>Bucket 列表</h3>
      <div class="bucket-grid">
        <div
          v-for="bucket in buckets"
          :key="bucket.name"
          class="bucket-card"
          @click="emit('select-bucket', bucket.name)"
        >
          <div class="bucket-icon">🪣</div>
          <div class="bucket-info">
            <div class="bucket-name">{{ bucket.name }}</div>
            <div class="bucket-meta">
              <span>{{ bucket.location || bucket.region }}</span>
              <span>{{ bucket.storageClass }}</span>
            </div>
            <div class="bucket-date">创建于 {{ formatDate(bucket.creationDate) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Object Browser -->
    <div v-if="currentBucket" class="object-browser">
      <!-- Breadcrumb -->
      <div class="breadcrumb">
        <button class="breadcrumb-btn" @click="emit('back-to-buckets')">⬅ Bucket 列表</button>
        <span class="breadcrumb-sep">/</span>
        <template v-for="(item, index) in breadcrumb" :key="index">
          <button
            class="breadcrumb-btn"
            :class="{ active: index === breadcrumb.length - 1 }"
            @click="emit('navigate-folder', item.prefix)"
          >
            {{ item.name }}
          </button>
          <span v-if="index < breadcrumb.length - 1" class="breadcrumb-sep">/</span>
        </template>
      </div>

      <!-- Object Table -->
      <div v-if="loading" class="loading-state">加载中...</div>

      <div v-else-if="objects && objects.length > 0" class="object-table-wrapper">
        <table class="object-table">
          <thead>
            <tr>
              <th class="col-icon"></th>
              <th class="col-name">名称</th>
              <th class="col-size">大小</th>
              <th class="col-modified">修改时间</th>
              <th class="col-class">存储类型</th>
            </tr>
          </thead>
          <tbody>
            <!-- Parent Directory -->
            <tr v-if="currentPrefix" class="folder-row" @click="emit('navigate-up')">
              <td class="col-icon">📁</td>
              <td class="col-name">..</td>
              <td class="col-size">-</td>
              <td class="col-modified">-</td>
              <td class="col-class">-</td>
            </tr>
            <!-- Objects -->
            <tr
              v-for="obj in objects"
              :key="obj.key"
              :class="{ 'folder-row': obj.isFolder }"
              @click="obj.isFolder ? emit('navigate-folder', obj.key) : null"
            >
              <td class="col-icon">{{ getObjectIcon(obj) }}</td>
              <td class="col-name">{{ getObjectName(obj.key) }}</td>
              <td class="col-size">{{ obj.isFolder ? '-' : formatFileSize(obj.size) }}</td>
              <td class="col-modified">{{ formatDate(obj.lastModified) }}</td>
              <td class="col-class">{{ obj.storageClass || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else class="empty-state">
        此目录为空
      </div>

      <div v-if="isTruncated" class="truncated-notice">
        结果已截断，仅显示前 100 个对象
      </div>
    </div>

    <div v-else-if="!loadingBuckets && (!buckets || buckets.length === 0)" class="empty-state">
      点击"获取 Bucket 列表"开始浏览
    </div>
  </section>
</template>

<style scoped>
.oss-browser {
  padding: 0;
}

.oss-config {
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
  align-items: end;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
}

.form-group select {
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.bucket-list,
.object-browser {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.bucket-list h3,
.object-browser h3 {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #444;
}

.bucket-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.bucket-card {
  display: flex;
  gap: 12px;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.bucket-card:hover {
  border-color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.bucket-icon {
  font-size: 32px;
}

.bucket-info {
  flex: 1;
  min-width: 0;
}

.bucket-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
  word-break: break-all;
}

.bucket-meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.bucket-date {
  font-size: 11px;
  color: #c0c4cc;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.breadcrumb-btn {
  background: none;
  border: none;
  color: #667eea;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}

.breadcrumb-btn:hover {
  background: #f0f2ff;
}

.breadcrumb-btn.active {
  color: #333;
  font-weight: 600;
  cursor: default;
}

.breadcrumb-sep {
  color: #c0c4cc;
  font-size: 12px;
}

.object-table-wrapper {
  overflow-x: auto;
  max-height: 500px;
  overflow-y: auto;
}

.object-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.object-table th {
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

.object-table td {
  padding: 10px 12px;
  border-bottom: 1px solid #ebeef5;
  color: #333;
}

.object-table tr:hover td {
  background: #f5f7fa;
}

.folder-row {
  cursor: pointer;
}

.folder-row td {
  color: #667eea;
  font-weight: 500;
}

.col-icon {
  width: 40px;
  text-align: center;
}

.col-name {
  min-width: 200px;
}

.col-size {
  width: 100px;
}

.col-modified {
  width: 180px;
}

.col-class {
  width: 100px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
  font-size: 14px;
}

.truncated-notice {
  text-align: center;
  padding: 12px;
  color: #e6a23c;
  font-size: 13px;
  font-weight: 500;
}
</style>
