import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import moment from 'moment';
import {
  FooterToolbar,
  ModalForm,
  PageContainer,
  ProDescriptions,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, Input, message,Modal,Tooltip } from 'antd';
import React, { useRef, useState } from 'react';
import type { FormValueType } from './components/UpdateForm';
import UpdateForm from './components/UpdateForm';
import ReactJson from 'react-json-view';
import type { TableListItem, TableListPagination,Pipeline } from './data';
import { addRule, removeRule, getJob, updateRule,getJobPipelines } from './service';
/**
 * 添加节点
 *
 * @param fields
 */

const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');

  try {
    await addRule({ ...fields });
    hide();
    message.success('添加成功');
    return true;
  } catch (error) {
    hide();
    message.error('添加失败请重试！');
    return false;
  }
};
/**
 * 更新节点
 *
 * @param fields
 */

const handleUpdate = async (fields: FormValueType, currentRow?: TableListItem) => {
  const hide = message.loading('正在配置');

  try {
    await updateRule({
      ...currentRow,
      ...fields,
    });
    hide();
    message.success('配置成功');
    return true;
  } catch (error) {
    hide();
    message.error('配置失败请重试！');
    return false;
  }
};
/**
 * 删除节点
 *
 * @param selectedRows
 */

const handleRemove = async (selectedRows: TableListItem[]) => {
  const hide = message.loading('正在删除');
  if (!selectedRows) return true;

  try {
    await removeRule({
      key: selectedRows.map((row) => row.id),
    });
    hide();
    message.success('删除成功，即将刷新');
    return true;
  } catch (error) {
    hide();
    message.error('删除失败，请重试');
    return false;
  }
};

const TableList: React.FC = () => {
  /** 新建窗口的弹窗 */
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  /** 分布更新窗口的弹窗 */

  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<TableListItem>();
  const [selectedRowsState, setSelectedRows] = useState<TableListItem[]>([]);

  // drawer
  const [showDrawerDetail, setShowDrawerDetail] = useState<boolean>(false);
  const [drawerContentType, setDrawerContentType] = useState<'detail' | 'log' | 'json'>('detail');
  const [jsonData, setJsonData] = useState(null);
  const [logContent, setLogContent] = useState('');
  const [drawerConfig, setDrawerConfig] = useState<{ title: string; width: number }>({
    title: '默认标题',
    width: 600,
  });
  // model
  const [visibleJsonModel, setVisibleJsonModel] = useState(false);
  const [modelConfig, setModelConfig] = useState<{ title: string; width: number }>({
    title: '重试操作',
    width: 800,
  });

  // 存储所有二级数据（按主表行ID索引）
  const [subTableData, setSubTableData] = useState<Record<number, Pipeline[]>>({});
  // 加载状态
  const [loadingSubTableIds, setLoadingSubTableIds] = useState<number[]>([]);

  // 修改后的展开事件处理
  const handleExpand = async (expanded: boolean, record: TableListItem) => {
    if (expanded) { // 只要展开就请求，无论是否已有数据
      try {
        setLoadingSubTableIds((prev) => [...prev, record.id]);
        const response = await getJobPipelines({jobId:record.id});
        setSubTableData((prev) => ({
          ...prev,
          [record.id]: response.data, // 总是覆盖旧数据
        }));
      } catch (error) {
        console.error('Failed to load pipelines:', error);
      } finally {
        setLoadingSubTableIds((prev) => prev.filter((id) => id !== record.id));
      }
    }
  };

  /** 国际化配置 */

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '任务名称',
      dataIndex: 'name',
      tip: '任务名称',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDrawerDetail(true);
              setDrawerConfig({title: '查看规则', width: 600});
              setDrawerContentType('detail');
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: '描述',
      dataIndex: 'desc',
      valueType: 'textarea',
    },
    {
      title: '执行人',
      dataIndex: 'owner',
      hideInForm: true,
    },
    {
      title: '执行服务器',
      dataIndex: 'locker',
      hideInForm: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
      valueEnum: {
        0: {
          text: '废弃',
          status: 'Default',
        },
        1: {
          text: '运行中',
          status: 'Processing',
        },
        2: {
          text: '成功',
          status: 'Success',
        },
        3: {
          text: '异常',
          status: 'Error',
        },
        'fail': {
          text: '异常',
          status: 'Error',
        },
        "pending": {
          text: '运行中',
          status: 'Processing',
        },
        'success': {
          text: '成功',
          status: 'Success',
        },
      },
    },
    {
      title: '结束时间',
      sorter: true,
      dataIndex: 'update_at',
      valueType: 'dateTime',
      render: (_, record) => {
        // 将秒级时间戳转为毫秒
        const timestamp = record.update_at * 1000;
        return moment(timestamp).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        <a
          key="config"
          onClick={() => {
            handleUpdateModalVisible(true);
            setCurrentRow(record);
          }}
        >
          暂停
        </a>,
        <a key="subscribeAlert" href="https://procomponents.ant.design/">
          回滚
        </a>,
      ],
    },
  ];

  /** 子表列配置 */
  const pipelineColumns: ProColumns<Pipeline>[] = [
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '原因',
      dataIndex: 'reason',
      render: (text) => {
        if (!text) return '-'; // 处理空值

        const maxLength = 30;
        const isOverflow = text.length > maxLength;
        const truncatedText = isOverflow ? `${text.slice(0, maxLength)}...` : text;

        return isOverflow ? (
          <Tooltip title={text}>
            <span style={{ cursor: 'pointer' }}>{truncatedText}</span>
          </Tooltip>
        ) : (
          <span>{truncatedText}</span>
        );
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
    },
    {
      title: '更新时间',
      dataIndex: 'update_at',
      valueType: 'dateTime',
      render: (_, record) => {
        // 将秒级时间戳转为毫秒
        const timestamp = record.update_at * 1000;
        return moment(timestamp).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      title: '操作',
      key: 'operation',
      valueType: 'option',
      render: () => [
        <a key="Input" onClick={() => {
          // 模拟获取日志内容
          const mockJson = `{"name": "John", "age": 30, "city": "New York"}`;
          setJsonData(JSON.parse(mockJson));
          setShowDrawerDetail(true);
          setDrawerContentType('json');
          setDrawerConfig({title: '查看入参', width: 600});
        }}>入参</a>,
        <a key="output" onClick={() => {
          // 模拟获取日志内容
          const mockJson = `{"name": "John", "age": 30, "city": "New York"}`;
          setJsonData(JSON.parse(mockJson));
          setShowDrawerDetail(true);
          setDrawerContentType('json');
          setDrawerConfig({title: '查看出参', width: 600});
        }}>出参</a>,
        <a key="log" onClick={() => {
          // 模拟获取日志内容
          const log = `2023-10-01 12:00:00 [INFO] Starting server...\n2023-10-01 12:00:05 [INFO] Server started successfully.\n`;
          setLogContent(log);
          setShowDrawerDetail(true);
          setDrawerContentType('log');
          setDrawerConfig({title: '查看日志', width: 600});
        }}>日志</a>,
        <a key="retry" onClick={() => {
          // 模拟获取日志内容
          const mockJson = `{"name": "John", "age": 30, "city": "New York"}`;
          setJsonData(JSON.parse(mockJson));
          setVisibleJsonModel(true);
        }}>重试</a>,
        <a key="log">跳过</a>,
        <a key="discard">废弃</a>,
      ],
    },
  ];

  const expandedRowRender = (record: TableListItem) => {
    const pipelines = subTableData[record.id] || [];
    const loading = loadingSubTableIds.includes(record.id);
    return (
      <ProTable
        columns={pipelineColumns}
        headerTitle={false}
        search={false}
        options={false}
        dataSource={pipelines}
        loading={loading}
        pagination={false}
      />
    );
  };

  /** 根据 drawerContentType 渲染不同的内容 */
  const renderDrawerContent = () => {
    switch (drawerContentType) {
      case 'detail':
        return (
          currentRow?.name && (
            <ProDescriptions<TableListItem>
              column={2}
              title={currentRow?.name}
              request={async () => ({
                data: currentRow || {},
              })}
              params={{
                id: currentRow?.name,
              }}
              columns={columns as ProDescriptionsItemProps<TableListItem>[]}
            />
          )
        );
      case 'log':
        return (
          <pre style={{ backgroundColor: '#000', color: '#fff', padding: '16px', borderRadius: '4px' }}>
            {logContent}
          </pre>
        );
      case 'json':
        return <ReactJson src={jsonData || {}} theme="apathy" />;
      default:
        return null;
    }
  };

  // 提交修改后的 JSON 数据
  const handleRetry = () => {
    console.log('提交的JSON数据:', jsonData);
    setVisibleJsonModel(false);
  };

  return (
    <PageContainer>
      <ProTable<TableListItem, TableListPagination>
        headerTitle="查询表格"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              handleModalVisible(true);
            }}
          >
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={getJob}
        columns={columns}
        expandable={{
          expandedRowRender,
          onExpand: handleExpand, // 监听展开事件
        }}
        // rowSelection={{
        //   onChange: (_, selectedRows) => {
        //     setSelectedRows(selectedRows);
        //   },
        // }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择{' '}
              <a
                style={{
                  fontWeight: 600,
                }}
              >
                {selectedRowsState.length}
              </a>{' '}
              项 &nbsp;&nbsp;
              <span>
                服务调用次数总计 {selectedRowsState.reduce((pre, item) => pre + item.callNo!, 0)} 万
              </span>
            </div>
          }
        >
          <Button
            onClick={async () => {
              await handleRemove(selectedRowsState);
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
          >
            批量删除
          </Button>
          <Button type="primary">批量审批</Button>
        </FooterToolbar>
      )}
      <ModalForm
        title="新建规则"
        width="400px"
        open={createModalVisible}
        onVisibleChange={handleModalVisible}
        onFinish={async (value) => {
          const success = await handleAdd(value as TableListItem);
          if (success) {
            handleModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '规则名称为必填项',
            },
          ]}
          width="md"
          name="name"
        />
        <ProFormTextArea width="md" name="desc" />
      </ModalForm>
      <UpdateForm
        onSubmit={async (value) => {
          const success = await handleUpdate(value, currentRow);

          if (success) {
            handleUpdateModalVisible(false);
            setCurrentRow(undefined);

            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
        onCancel={() => {
          handleUpdateModalVisible(false);
          setCurrentRow(undefined);
        }}
        updateModalVisible={updateModalVisible}
        values={currentRow || {}}
      />

      <Drawer
        width={drawerConfig.width}
        title={drawerConfig.title}
        open={showDrawerDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDrawerDetail(false);
          setLogContent('');
        }}
        closable={false}
      >
        {renderDrawerContent()}
      </Drawer>

      {/* 模态框 */}
      <Modal
        title={modelConfig.title}
        open={visibleJsonModel}
        onCancel={() => {
          setVisibleJsonModel(false);
        }}
        onOk={handleRetry}
        width={modelConfig.width}
      >
        {/* JSON 编辑器 */}
        <ReactJson
          src={jsonData} // 数据源
          onEdit={(edit) => {
            setJsonData(edit.updated_src); // 更新 JSON 数据
          }}
          onAdd={(add) => {
            setJsonData(add.updated_src); // 添加新字段
          }}
          onDelete={(del) => {
            setJsonData(del.updated_src); // 删除字段
          }}
          name={false} // 不显示根节点名称
          displayDataTypes={false} // 不显示数据类型
          enableClipboard={false} // 禁用剪贴板功能
          collapsed={false} // 默认展开所有节点
        />
      </Modal>

</PageContainer>
  );
};

export default TableList;
