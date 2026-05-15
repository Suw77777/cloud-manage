<script setup>
defineProps({
  regionResults: { type: Array, default: () => [] },
  totalInstances: { type: Number, default: 0 },
  regionOptions: { type: Array, default: () => [] },
})

const emit = defineEmits(['viewDetail', 'requestAction'])

function getStatusClass(status) {
  const s = (status || '').toLowerCase()
  if (s === 'running') return 'status-running'
  if (s === 'stopped') return 'status-stopped'
  if (s === 'starting') return 'status-starting'
  if (s === 'stopping') return 'status-stopping'
  return 'status-default'
}

function getRegionLabel(value, options) {
  const found = options.find(r => r.value === value)
  return found ? found.label : value
}
</script>

<template>
  <section class="results-section">
    <h2>ECS 实例列表</h2>
    <div v-if="regionResults.length === 0" class="empty-state">
      暂无数据，请输入凭据并选择 Region 后点击"查询 ECS"
    </div>

    <div v-if="regionResults.length > 0">
      <div class="total-summary">
        共 {{ totalInstances }} 个实例，跨 {{ regionResults.length }} 个 Region
      </div>

      <div v-for="region in regionResults" :key="region.region" class="region-group">
        <div class="region-header">
          <span class="region-name">{{ getRegionLabel(region.region, regionOptions) }}</span>
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
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inst in region.instances" :key="inst.instanceId">
                <td>
                  <a class="link" @click="emit('viewDetail', inst.instanceId, region.region)">
                    {{ inst.instanceId }}
                  </a>
                </td>
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
                <td>
                  <div class="action-btns">
                    <button
                      v-if="inst.status === 'Stopped'"
                      class="btn btn-sm btn-success"
                      @click="emit('requestAction', 'start', inst.instanceId, inst.instanceName, region.region)"
                    >
                      启动
                    </button>
                    <button
                      v-if="inst.status === 'Running'"
                      class="btn btn-sm btn-warning"
                      @click="emit('requestAction', 'stop', inst.instanceId, inst.instanceName, region.region)"
                    >
                      停止
                    </button>
                    <button
                      v-if="inst.status === 'Running'"
                      class="btn btn-sm btn-info"
                      @click="emit('requestAction', 'reboot', inst.instanceId, inst.instanceName, region.region)"
                    >
                      重启
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>
</template>
