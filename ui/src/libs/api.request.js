import HttpRequest from '@/libs/axios'
var baseUrl = '/api/data/v1'
if (process.env.NODE_ENV === 'development') {
  baseUrl = 'http://192.168.8.244:8889/api'
}
const axios = new HttpRequest(baseUrl)
export default axios
