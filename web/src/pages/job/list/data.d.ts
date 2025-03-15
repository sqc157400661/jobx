export type TableListItem = {
  id: number;
  root_id: number;
  disabled?: boolean;
  href: string;
  avatar: string;
  name: string;
  owner: string;
  desc: string;
  callNo: number;
  biz_id: string;
  locker: string;
  tenant: string;
  phase: string;
  reason:string;
  input:object;
  env:object;
  status: string;
  create_at: Date;
  update_at: Date;
  progress: number;
  pipelines: Pipeline[];
};

export type Pipeline = {
  name: string;
  action: string;
  owner: string;
  desc: string;
  reason: string;
  status: string;
  updatedAt: Date;
  createdAt: Date;
};

export type TableListPagination = {
  total: number;
  pageSize: number;
  current: number;
};

export type TableListData = {
  list: TableListItem[];
  pagination: Partial<TableListPagination>;
};

export type TableListParams = {
  status?: string;
  name?: string;
  desc?: string;
  key?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};
