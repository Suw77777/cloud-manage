<script setup>
defineProps({
  logs: { type: Array, default: () => [] },
})

const emit = defineEmits(['clear'])

const actionLabels = { start: '启动', stop: '停止', reboot: '重启' }

function getActionLabel(action) {
  return actionLabels[action] || action
}
</script>

<template>
  <section v-if="logs.length > 0" class="results-section">
    <div class="log-header">
      <h2>操作日志</h2>
      <button class="btn btn-sm btn-secondary" @click="emit('clear')">清空日志</button>
    </div>
    <div class="table-wrapper">
      <table class="ecs-table">
        <thead>
          <tr>
            <th>时间</th>
            <th>操作</th>
            <th>InstanceId</th>
            <th>Region</th>
            <th>结果</th>
            <th>消息</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(log, idx) in logs" :key="idx">
            <td>{{ log.time }}</td>
            <td>{{ getActionLabel(log.action) }}</td>
            <td>{{ log.instanceId }}</td>
            <td>{{ log.region }}</td>
            <td>
              <span class="status-badge" :class="log.success ? 'status-running' : 'status-stopped'">
                {{ log.success ? '成功' : '失败' }}
              </span>
            </td>
            <td>{{ log.message }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
