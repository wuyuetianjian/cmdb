const AUTH_STORAGE_KEY = 'cmdb_auth_user'

export const getAuthUser = (): string | null => {
  return localStorage.getItem(AUTH_STORAGE_KEY)
}

export const setAuthUser = (username: string): void => {
  localStorage.setItem(AUTH_STORAGE_KEY, username)
}

const AUTH_TAGS_KEY = 'cmdb_auth_tags'

export const setAuthTags = (tags: string[]): void => {
  localStorage.setItem(AUTH_TAGS_KEY, JSON.stringify(tags))
}

export const getAuthTags = (): string[] => {
  const raw = localStorage.getItem(AUTH_TAGS_KEY)
  if (!raw) {
    return []
  }
  try {
    return JSON.parse(raw) as string[]
  } catch (error) {
    return []
  }
}

export const clearAuthUser = (): void => {
  localStorage.removeItem(AUTH_STORAGE_KEY)
  localStorage.removeItem(AUTH_TAGS_KEY)
}

export const hasSSOCookie = (): boolean => {
  return document.cookie.split(';').some((item) => item.trim().startsWith('sso_session='))
}

export const hasLocalCookie = (): boolean => {
  return document.cookie.split(';').some((item) => item.trim().startsWith('local_session='))
}

export const isAuthenticated = (): boolean => {
  return Boolean(getAuthUser()) || hasSSOCookie() || hasLocalCookie()
}
