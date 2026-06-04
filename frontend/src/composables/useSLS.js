import { ref } from 'vue'
import { ListSLSLogStores, QuerySLSLogs, QuerySLSLogsStream, ExportSLSLogs } from '../../wailsjs/go/main/App'

export function useSLS(accessKeyId, accessKeySecret) {
  const region = ref('cn-hangzhou')
  const project = ref('')
  const logstore = ref('')
  const query = ref('')
  const timeRange = ref('1h')
  const maxLines = ref(100)

  const logStores = ref([])
  const logs = ref([])
  const logCount = ref(0)
  const hasMore = ref(false)

  const loading = ref(false)
  const loadingStores = ref(false)
  const streaming = ref(false)
  const streamProgress = ref(0)
  const error = ref('')
  const success = ref('')

  // 分页相关
  const pageSize = ref(50)
  const currentPage = ref(1)
  const totalLogs = ref(0)

  async function fetchLogStores() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }
    if (!project.value.trim()) {
      error.value = '请输入 Project 名称'
      return
    }

    loadingStores.value = true
    error.value = ''

    try {
      const result = await ListSLSLogStores(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        project.value.trim()
      )
      if (result.success) {
        logStores.value = (result.logStores || []).map(name => ({ logstoreName: name }))
        if (logStores.value.length > 0 && !logstore.value) {
          logstore.value = logStores.value[0].logstoreName
        }
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 Logstore 失败: ' + (err.message || err)
    } finally {
      loadingStores.value = false
    }
  }

  function getTimeRange() {
    const now = Math.floor(Date.now() / 1000)
    let from = now - 3600

    switch (timeRange.value) {
      case '15m': from = now - 900; break
      case '30m': from = now - 1800; break
      case '1h': from = now - 3600; break
      case '6h': from = now - 21600; break
      case '12h': from = now - 43200; break
      case '24h': from = now - 86400; break
      case '7d': from = now - 604800; break
    }

    return { from, to: now }
  }

  async function queryLogs() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }
    if (!project.value.trim()) {
      error.value = '请输入 Project 名称'
      return
    }
    if (!logstore.value) {
      error.value = '请选择 Logstore'
      return
    }

    loading.value = true
    error.value = ''
    success.value = ''
    logs.value = []
    currentPage.value = 1

    const { from, to } = getTimeRange()

    try {
      const result = await QuerySLSLogs(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        project.value.trim(),
        logstore.value,
        query.value || '*',
        from,
        to,
        maxLines.value
      )
      if (result.success) {
        logs.value = result.logs || []
        logCount.value = result.count
        totalLogs.value = logs.value.length
        hasMore.value = result.hasMore
        success.value = result.message
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '查询日志失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  // 流式查询 - 边查边显示
  async function queryLogsStream() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }
    if (!project.value.trim()) {
      error.value = '请输入 Project 名称'
      return
    }
    if (!logstore.value) {
      error.value = '请选择 Logstore'
      return
    }

    streaming.value = true
    error.value = ''
    success.value = ''
    logs.value = []
    streamProgress.value = 0

    const { from, to } = getTimeRange()

    try {
      // 使用流式查询
      const result = await QuerySLSLogsStream(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        project.value.trim(),
        logstore.value,
        query.value || '*',
        from,
        to,
        maxLines.value,
        (progress, newLogs) => {
          // 回调函数：更新进度和追加日志
          streamProgress.value = progress
          if (newLogs && newLogs.length > 0) {
            logs.value = [...logs.value, ...newLogs]
            totalLogs.value = logs.value.length
          }
        }
      )

      if (result.success) {
        success.value = `流式查询完成，共 ${logs.value.length} 条日志`
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '流式查询失败: ' + (err.message || err)
    } finally {
      streaming.value = false
      streamProgress.value = 100
    }
  }

  // 分页获取当前页数据
  function getPagedLogs() {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return logs.value.slice(start, end)
  }

  function setPage(page) {
    currentPage.value = page
  }

  function setPageSize(size) {
    pageSize.value = size
    currentPage.value = 1
  }

  function clearResults() {
    logs.value = []
    logCount.value = 0
    totalLogs.value = 0
    hasMore.value = false
    error.value = ''
    success.value = ''
    currentPage.value = 1
    streamProgress.value = 0
  }

  function clearAll() {
    project.value = ''
    logstore.value = ''
    query.value = ''
    logStores.value = []
    clearResults()
  }

  async function exportLogs(format = 'csv') {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }
    if (!project.value.trim()) {
      error.value = '请输入 Project 名称'
      return
    }
    if (!logstore.value) {
      error.value = '请选择 Logstore'
      return
    }

    loading.value = true
    error.value = ''
    success.value = ''

    const { from, to } = getTimeRange()
    const fromStr = new Date(from * 1000).toISOString()
    const toStr = new Date(to * 1000).toISOString()

    try {
      const result = await ExportSLSLogs(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        project.value.trim(),
        logstore.value,
        query.value || '*',
        fromStr,
        toStr,
        format,
        maxLines.value
      )
      if (result.filePath) {
        success.value = `导出成功: ${result.filePath} (${result.count} 条)`
      } else {
        error.value = '导出失败'
      }
    } catch (err) {
      error.value = '导出失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  return {
    region,
    project,
    logstore,
    query,
    timeRange,
    maxLines,
    logStores,
    logs,
    logCount,
    hasMore,
    loading,
    loadingStores,
    streaming,
    streamProgress,
    error,
    success,
    // 分页
    pageSize,
    currentPage,
    totalLogs,
    fetchLogStores,
    queryLogs,
    queryLogsStream,
    exportLogs,
    getPagedLogs,
    setPage,
    setPageSize,
    clearResults,
    clearAll
  }
}
