<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  items: { type: Array, required: true },
  itemHeight: { type: Number, default: 40 },
  maxHeight: { type: Number, default: 500 },
  buffer: { type: Number, default: 5 }
})

const containerRef = ref(null)
const scrollTop = ref(0)
const containerHeight = ref(0)

const visibleCount = computed(() => Math.ceil(containerHeight.value / props.itemHeight) + props.buffer * 2)
const startIndex = computed(() => Math.max(0, Math.floor(scrollTop.value / props.itemHeight) - props.buffer))
const endIndex = computed(() => Math.min(props.items.length, startIndex.value + visibleCount.value))
const visibleItems = computed(() => props.items.slice(startIndex.value, endIndex.value))
const totalHeight = computed(() => props.items.length * props.itemHeight)
const offsetY = computed(() => startIndex.value * props.itemHeight)

function onScroll(e) {
  scrollTop.value = e.target.scrollTop
}

function updateHeight() {
  if (containerRef.value) {
    containerHeight.value = containerRef.value.clientHeight
  }
}

onMounted(() => {
  updateHeight()
  window.addEventListener('resize', updateHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateHeight)
})

watch(() => props.items.length, () => {
  updateHeight()
})
</script>

<template>
  <div
    ref="containerRef"
    class="virtual-scroller"
    :style="{ maxHeight: maxHeight + 'px', overflow: 'auto' }"
    @scroll="onScroll"
  >
    <div class="virtual-scroller-phantom" :style="{ height: totalHeight + 'px' }">
      <div class="virtual-scroller-content" :style="{ transform: `translateY(${offsetY}px)` }">
        <div
          v-for="(item, index) in visibleItems"
          :key="startIndex + index"
          class="virtual-scroller-item"
          :style="{ height: itemHeight + 'px' }"
        >
          <slot :item="item" :index="startIndex + index"></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.virtual-scroller {
  position: relative;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.virtual-scroller-phantom {
  position: relative;
}

.virtual-scroller-content {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
}

.virtual-scroller-item {
  display: flex;
  align-items: center;
  padding: 0 12px;
  border-bottom: 1px solid #ebeef5;
  box-sizing: border-box;
}

.virtual-scroller-item:hover {
  background: #f5f7fa;
}
</style>
