<script setup>
import { useECS } from './composables/useECS'
import EcsResultTable from './components/EcsResultTable.vue'
import InstanceDetailModal from './components/InstanceDetailModal.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import OperationLog from './components/OperationLog.vue'

const ecs = useECS()

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

function toggleRegion(regionValue) {
  const idx = ecs.selectedRegions.value.indexOf(regionValue)
  if (idx === -1) {
    ecs.selectedRegions.value.push(regionValue)
  } else {
    ecs.selectedRegions.value.splice(idx, 1)
  }
}
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>Cloud 管理小助手</h1>
      <span class="version">v0.0.3</span>
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
              v-model="ecs.accessKeyId.value"
              type="text"
              placeholder="请输入 AccessKey ID"
              autocomplete="off"
            />
          </div>
          <div class="form-group">
            <label for="accessKeySecret">AccessKey Secret</label>
            <input
              id="accessKeySecret"
              v-model="ecs.accessKeySecret.value"
              type="password"
              placeholder="请输入 AccessKey Secret"
              autocomplete="off"
            />
          </div>
          <div class="form-group">
            <label for="env">环境</label>
            <select id="env" v-model="ecs.env.value">
              <option v-for="opt in envOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>

        <div class="region-section">
          <label>Region（可多选）</label>
          <div class="region-grid">
            <label
              v-for="r in regionOptions"
              :key="r.value"
              class="region-checkbox"
              :class="{ checked: ecs.selectedRegions.value.includes(r.value) }"
            >
              <input
                type="checkbox"
                :value="r.value"
                :checked="ecs.selectedRegions.value.includes(r.value)"
                @change="toggleRegion(r.value)"
              />
              <span class="region-label">{{ r.label }}</span>
              <span class="region-value">{{ r.value }}</span>
            </label>
          </div>
          <div class="selected-summary">
            已选择 {{ ecs.selectedRegions.value.length }} 个 Region
          </div>
        </div>

        <div class="button-group">
          <button class="btn btn-primary" :disabled="ecs.loading.value" @click="ecs.queryECS">
            {{ ecs.loading.value ? '查询中...' : '查询 ECS' }}
          </button>
          <button class="btn btn-secondary" @click="ecs.clearInputs">清空输入</button>
          <button class="btn btn-secondary" @click="ecs.clearResults">清空结果</button>
        </div>
      </section>

      <!-- Error / Success Messages -->
      <div v-if="ecs.errorMessage.value" class="message error-message">
        {{ ecs.errorMessage.value }}
      </div>
      <div v-if="ecs.successMessage.value" class="message success-message">
        {{ ecs.successMessage.value }}
      </div>

      <!-- Results -->
      <EcsResultTable
        :region-results="ecs.regionResults.value"
        :total-instances="ecs.totalInstances.value"
        :region-options="regionOptions"
        @view-detail="ecs.viewDetail"
        @request-action="ecs.requestAction"
      />

      <!-- Operation Log -->
      <OperationLog
        :logs="ecs.operationLogs.value"
        @clear="ecs.clearLogs"
      />
    </main>

    <!-- Detail Modal -->
    <InstanceDetailModal
      v-if="ecs.showDetail.value"
      :loading="ecs.detailLoading.value"
      :data="ecs.detailData.value"
      @close="ecs.closeDetail"
    />

    <!-- Confirm Dialog -->
    <ConfirmDialog
      v-if="ecs.showConfirm.value"
      :action="ecs.confirmAction.value"
      :instance-id="ecs.confirmInstanceId.value"
      :instance-name="ecs.confirmInstanceName.value"
      :region="ecs.confirmRegion.value"
      :is-prod="ecs.env.value === 'prod'"
      @confirm="ecs.executeConfirm"
      @cancel="ecs.cancelConfirm"
    />
  </div>
</template>

<style>
@import './assets/main.css';
</style>
