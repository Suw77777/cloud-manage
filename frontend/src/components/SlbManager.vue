<script setup>
import { ref } from 'vue'

const props = defineProps({
  region: String,
  slbs: Array,
  slbDetail: Object,
  listeners: Array,
  loading: Boolean,
  loadingDetail: Boolean,
  regionOptions: Array
})

const emit = defineEmits([
  'update:region',
  'fetch-slbs',
  'fetch-detail',
  'fetch-listeners',
  'clear'
])

const selectedSlb = ref(null)
const showDetail = ref(false)
const showListeners = ref(false)

function viewDetail(slb) {
  selectedSlb.value = slb
  showDetail.value = true
  showListeners.value = false
  emit('fetch-detail', slb.loadBalancerId)
}

function viewListeners(slb) {
  selectedSlb.value = slb
  showListeners.value = true
  showDetail.value = false
  emit('fetch-listeners', slb.loadBalancerId)
}

function closeDetail() {
  showDetail.value = false
  showListeners.value = false
  selectedSlb.value = null
}
</script>

<template>
  <div class="slb-manager">
    <!-- Region Selection -->
    <section class="input-section">
      <h2>负载均衡管理</h2>
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
        <button class="btn btn-primary" :disabled="loading" @click="emit('fetch-slbs')">
          {{ loading ? '查询中...' : '查询 SLB' }}
        </button>
        <button class="btn btn-secondary" @click="emit('clear')">清空</button>
      </div>
    </section>

    <!-- SLB List -->
    <section v-if="slbs.length > 0" class="result-section">
      <h3>SLB 列表 ({{ slbs.length }})</h3>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>SLB ID</th>
              <th>名称</th>
              <th>地址</th>
              <th>类型</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="slb in slbs" :key="slb.loadBalancerId">
              <td>{{ slb.loadBalancerId }}</td>
              <td>{{ slb.loadBalancerName || '-' }}</td>
              <td>{{ slb.address }}</td>
              <td>{{ slb.addressType }}</td>
              <td>
                <span class="status-badge" :class="slb.status === 'active' ? 'status-ok' : 'status-warn'">
                  {{ slb.status }}
                </span>
              </td>
              <td>
                <button class="btn btn-sm" @click="viewDetail(slb)">详情</button>
                <button class="btn btn-sm" @click="viewListeners(slb)">监听器</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- SLB Detail -->
    <section v-if="showDetail && slbDetail" class="detail-section">
      <h3>SLB 详情</h3>
      <div class="detail-grid">
        <div class="detail-item">
          <span class="label">ID:</span>
          <span class="value">{{ slbDetail.loadBalancerId }}</span>
        </div>
        <div class="detail-item">
          <span class="label">名称:</span>
          <span class="value">{{ slbDetail.loadBalancerName || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">地址:</span>
          <span class="value">{{ slbDetail.address }}</span>
        </div>
        <div class="detail-item">
          <span class="label">类型:</span>
          <span class="value">{{ slbDetail.addressType }}</span>
        </div>
        <div class="detail-item">
          <span class="label">状态:</span>
          <span class="value">{{ slbDetail.status }}</span>
        </div>
        <div class="detail-item">
          <span class="label">区域:</span>
          <span class="value">{{ slbDetail.regionId }}</span>
        </div>
        <div class="detail-item">
          <span class="label">VPC ID:</span>
          <span class="value">{{ slbDetail.vpcId || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">VSwitch ID:</span>
          <span class="value">{{ slbDetail.vswitchId || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">监听器数:</span>
          <span class="value">{{ slbDetail.listenerCount }}</span>
        </div>
        <div class="detail-item">
          <span class="label">带宽:</span>
          <span class="value">{{ slbDetail.bandwidth }} Mbps</span>
        </div>
        <div class="detail-item">
          <span class="label">创建时间:</span>
          <span class="value">{{ slbDetail.creationTime }}</span>
        </div>
      </div>
      <button class="btn btn-secondary" @click="closeDetail">关闭</button>
    </section>

    <!-- Listener List -->
    <section v-if="showListeners && listeners.length > 0" class="result-section">
      <h3>监听器列表 ({{ listeners.length }})</h3>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>端口</th>
              <th>协议</th>
              <th>状态</th>
              <th>带宽</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(l, idx) in listeners" :key="idx">
              <td>{{ l.listenerPort }}</td>
              <td>{{ l.listenerProtocol }}</td>
              <td>
                <span class="status-badge" :class="l.status === 'running' ? 'status-ok' : 'status-warn'">
                  {{ l.status }}
                </span>
              </td>
              <td>{{ l.bandwidth }} Mbps</td>
            </tr>
          </tbody>
        </table>
      </div>
      <button class="btn btn-secondary" @click="closeDetail">关闭</button>
    </section>

    <!-- Empty State -->
    <section v-if="!loading && slbs.length === 0" class="empty-state">
      点击"查询 SLB"获取负载均衡列表
    </section>
  </div>
</template>
