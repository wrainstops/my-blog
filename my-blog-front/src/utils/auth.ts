export const TOKENKEY = 'Authorization'
export const USERINFO = 'Info'

export function setToken(token: string) {
  return localStorage.setItem(TOKENKEY, token)
}

export function getToken() {
  return localStorage.getItem(TOKENKEY)
}

export function removeToken() {
  return localStorage.removeItem(TOKENKEY)
}

export function setLocalUserInfo(userInfo: any) {
  return localStorage.setItem(USERINFO, JSON.stringify(userInfo))
}

export function getLocalUserInfo() {
  const userInfo = localStorage.getItem(USERINFO)
  return userInfo ? JSON.parse(userInfo) : null
}

export function removeLocalUserInfo() {
  return localStorage.removeItem(USERINFO)
}
