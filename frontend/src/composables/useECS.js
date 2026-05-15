import { ref, computed } from 'vue'
import {
  QueryECSMultiRegion,
  GetECSDetail,
  StartECS,
  StopECS,
  RebootECS,
} from '../../wailsjs/go/main/App'

export function useECS() {
  const accessKeyId = ref('')
  const accessKeySecret = ref('')
  const selectedRegions = ref(['cn-hangzhou'])
  const env = ref('dev')
  const regionResults = ref([])
  const errorMessage = ref('')
  const successMessage = ref('')
  const loading = ref(false)

  // Detail
  const showDetail = ref(false)
  const detailLoading = ref(false)
  const detailData = ref(null)

  // Confirm
  const showConfirm = ref(false)
  const confirmAction = ref('')
  const confirmInstanceId = ref('')
  const confirmInstanceName = ref('')
  const confirmRegion = ref('')

  // Operation log
  const operationLogs = ref([])

  const totalInstances = computed(() => {
    return regionResults.value.reduce((sum, r) => sum + (r.instances ? r.instances.length : 0), 0)
  })

  function addLog(action, instanceId, region, success, message) {
    operationLogs.value.unshift({
      time: new Date().toLocaleTimeString(),
      action,
      instanceId,
      region,
      success,
      message,
    })
    if (operationLogs.value.length > 50) {
      operationLogs.value.pop()
    }
  }

  async function queryECS() {
    errorMessage.value = ''
    successMessage.value = ''
    regionResults.value = []

    if (!accessKeyId.value.trim()) {
      errorMessage.value = '请输入 AccessKey ID'
      return
    }
    if (!accessKeySecret.value.trim()) {
      errorMessage.value = '请输入 AccessKey Secret'
      return
    }
    if (selectedRegions.value.length === 0) {
      errorMessage.value = '请至少选择一个 Region'
      return
    }

    loading.value = true
    try {
      const result = await QueryECSMultiRegion(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        selectedRegions.value,
        env.value
      )
      if (result.success) {
        regionResults.value = result.regions || []
        successMessage.value = result.message
      } else {
        errorMessage.value = result.message
      }
    } catch (err) {
      errorMessage.value = '调用失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  async function viewDetail(instanceId, region) {
    showDetail.value = true
    detailLoading.value = true
    detailData.value = null

    try {
      const result = await GetECSDetail(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region,
        instanceId
      )
      if (result.success) {
        detailData.value = JSON.parse(result.message)
      } else {
        detailData.value = { error: result.message }
      }
    } catch (err) {
      detailData.value = { error: err.message || String(err) }
    } finally {
      detailLoading.value = false
    }
  }

  function closeDetail() {
    showDetail.value = false
    detailData.value = null
  }

  function requestAction(action, instanceId, instanceName, region) {
    confirmAction.value = action
    confirmInstanceId.value = instanceId
    confirmInstanceName.value = instanceName
    confirmRegion.value = region
    showConfirm.value = true
  }

  function cancelConfirm() {
    showConfirm.value = false
    confirmAction.value = ''
    confirmInstanceId.value = ''
    confirmInstanceName.value = ''
    confirmRegion.value = ''
  }

  async function executeConfirm() {
    const action = confirmAction.value
    const instanceId = confirmInstanceId.value
    const region = confirmRegion.value

    showConfirm.value = false

    try {
      let result
      if (action === 'start') {
        result = await StartECS(accessKeyId.value.trim(), accessKeySecret.value.trim(), region, instanceId)
      } else if (action === 'stop') {
        result = await StopECS(accessKeyId.value.trim(), accessKeySecret.value.trim(), region, instanceId, false)
      } else if (action === 'reboot') {
        result = await RebootECS(accessKeyId.value.trim(), accessKeySecret.value.trim(), region, instanceId, false)
      }

      if (result) {
        addLog(action, instanceId, region, result.success, result.message)
        if (result.success) {
          successMessage.value = result.message
        } else {
          errorMessage.value = result.message
        }
      }
    } catch (err) {
      addLog(action, instanceId, region, false, err.message || String(err))
      errorMessage.value = '操作失败: ' + (err.message || err)
    }

    cancelConfirm()
  }

  function clearInputs() {
    accessKeyId.value = ''
    accessKeySecret.value = ''
    selectedRegions.value = ['cn-hangzhou']
    env.value = 'dev'
    errorMessage.value = ''
    successMessage.value = ''
  }

  function clearResults() {
    regionResults.value = []
    errorMessage.value = ''
    successMessage.value = ''
  }

  function clearLogs() {
    operationLogs.value = []
  }

  return {
    // State
    accessKeyId,
    accessKeySecret,
    selectedRegions,
    env,
    regionResults,
    errorMessage,
    successMessage,
    loading,
    totalInstances,
    // Detail
    showDetail,
    detailLoading,
    detailData,
    viewDetail,
    closeDetail,
    // Confirm
    showConfirm,
    confirmAction,
    confirmInstanceId,
    confirmInstanceName,
    confirmRegion,
    requestAction,
    cancelConfirm,
    executeConfirm,
    // Log
    operationLogs,
    clearLogs,
    // Actions
    queryECS,
    clearInputs,
    clearResults,
  }
}
