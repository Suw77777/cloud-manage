import { ref } from 'vue'
import { GetECSMetrics, GetECSMetricsMultiRegion } from '../../wailsjs/go/main/App'

export function useCMS(accessKeyId, accessKeySecret) {
  const metricsLoading = ref(false)
  const metricsData = ref([])
  const metricsError = ref('')
  const metricsSuccess = ref('')
  const selectedInstances = ref([])

  async function queryMetrics(instances) {
    if (!accessKeyId.value.trim()) {
      metricsError.value = '请输入 AccessKey ID'
      return
    }
    if (!accessKeySecret.value.trim()) {
      metricsError.value = '请输入 AccessKey Secret'
      return
    }
    if (!instances || instances.length === 0) {
      metricsError.value = '请至少选择一个实例'
      return
    }

    metricsLoading.value = true
    metricsError.value = ''
    metricsSuccess.value = ''
    metricsData.value = []

    try {
      const result = await GetECSMetricsMultiRegion(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        instances
      )
      if (result.success) {
        metricsData.value = result.metrics || []
        metricsSuccess.value = result.message
      } else {
        metricsError.value = result.message
      }
    } catch (err) {
      metricsError.value = '调用失败: ' + (err.message || err)
    } finally {
      metricsLoading.value = false
    }
  }

  async function querySingleMetric(instanceId, region) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      return null
    }

    try {
      const result = await GetECSMetrics(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region,
        instanceId
      )
      if (result.success && result.metrics && result.metrics.length > 0) {
        return result.metrics[0]
      }
      return null
    } catch (err) {
      console.error('Failed to query metrics:', err)
      return null
    }
  }

  function clearMetrics() {
    metricsData.value = []
    metricsError.value = ''
    metricsSuccess.value = ''
  }

  function toggleInstance(instanceId, region) {
    const idx = selectedInstances.value.findIndex(
      i => i.instanceId === instanceId && i.region === region
    )
    if (idx === -1) {
      selectedInstances.value.push({ instanceId, region })
    } else {
      selectedInstances.value.splice(idx, 1)
    }
  }

  function selectAllInstances(instances) {
    selectedInstances.value = instances.map(i => ({
      instanceId: i.instanceId,
      region: i.regionId
    }))
  }

  function clearSelection() {
    selectedInstances.value = []
  }

  return {
    metricsLoading,
    metricsData,
    metricsError,
    metricsSuccess,
    selectedInstances,
    queryMetrics,
    querySingleMetric,
    clearMetrics,
    toggleInstance,
    selectAllInstances,
    clearSelection
  }
}
