<script setup>
import { ref, onMounted } from 'vue'
import { useECS } from './composables/useECS'
import { useCMS } from './composables/useCMS'
import { useSLS } from './composables/useSLS'
import { useOSS } from './composables/useOSS'
import { useVPC } from './composables/useVPC'
import { useSLB } from './composables/useSLB'
import { useConfig } from './composables/useConfig'
import EcsResultTable from './components/EcsResultTable.vue'
import InstanceDetailModal from './components/InstanceDetailModal.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import OperationLog from './components/OperationLog.vue'
import CmsMonitor from './components/CmsMonitor.vue'
import SlsQuery from './components/SlsQuery.vue'
import OssBrowser from './components/OssBrowser.vue'
import VpcManager from './components/VpcManager.vue'
import SlbManager from './components/SlbManager.vue'

const ecs = useECS()
const cms = useCMS(ecs.accessKeyId, ecs.accessKeySecret)
const sls = useSLS(ecs.accessKeyId, ecs.accessKeySecret)
const oss = useOSS(ecs.accessKeyId, ecs.accessKeySecret)
const vpc = useVPC(ecs.accessKeyId, ecs.accessKeySecret)
const slb = useSLB(ecs.accessKeyId, ecs.accessKeySecret)
const cfg = useConfig()

const activeTab = ref('ecs')
const theme = ref('auto')

function applyTheme(t) {
  if (t === 'dark' || (t === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}

function toggleTheme() {
  const themes = ['auto', 'light', 'dark']
  const idx = themes.indexOf(theme.value)
  theme.value = themes[(idx + 1) % themes.length]
  applyTheme(theme.value)
}

async function handleProfileSwitch(name) {
  await cfg.switchProfile(name)
  const creds = await cfg.getProfileCredentials(name)
  if (creds) {
    ecs.accessKeyId.value = creds.accessKeyId || ''
    ecs.accessKeySecret.value = creds.accessKeySecret || ''
    if (creds.region) {
      ecs.region.value = creds.region
    }
  }
}

onMounted(async () => {
  applyTheme(theme.value)
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'auto') applyTheme('auto')
  })
  await cfg.loadProfiles()
  // Auto-fill credentials from config if no env vars
  if (cfg.currentProfile.value) {
    const creds = await cfg.getProfileCredentials(cfg.currentProfile.value)
    if (creds && !ecs.accessKeyId.value) {
      ecs.accessKeyId.value = creds.accessKeyId || ''
      ecs.accessKeySecret.value = creds.accessKeySecret || ''
      if (creds.region) {
        ecs.region.value = creds.region
      }
    }
  }
})

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
      <span class="version">v0.2.0</span>
      <div v-if="cfg.profiles.value.length > 0" class="profile-selector">
        <select v-model="cfg.currentProfile.value" @change="handleProfileSwitch($event.target.value)" class="profile-select">
          <option v-for="p in cfg.profiles.value" :key="p" :value="p">{{ p }}</option>
        </select>
      </div>
      <button class="theme-toggle" @click="toggleTheme" :title="theme === 'auto' ? '跟随系统' : theme === 'dark' ? '暗色' : '浅色'">
        {{ theme === 'auto' ? '🖥️' : theme === 'dark' ? '🌙' : '☀️' }}
      </button>
    </header>

    <!-- Tab Navigation -->
    <nav class="tab-nav">
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'ecs' }"
        @click="activeTab = 'ecs'"
      >
        ECS 实例管理
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'cms' }"
        @click="activeTab = 'cms'"
      >
        云监控
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'sls' }"
        @click="activeTab = 'sls'"
      >
        SLS 日志
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'oss' }"
        @click="activeTab = 'oss'"
      >
        OSS 存储
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'vpc' }"
        @click="activeTab = 'vpc'"
      >
        VPC 网络
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'slb' }"
        @click="activeTab = 'slb'"
      >
        负载均衡
      </button>
    </nav>

    <main class="app-main">
      <!-- ECS Tab -->
      <div v-show="activeTab === 'ecs'">
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
      </div>

      <!-- CMS Tab -->
      <div v-show="activeTab === 'cms'">
        <!-- Connection Info Display -->
        <section class="input-section">
          <h2>监控配置</h2>
          <div class="form-grid">
            <div class="form-group">
              <label>AccessKey ID</label>
              <input
                :value="ecs.accessKeyId.value"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
            <div class="form-group">
              <label>AccessKey Secret</label>
              <input
                :value="ecs.accessKeySecret.value ? '********' : ''"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
          </div>
        </section>

        <!-- Error / Success Messages -->
        <div v-if="cms.metricsError.value" class="message error-message">
          {{ cms.metricsError.value }}
        </div>
        <div v-if="cms.metricsSuccess.value" class="message success-message">
          {{ cms.metricsSuccess.value }}
        </div>

        <!-- CMS Monitor Component -->
        <CmsMonitor
          :loading="cms.metricsLoading.value"
          :metrics="cms.metricsData.value"
          :selected-instances="cms.selectedInstances.value"
          :region-results="ecs.regionResults.value"
          :region-options="regionOptions"
          @query="cms.queryMetrics"
          @toggle-instance="cms.toggleInstance"
          @select-all="cms.selectAllInstances"
          @clear-selection="cms.clearSelection"
        />
      </div>

      <!-- SLS Tab -->
      <div v-show="activeTab === 'sls'">
        <!-- Connection Info Display -->
        <section class="input-section">
          <h2>SLS 配置</h2>
          <div class="form-grid">
            <div class="form-group">
              <label>AccessKey ID</label>
              <input
                :value="ecs.accessKeyId.value"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
            <div class="form-group">
              <label>AccessKey Secret</label>
              <input
                :value="ecs.accessKeySecret.value ? '********' : ''"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
          </div>
        </section>

        <!-- Error / Success Messages -->
        <div v-if="sls.error.value" class="message error-message">
          {{ sls.error.value }}
        </div>
        <div v-if="sls.success.value" class="message success-message">
          {{ sls.success.value }}
        </div>

        <!-- SLS Query Component -->
        <SlsQuery
          :region="sls.region.value"
          :project="sls.project.value"
          :logstore="sls.logstore.value"
          :query="sls.query.value"
          :time-range="sls.timeRange.value"
          :max-lines="sls.maxLines.value"
          :log-stores="sls.logStores.value"
          :logs="sls.logs.value"
          :log-count="sls.logCount.value"
          :has-more="sls.hasMore.value"
          :loading="sls.loading.value"
          :loading-stores="sls.loadingStores.value"
          :streaming="sls.streaming.value"
          :stream-progress="sls.streamProgress.value"
          :region-options="regionOptions"
          :page-size="sls.pageSize.value"
          :current-page="sls.currentPage.value"
          :total-logs="sls.totalLogs.value"
          @update:region="sls.region.value = $event"
          @update:project="sls.project.value = $event"
          @update:logstore="sls.logstore.value = $event"
          @update:query="sls.query.value = $event"
          @update:timeRange="sls.timeRange.value = $event"
          @update:maxLines="sls.maxLines.value = $event"
          @fetch-stores="sls.fetchLogStores"
          @query-logs="sls.queryLogs"
          @query-logs-stream="sls.queryLogsStream"
          @page-change="sls.setPage"
          @page-size-change="sls.setPageSize"
          @clear="sls.clearAll"
        />
      </div>

      <!-- OSS Tab -->
      <div v-show="activeTab === 'oss'">
        <!-- Connection Info Display -->
        <section class="input-section">
          <h2>OSS 配置</h2>
          <div class="form-grid">
            <div class="form-group">
              <label>AccessKey ID</label>
              <input
                :value="ecs.accessKeyId.value"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
            <div class="form-group">
              <label>AccessKey Secret</label>
              <input
                :value="ecs.accessKeySecret.value ? '********' : ''"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
          </div>
        </section>

        <!-- Error / Success Messages -->
        <div v-if="oss.error.value" class="message error-message">
          {{ oss.error.value }}
        </div>
        <div v-if="oss.success.value" class="message success-message">
          {{ oss.success.value }}
        </div>

        <!-- OSS Browser Component -->
        <OssBrowser
          :region="oss.region.value"
          :current-bucket="oss.currentBucket.value"
          :current-prefix="oss.currentPrefix.value"
          :buckets="oss.buckets.value"
          :objects="oss.objects.value"
          :is-truncated="oss.isTruncated.value"
          :breadcrumb="oss.breadcrumb.value"
          :loading="oss.loading.value"
          :loading-buckets="oss.loadingBuckets.value"
          :region-options="regionOptions"
          @update:region="oss.region.value = $event"
          @fetch-buckets="oss.fetchBuckets"
          @fetch-objects="oss.fetchObjects"
          @navigate-folder="oss.navigateToFolder"
          @navigate-root="oss.navigateToRoot"
          @navigate-up="oss.navigateUp"
          @select-bucket="oss.selectBucket"
          @back-to-buckets="oss.backToBuckets"
          @clear="oss.clearResults"
        />
      </div>

      <!-- VPC Tab -->
      <div v-show="activeTab === 'vpc'">
        <!-- Connection Info Display -->
        <section class="input-section">
          <h2>VPC 配置</h2>
          <div class="form-grid">
            <div class="form-group">
              <label>AccessKey ID</label>
              <input
                :value="ecs.accessKeyId.value"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
            <div class="form-group">
              <label>AccessKey Secret</label>
              <input
                :value="ecs.accessKeySecret.value ? '********' : ''"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
          </div>
        </section>

        <!-- Error / Success Messages -->
        <div v-if="vpc.error.value" class="message error-message">
          {{ vpc.error.value }}
        </div>
        <div v-if="vpc.success.value" class="message success-message">
          {{ vpc.success.value }}
        </div>

        <!-- VPC Manager Component -->
        <VpcManager
          :region="vpc.region.value"
          :vpcs="vpc.vpcs.value"
          :vpc-detail="vpc.vpcDetail.value"
          :vswitches="vpc.vswitches.value"
          :loading="vpc.loading.value"
          :loading-detail="vpc.loadingDetail.value"
          :region-options="regionOptions"
          @update:region="vpc.region.value = $event"
          @fetch-vpcs="vpc.fetchVPCs"
          @fetch-detail="vpc.fetchVPCDetail"
          @fetch-vswitches="vpc.fetchVSwitches"
          @clear="vpc.clearResults"
        />
      </div>

      <!-- SLB Tab -->
      <div v-show="activeTab === 'slb'">
        <!-- Connection Info Display -->
        <section class="input-section">
          <h2>SLB 配置</h2>
          <div class="form-grid">
            <div class="form-group">
              <label>AccessKey ID</label>
              <input
                :value="ecs.accessKeyId.value"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
            <div class="form-group">
              <label>AccessKey Secret</label>
              <input
                :value="ecs.accessKeySecret.value ? '********' : ''"
                type="text"
                placeholder="请在 ECS 页面输入"
                disabled
              />
            </div>
          </div>
        </section>

        <!-- Error / Success Messages -->
        <div v-if="slb.error.value" class="message error-message">
          {{ slb.error.value }}
        </div>
        <div v-if="slb.success.value" class="message success-message">
          {{ slb.success.value }}
        </div>

        <!-- SLB Manager Component -->
        <SlbManager
          :region="slb.region.value"
          :slbs="slb.slbs.value"
          :slb-detail="slb.slbDetail.value"
          :listeners="slb.listeners.value"
          :loading="slb.loading.value"
          :loading-detail="slb.loadingDetail.value"
          :region-options="regionOptions"
          @update:region="slb.region.value = $event"
          @fetch-slbs="slb.fetchSLBs"
          @fetch-detail="slb.fetchSLBDetail"
          @fetch-listeners="slb.fetchListeners"
          @clear="slb.clearResults"
        />
      </div>
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
