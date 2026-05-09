<script setup>
defineProps({
  loading: { type: Boolean, default: false },
  data: { type: Object, default: null },
})

const emit = defineEmits(['close'])

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
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <div class="modal-header">
        <h3>实例详情</h3>
        <button class="modal-close" @click="emit('close')">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="loading" class="empty-state">加载中...</div>
        <div v-else-if="data && data.error" class="region-error-msg">
          {{ data.error }}
        </div>
        <div v-else-if="data" class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">InstanceId</span>
            <span class="detail-value">{{ data.instanceId }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">InstanceName</span>
            <span class="detail-value">{{ data.instanceName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Description</span>
            <span class="detail-value">{{ data.description || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">HostName</span>
            <span class="detail-value">{{ data.hostName || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Status</span>
            <span class="detail-value">
              <span class="status-badge" :class="getStatusClass(data.status)">
                {{ data.status }}
              </span>
            </span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Region</span>
            <span class="detail-value">{{ data.regionId }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Zone</span>
            <span class="detail-value">{{ data.zoneId }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">InstanceType</span>
            <span class="detail-value">{{ data.instanceType }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">CPU / Memory</span>
            <span class="detail-value">{{ data.cpu }} vCPU / {{ data.memory }} MB</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">ImageId</span>
            <span class="detail-value">{{ data.imageId }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">InternetChargeType</span>
            <span class="detail-value">{{ data.internetChargeType }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">CreationTime</span>
            <span class="detail-value">{{ data.creationTime }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">ExpiredTime</span>
            <span class="detail-value">{{ data.expiredTime || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">StoppedMode</span>
            <span class="detail-value">{{ data.stoppedMode || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">PublicIp</span>
            <span class="detail-value">{{ (data.publicIp || []).join(', ') || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">PrivateIp</span>
            <span class="detail-value">{{ (data.privateIp || []).join(', ') || '-' }}</span>
          </div>
          <div class="detail-item full-width">
            <span class="detail-label">SecurityGroupIds</span>
            <span class="detail-value">{{ (data.securityGroupIds || []).join(', ') || '-' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
