// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';
import { TableListItem } from './data';

/** 获取Job列表 GET /api/job */
export async function getJob(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data: TableListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  }>('/api/v1/job/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  }).then((response) => {
    // 将接口返回的数据格式转换为 Ant Design 表格所需的格式
    if (response.status === 'success') {
      return {
        success: true,
        data: response.data.list, // 列表数据
        total: response.data.count, // 总条数
      };
    } else {
      // 如果接口返回的状态不是 success，返回一个错误格式
      return {
        success: false,
        data: [],
        total: 0,
      };
    }
  });
}

/** 新建规则 PUT /api/rule */
export async function updateRule(data: { [key: string]: any }, options?: { [key: string]: any }) {
  return request<TableListItem>('/api/rule', {
    data,
    method: 'PUT',
    ...(options || {}),
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(data: { [key: string]: any }, options?: { [key: string]: any }) {
  return request<TableListItem>('/api/rule', {
    data,
    method: 'POST',
    ...(options || {}),
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(data: { key: number[] }, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/rule', {
    data,
    method: 'DELETE',
    ...(options || {}),
  });
}
