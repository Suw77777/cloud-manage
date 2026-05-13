<script setup>
const props = defineProps({
  currentPage: { type: Number, required: true },
  totalPages: { type: Number, required: true },
  totalItems: { type: Number, default: 0 },
  pageSize: { type: Number, default: 20 },
  hasMore: { type: Boolean, default: false },
  hasPrev: { type: Boolean, default: false }
})

const emit = defineEmits(['page-change', 'page-size-change'])

const pageSizeOptions = [10, 20, 50, 100]

function goToPage(page) {
  if (page >= 1 && page <= props.totalPages) {
    emit('page-change', page)
  }
}

function onPageSizeChange(e) {
  emit('page-size-change', Number(e.target.value))
}
</script>

<template>
  <div class="pagination-bar">
    <div class="pagination-info">
      <span>共 {{ totalItems }} 条</span>
      <select :value="pageSize" @change="onPageSizeChange" class="page-size-select">
        <option v-for="size in pageSizeOptions" :key="size" :value="size">
          {{ size }} 条/页
        </option>
      </select>
    </div>

    <div class="pagination-controls">
      <button
        class="page-btn"
        :disabled="!hasPrev"
        @click="goToPage(1)"
      >
        «
      </button>
      <button
        class="page-btn"
        :disabled="!hasPrev"
        @click="goToPage(currentPage - 1)"
      >
        ‹
      </button>

      <span class="page-current">
        {{ currentPage }} / {{ totalPages || 1 }}
      </span>

      <button
        class="page-btn"
        :disabled="!hasMore"
        @click="goToPage(currentPage + 1)"
      >
        ›
      </button>
      <button
        class="page-btn"
        :disabled="!hasMore"
        @click="goToPage(totalPages)"
      >
        »
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  font-size: 13px;
  color: #606266;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-size-select {
  padding: 4px 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 12px;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  padding: 4px 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  font-size: 12px;
  color: #606266;
}

.page-btn:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-current {
  padding: 0 12px;
  font-weight: 500;
}
</style>
