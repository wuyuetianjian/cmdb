const AUTH_STORAGE_KEY = 'cmdb_auth_user'

export const getAuthUser = (): string | null => {
  return localStorage.getItem(AUTH_STORAGE_KEY)
}

export const setAuthUser = (username: string): void => {
  localStorage.setItem(AUTH_STORAGE_KEY, username)
}

export const clearAuthUser = (): void => {
  localStorage.removeItem(AUTH_STORAGE_KEY)
}

export const hasSSOCookie = (): boolean => {
  return document.cookie.split(';').some((item) => item.trim().startsWith('sso_session='))
}

export const isAuthenticated = (): boolean => {
  return Boolean(getAuthUser()) || hasSSOCookie()
}
