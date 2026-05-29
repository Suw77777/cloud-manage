<script setup>
import { ref } from 'vue'

const props = defineProps({
  region: String,
  vpcs: Array,
  vpcDetail: Object,
  vswitches: Array,
  loading: Boolean,
  loadingDetail: Boolean,
  regionOptions: Array
})

const emit = defineEmits([
  'update:region',
  'fetch-vpcs',
  'fetch-detail',
  'fetch-vswitches',
  'clear'
])

const selectedVpc = ref(null)
const showDetail = ref(false)
const showVSwitches = ref(false)

function viewDetail(vpc) {
  selectedVpc.value = vpc
  showDetail.value = true
  showVSwitches.value = false
  emit('fetch-detail', vpc.vpcId)
}

function viewVSwitches(vpc) {
  selectedVpc.value = vpc
  showVSwitches.value = true
  showDetail.value = false
  emit('fetch-vswitches', vpc.vpcId)
}

function closeDetail() {
  showDetail.value = false
  showVSwitches.value = false
  selectedVpc.value = null
}
</script>

<template>
  <div class="vpc-manager">
    <!-- Region Selection -->
    <section class="input-section">
      <h2>VPC 网络管理</h2>
      <div class="form-grid">
        <div class="form-group">
          <label>Region</label>
          <select :value="region" @change="emit('update:region', $event.target.value)">
            <option v-for="r in regionOptions" :key="r.value" :value="r.value">
              {{ r.label }}
            </option>
          </select>
        </div>
      </div>
      <div class="button-group">
        <button class="btn btn-primary" :disabled="loading" @click="emit('fetch-vpcs')">
          {{ loading ? '查询中...' : '查询 VPC' }}
        </button>
        <button class="btn btn-secondary" @click="emit('clear')">清空</button>
      </div>
    </section>

    <!-- VPC List -->
    <section v-if="vpcs.length > 0" class="result-section">
      <h3>VPC 列表 ({{ vpcs.length }})</h3>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>VPC ID</th>
              <th>名称</th>
              <th>CIDR</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="vpc in vpcs" :key="vpc.vpcId">
              <td>{{ vpc.vpcId }}</td>
              <td>{{ vpc.vpcName || '-' }}</td>
              <td>{{ vpc.cidrBlock }}</td>
              <td>
                <span class="status-badge" :class="vpc.status === 'Available' ? 'status-ok' : 'status-warn'">
                  {{ vpc.status }}
                </span>
              </td>
              <td>
                <button class="btn btn-sm" @click="viewDetail(vpc)">详情</button>
                <button class="btn btn-sm" @click="viewVSwitches(vpc)">VSwitch</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- VPC Detail -->
    <section v-if="showDetail && vpcDetail" class="detail-section">
      <h3>VPC 详情</h3>
      <div class="detail-grid">
        <div class="detail-item">
          <span class="label">ID:</span>
          <span class="value">{{ vpcDetail.vpcId }}</span>
        </div>
        <div class="detail-item">
          <span class="label">名称:</span>
          <span class="value">{{ vpcDetail.vpcName || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">CIDR:</span>
          <span class="value">{{ vpcDetail.cidrBlock }}</span>
        </div>
        <div class="detail-item">
          <span class="label">状态:</span>
          <span class="value">{{ vpcDetail.status }}</span>
        </div>
        <div class="detail-item">
          <span class="label">区域:</span>
          <span class="value">{{ vpcDetail.regionId }}</span>
        </div>
        <div class="detail-item">
          <span class="label">描述:</span>
          <span class="value">{{ vpcDetail.description || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">创建时间:</span>
          <span class="value">{{ vpcDetail.creationTime }}</span>
        </div>
        <div class="detail-item full-width">
          <span class="label">VSwitch IDs:</span>
          <span class="value">{{ vpcDetail.vswitchIds?.join(', ') || '-' }}</span>
        </div>
      </div>
      <button class="btn btn-secondary" @click="closeDetail">关闭</button>
    </section>

    <!-- VSwitch List -->
    <section v-if="showVSwitches && vswitches.length > 0" class="result-section">
      <h3>VSwitch 列表 ({{ vswitches.length }})</h3>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>VSwitch ID</th>
              <th>名称</th>
              <th>CIDR</th>
              <th>可用区</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="vs in vswitches" :key="vs.vswitchId">
              <td>{{ vs.vswitchId }}</td>
              <td>{{ vs.vswitchName || '-' }}</td>
              <td>{{ vs.cidrBlock }}</td>
              <td>{{ vs.zoneId }}</td>
              <td>
                <span class="status-badge" :class="vs.status === 'Available' ? 'status-ok' : 'status-warn'">
                  {{ vs.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <button class="btn btn-secondary" @click="closeDetail">关闭</button>
    </section>

    <!-- Empty State -->
    <section v-if="!loading && vpcs.length === 0" class="empty-state">
      点击"查询 VPC"获取 VPC 列表
    </section>
  </div>
</template>
