<script setup>
const props = defineProps({
  action: { type: String, default: '' },
  instanceId: { type: String, default: '' },
  instanceName: { type: String, default: '' },
  region: { type: String, default: '' },
  isProd: { type: Boolean, default: false },
})

const emit = defineEmits(['confirm', 'cancel'])

const actionLabels = { start: '启动', stop: '停止', reboot: '重启' }

function getActionLabel(action) {
  return actionLabels[action] || action
}
</script>

<template>
  <div class="modal-overlay" @click.self="emit('cancel')">
    <div class="modal modal-confirm">
      <div class="modal-header">
        <h3>操作确认</h3>
        <button class="modal-close" @click="emit('cancel')">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="isProd" class="prod-warning">
          当前环境为生产环境，操作不可逆，请谨慎确认！
        </div>
        <p class="confirm-text">
          确定要<strong>{{ getActionLabel(action) }}</strong>实例吗？
        </p>
        <div class="confirm-detail">
          <div>InstanceId: {{ instanceId }}</div>
          <div>InstanceName: {{ instanceName }}</div>
          <div>Region: {{ region }}</div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="emit('cancel')">取消</button>
        <button
          class="btn"
          :class="isProd ? 'btn-danger' : 'btn-primary'"
          @click="emit('confirm')"
        >
          {{ isProd ? '确认执行（生产环境）' : '确认' }}
        </button>
      </div>
    </div>
  </div>
</template>
