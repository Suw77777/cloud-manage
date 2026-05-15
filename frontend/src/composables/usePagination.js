import { ref, computed } from 'vue'

export function usePagination(options = {}) {
  const pageSize = ref(options.pageSize || 20)
  const currentPage = ref(1)
  const totalItems = ref(0)

  const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value))
  const hasMore = computed(() => currentPage.value < totalPages.value)
  const hasPrev = computed(() => currentPage.value > 1)

  const startIndex = computed(() => (currentPage.value - 1) * pageSize.value)
  const endIndex = computed(() => Math.min(startIndex.value + pageSize.value, totalItems.value))

  function setPage(page) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page
    }
  }

  function nextPage() {
    if (hasMore.value) {
      currentPage.value++
    }
  }

  function prevPage() {
    if (hasPrev.value) {
      currentPage.value--
    }
  }

  function setTotal(total) {
    totalItems.value = total
    // 如果当前页超出范围，重置到第一页
    if (currentPage.value > totalPages.value && totalPages.value > 0) {
      currentPage.value = 1
    }
  }

  function reset() {
    currentPage.value = 1
    totalItems.value = 0
  }

  function setPageSize(size) {
    pageSize.value = size
    currentPage.value = 1
  }

  return {
    pageSize,
    currentPage,
    totalItems,
    totalPages,
    hasMore,
    hasPrev,
    startIndex,
    endIndex,
    setPage,
    nextPage,
    prevPage,
    setTotal,
    reset,
    setPageSize
  }
}
