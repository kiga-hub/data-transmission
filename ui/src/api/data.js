import axios from '@/libs/api.request'

export const updateModelList = (query) => {
  return axios.request({
    url: '/upgrade/update_list',
    method: 'get',
    params: query
  })
}

export const getModelList = () => {
  return axios.request({
    url: '/upgrade/source_list',
    method: 'get'
  })
}

export const getUpgradeList = () => {
  return axios.request({
    url: '/upgrade/log_list',
    method: 'get'
  })
}

export const getUpgradeDetail = (query) => {
  return axios.request({
    url: '/upgrade/detail',
    method: 'get',
    params: query
  })
}

export const startUpgrade = (data) => {
  return axios.request({
    url: '/upgrade/start',
    method: 'post',
    data: data
  })
}

export const deleteUpgradeDetail = (data) => {
  return axios.request({
    url: '/upgrade/delete',
    method: 'delete',
    params: data
  })
}
