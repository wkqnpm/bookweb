/*
 * @Descripttion:
 * @version:
 * @Author: wkq
 * @Date: 2022-06-12 09:27:35
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2022-07-04 16:41:56
 */
import {
  Table, Button, Modal, Form, Input, Popconfirm,
} from 'antd';
import React, { useEffect, useState } from 'react';
import { get, post } from '@/server/request';

const url = '/list';
export default function Mani() {
  const [config, setConf] = useState({
    isShow: false,
    title: '',
  });
  const [data, setData] = useState([]);
  const [form] = Form.useForm();
  const [record, setRecord] = useState(null);
  const [pageConfig, setPageConf] = useState({
    total: 0,
    pageSize: 10,
    pageNum: 1,
  });
  const InitData = (pageNum = 1, pageSize = 10) => {
    get(url, { pageNum, pageSize }).then((res) => {
      if (res && res.code === 200) {
        setData(res.data || []);
        setPageConf({
          pageNum,
          pageSize,
          total: res.total,
        });
      }
    });
  };
  // 确认删除
  const confirmDel = (row: any) => {
    post('/del', { id: row.id }).then((res) => {
      if (res && res.code === 200) {
        InitData();
      }
    });
  };
  const AddData = (params: any) => {
    post('/add', params).then((res) => {
      if (res && res.code === 200) {
        InitData();
        setConf({ title: '', isShow: false });
      }
    });
  };
  function EditData(values) {
    post('/edit', { id: record.id, bookname: values.bookname }).then((res) => {
      if (res && res.code === 200) {
        InitData();
        setConf({ title: '', isShow: false });
      }
    });
  }
  // 提交
  const onFinish = (values: any) => {
    if (config.title === 'Add') {
      AddData(values);
    } else if (config.title === 'Edit') {
      EditData(values);
    }
  };
  // 重置表单
  const onReset = () => {
    form.resetFields();
  };

  useEffect(() => {
    InitData();
  }, []);
  function handleEdit(row) {
    setConf({ isShow: true, title: 'Edit' });
    setRecord(row);
    form.setFieldsValue(row);
  }
  function pageChange(page: number, pageSize:number) {
    InitData(page, pageSize);
  }
  const columns = [
    {
      title: 'BookName',
      dataIndex: 'bookname',
    },
    {
      title: 'Author',
      dataIndex: 'author',
    },
    {
      title: 'CreateTime',
      dataIndex: 'create_time',
      render(val) {
        return <span>{new Date(val * 1000).toLocaleString()}</span>;
      },
    },
    {
      title: 'Actions',
      render(_: any, row: any) {
        return (
          <>
            <Button type="link" onClick={handleEdit.bind(this, row)}>Edit</Button>
            <Popconfirm
              title="Are you sure to delete this task?"
              onConfirm={() => confirmDel(row)}
              onCancel={() => { }}
              okText="Yes"
              cancelText="No"
            >
              <Button type="link" danger>Del</Button>
            </Popconfirm>
          </>
        );
      },
    },
  ];
  return (
    <>
      <Button onClick={() => setConf({ title: 'Add', isShow: true })}>Add</Button>
      <Table
        dataSource={data}
        rowKey="id"
        columns={columns}
        pagination={{
          current: pageConfig.pageNum,
          pageSize: pageConfig.pageSize,
          total: pageConfig.total,
          onChange: pageChange,
          showTotal: (total) => `Total ${total} items`,
          size: 'small',
          showSizeChanger: true,
          showQuickJumper: true,
        }}
      />
      <Modal
        visible={config.isShow}
        footer={null}
        title={config.title}
        onCancel={() => { setConf({ title: '', isShow: false }); }}
        destroyOnClose
      >
        <Form form={form} name="control-hooks" onFinish={onFinish}>
          <Form.Item name="bookname" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          {
             config.title === 'Add'
               ? (
                 <Form.Item name="author" label="Author" rules={[{ required: true }]}>
                   <Input />
                 </Form.Item>
               ) : null
          }

          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ marginRight: '10px' }}>
              Submit
            </Button>
            <Button htmlType="button" onClick={onReset}>
              Reset
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}
