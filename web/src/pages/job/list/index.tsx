import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  ModalForm,
  PageContainer,
  ProDescriptions,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, Input, message } from 'antd';
import React, { useRef, useState } from 'react';
import type { FormValueType } from './components/UpdateForm';
import UpdateForm from './components/UpdateForm';
import ReactJson from 'react-json-view';
import type { TableListItem, TableListPagination,Pipeline } from './data';
import { addRule, removeRule, getJob, updateRule } from './service';
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
      key: selectedRows.map((row) => row.key),
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
  /** 国际化配置 */

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '规则名称',
      dataIndex: 'name',
      tip: '规则名称是唯一的 key',
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
      title: '服务调用次数',
      dataIndex: 'callNo',
      sorter: true,
      hideInForm: true,
      renderText: (val: string) => `${val}万`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
      valueEnum: {
        0: {
          text: '关闭',
          status: 'Default',
        },
        1: {
          text: '运行中',
          status: 'Processing',
        },
        2: {
          text: '已上线',
          status: 'Success',
        },
        3: {
          text: '异常',
          status: 'Error',
        },
      },
    },
    {
      title: '上次调度时间',
      sorter: true,
      dataIndex: 'updatedAt',
      valueType: 'dateTime',
      renderFormItem: (item, { defaultRender, ...rest }, form) => {
        const status = form.getFieldValue('status');

        if (`${status}` === '0') {
          return false;
        }

        if (`${status}` === '3') {
          return <Input {...rest} placeholder="请输入异常原因！" />;
        }

        return defaultRender(item);
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
          配置
        </a>,
        <a key="subscribeAlert" href="https://procomponents.ant.design/">
          订阅警报
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
    },
    {
      title: '状态',
      dataIndex: 'status',
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      valueType: 'dateTime',
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
        <a key="retry">重试</a>,
        <a key="log">跳过</a>,
        <a key="discard">废弃</a>,
      ],
    },
  ];

  const expandedRowRender = (record: TableListItem) => {
    return (
      <ProTable
        columns={pipelineColumns}
        headerTitle={false}
        search={false}
        options={false}
        dataSource={record.pipelines}
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

  return (
    <PageContainer>
      <ProTable<TableListItem, TableListPagination>
        headerTitle="查询表格"
        actionRef={actionRef}
        rowKey="key"
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

</PageContainer>
  );
};

export default TableList;
