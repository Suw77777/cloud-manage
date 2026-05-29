import { ref } from 'vue'
import { ListSLBs, GetSLBDetail, ListSLBListeners } from '../../wailsjs/go/main/App'

export function useSLB(accessKeyId, accessKeySecret) {
  const region = ref('cn-hangzhou')
  const slbs = ref([])
  const slbDetail = ref(null)
  const listeners = ref([])

  const loading = ref(false)
  const loadingDetail = ref(false)
  const error = ref('')
  const success = ref('')

  async function fetchSLBs() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loading.value = true
    error.value = ''
    success.value = ''

    try {
      const result = await ListSLBs(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value
      )
      if (result.success) {
        slbs.value = result.slbs || []
        success.value = result.message
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 SLB 列表失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  async function fetchSLBDetail(slbId) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loadingDetail.value = true
    error.value = ''

    try {
      const result = await GetSLBDetail(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        slbId
      )
      if (result.success) {
        slbDetail.value = result.detail
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 SLB 详情失败: ' + (err.message || err)
    } finally {
      loadingDetail.value = false
    }
  }

  async function fetchListeners(slbId) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loading.value = true
    error.value = ''

    try {
      const result = await ListSLBListeners(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        slbId
      )
      if (result.success) {
        listeners.value = result.listeners || []
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取监听器列表失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  function clearResults() {
    slbs.value = []
    slbDetail.value = null
    listeners.value = []
    error.value = ''
    success.value = ''
  }

  return {
    region,
    slbs,
    slbDetail,
    listeners,
    loading,
    loadingDetail,
    error,
    success,
    fetchSLBs,
    fetchSLBDetail,
    fetchListeners,
    clearResults
  }
}
