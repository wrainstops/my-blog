import { ElMessage } from 'element-plus'
import axios from 'axios'
import { TOKENKEY, getToken } from './auth'

const service = axios.create({
  baseURL: '/api',
  withCredentials: false, // 跨域请求时 不携带cookie
  timeout: 5000 // 请求超时时间，超过5000ms 执行错误函数
})

service.interceptors.request.use(
  (config : any) => {
    config.headers[TOKENKEY] = 'Bearer ' + getToken()
    if (config.headers['Content-Type'] && config.headers['Content-Type'].indexOf('form-data') > -1) {
      return config
    }
    const data = config.data
    if (typeof (config.data) !== 'undefined') {
      let newData: any = {}
      const dataArr = Object.keys(data)
      if (dataArr.indexOf('NOFILTERNULL') > -1) {
        newData = config.data
        delete newData.NOFILTERNULL
      } else {
        Object.keys(data).forEach(item => {
          if (data[item] !== '' && data[item] !== undefined && data[item] !== null) {
            newData[item] = data[item]
          }
        })
      }
      config.data = newData
    }
    return config
  },
  (error: any) => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response: any) => {
    const res = response.data
    if (response.status !== 200) {
      return Promise.reject(new Error(response.statusText || 'Error'))
    } else {
      if (response.headers['content-type'] === 'application/octet-stream' || response.headers['content-type'] === 'application/vnd.ms-excel') {
        if (response.headers['content-disposition'].split('filename=')[1]) {
          const fullFilename = decodeURIComponent(response.headers['content-disposition'].split('filename=')[1])
          return { data: res, filename: fullFilename }
        } else {
          throw new Error('返回值没有filename关键字')
        }
      } else {
        if (res.code && res.code !== 200) {
          console.error(res)
          ElMessage.error(res.message)
          return Promise.reject(res.message)
        }
        if (res.data) {
          // 有data包着的返回data
          return res.data
        } else {
          return res
        }
      }
    }
  },
  (error: any) => {
    const { status } = error.response
    let errorMessage = ''
    switch (status) {
      case 403:
        errorMessage = '登录过期,请重新登录'
        break
      case 404:
        errorMessage = '接口不存在,请核实请求路径'
        break
      case 401:
        errorMessage = error.response.data.message || error.response.statusText
        break
      default:
        errorMessage = error.response.data.message || error.response.statusText
    }
    ElMessage.error(errorMessage)
    return Promise.reject(error)
  }
)

export default service
export const $get = service.get
export const $post = service.post
export const $delete = service.delete
