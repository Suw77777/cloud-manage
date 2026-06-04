import { ref, computed } from 'vue'

export function useConfig() {
  const profiles = ref([])
  const currentProfile = ref('')
  const loading = ref(false)
  const error = ref('')

  async function loadProfiles() {
    loading.value = true
    error.value = ''
    try {
      // Use Wails binding if available
      if (window.go?.main?.App?.ListConfigProfiles) {
        const result = await window.go.main.App.ListConfigProfiles()
        profiles.value = result || []
      }
      if (window.go?.main?.App?.GetCurrentProfile) {
        currentProfile.value = await window.go.main.App.GetCurrentProfile() || ''
      }
    } catch (e) {
      error.value = e.message || '加载配置失败'
    } finally {
      loading.value = false
    }
  }

  async function switchProfile(name) {
    loading.value = true
    error.value = ''
    try {
      if (window.go?.main?.App?.SwitchProfile) {
        await window.go.main.App.SwitchProfile(name)
        currentProfile.value = name
      }
    } catch (e) {
      error.value = e.message || '切换账号失败'
    } finally {
      loading.value = false
    }
  }

  async function getProfileCredentials(name) {
    try {
      if (window.go?.main?.App?.GetProfileCredentials) {
        return await window.go.main.App.GetProfileCredentials(name)
      }
    } catch (e) {
      error.value = e.message || '获取凭证失败'
    }
    return null
  }

  return {
    profiles,
    currentProfile,
    loading,
    error,
    loadProfiles,
    switchProfile,
    getProfileCredentials,
  }
}
