import { ref } from 'vue'
import { ListOSSBuckets, ListOSSObjects } from '../../wailsjs/go/main/App'

export function useOSS(accessKeyId, accessKeySecret) {
  const region = ref('cn-hangzhou')
  const currentBucket = ref('')
  const currentPrefix = ref('')
  const buckets = ref([])
  const objects = ref([])
  const isTruncated = ref(false)
  const breadcrumb = ref([])

  const loading = ref(false)
  const loadingBuckets = ref(false)
  const error = ref('')
  const success = ref('')

  async function fetchBuckets() {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    loadingBuckets.value = true
    error.value = ''
    success.value = ''

    try {
      const result = await ListOSSBuckets(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value
      )
      if (result.success) {
        buckets.value = result.buckets || []
        success.value = result.message
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取 Bucket 列表失败: ' + (err.message || err)
    } finally {
      loadingBuckets.value = false
    }
  }

  async function fetchObjects(bucket, prefix) {
    if (!accessKeyId.value.trim() || !accessKeySecret.value.trim()) {
      error.value = '请先输入 AccessKey'
      return
    }

    if (bucket) {
      currentBucket.value = bucket
    }
    if (prefix !== undefined) {
      currentPrefix.value = prefix
    }

    loading.value = true
    error.value = ''
    success.value = ''

    try {
      const result = await ListOSSObjects(
        accessKeyId.value.trim(),
        accessKeySecret.value.trim(),
        region.value,
        currentBucket.value,
        currentPrefix.value,
        100
      )
      if (result.success) {
        objects.value = result.objects || []
        isTruncated.value = result.isTruncated
        success.value = result.message
        updateBreadcrumb()
      } else {
        error.value = result.message
      }
    } catch (err) {
      error.value = '获取对象列表失败: ' + (err.message || err)
    } finally {
      loading.value = false
    }
  }

  function navigateToFolder(prefix) {
    currentPrefix.value = prefix
    fetchObjects()
  }

  function navigateToRoot() {
    currentPrefix.value = ''
    fetchObjects()
  }

  function navigateUp() {
    const parts = currentPrefix.value.split('/').filter(p => p)
    if (parts.length > 0) {
      parts.pop()
      currentPrefix.value = parts.length > 0 ? parts.join('/') + '/' : ''
    }
    fetchObjects()
  }

  function updateBreadcrumb() {
    const parts = currentPrefix.value.split('/').filter(p => p)
    breadcrumb.value = [{ name: currentBucket.value, prefix: '' }]
    let currentPath = ''
    for (const part of parts) {
      currentPath += part + '/'
      breadcrumb.value.push({ name: part, prefix: currentPath })
    }
  }

  function selectBucket(bucketName) {
    currentBucket.value = bucketName
    currentPrefix.value = ''
    objects.value = []
    breadcrumb.value = []
    fetchObjects(bucketName, '')
  }

  function backToBuckets() {
    currentBucket.value = ''
    currentPrefix.value = ''
    objects.value = []
    breadcrumb.value = []
  }

  function formatFileSize(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  function clearResults() {
    buckets.value = []
    objects.value = []
    error.value = ''
    success.value = ''
  }

  return {
    region,
    currentBucket,
    currentPrefix,
    buckets,
    objects,
    isTruncated,
    breadcrumb,
    loading,
    loadingBuckets,
    error,
    success,
    fetchBuckets,
    fetchObjects,
    navigateToFolder,
    navigateToRoot,
    navigateUp,
    selectBucket,
    backToBuckets,
    formatFileSize,
    clearResults
  }
}
