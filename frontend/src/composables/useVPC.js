import { ref } from 'vue'
import { ListVPCs, GetVPCDetail, ListVSwitches } from '../../wailsjs/go/main/App'

export function useVPC(accessKeyId, accessKeySecret) {
  const region = ref('cn-hangzhou')
  const vpcs = ref([])
  const vpcDetail = ref(null)
  const vswitches = ref([])

  const loading = ref(false)
  const loadingDetail = ref(false)
  const error = ref('')
  const success = ref('')

  async function fetchVPCs() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loading.value = true
    error.value = ''
    success.value = ''

    try {
      const result = await ListVPCs(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value
      )
      if (result.success) {
        vpcs.value = result.vpcs || []
        success.value = result.message
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 VPC 列表失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  async function fetchVPCDetail(vpcId) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loadingDetail.value = true
    error.value = ''

    try {
      const result = await GetVPCDetail(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        vpcId
      )
      if (result.success) {
        vpcDetail.value = result.detail
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 VPC 详情失败: ' + (err.message || err)
    } finally {
      loadingDetail.value = false
    }
  }

  async function fetchVSwitches(vpcId) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loading.value = true
    error.value = ''

    try {
      const result = await ListVSwitches(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        vpcId
      )
      if (result.success) {
        vswitches.value = result.vswitches || []
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 VSwitch 列表失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  function clearResults() {
    vpcs.value = []
    vpcDetail.value = null
    vswitches.value = []
    error.value = ''
    success.value = ''
  }

  return {
    region,
    vpcs,
    vpcDetail,
    vswitches,
    loading,
    loadingDetail,
    error,
    success,
    fetchVPCs,
    fetchVPCDetail,
    fetchVSwitches,
    clearResults
  }
}
