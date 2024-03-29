import React from 'react'
import { useTranslation } from 'react-i18next'

import { Form, Input, Button, InputNumber } from '~/uikit'
import { StepCompnent } from '~/components/CreateTaskConfig/shared'

const IPv4Pattern =
  /\b(?:(?:2(?:[0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9])\.){3}(?:(?:2([0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9]))\b/

const TargetInfo: StepCompnent = ({ prev, initialValues }) => {
  const [t] = useTranslation()
  return (
    <Form
      wrapperCol={{ span: 16 }}
      labelCol={{ span: 4 }}
      initialValues={initialValues}
      validateTrigger="onBlur"
      name="targetInfo"
    >
      <div className="flex">
        <div className="flex-1">
          <Form.Item
            label={t('ip')}
            name={['target_config', 'host']}
            rules={[
              { required: true, message: t('host is required') },
              {
                pattern: IPv4Pattern,
                message: t('host bad format'),
              },
            ]}
          >
            <Input placeholder="1.1.1.1" />
          </Form.Item>

          <Form.Item
            label={t('port')}
            name={['target_config', 'port']}
            rules={[{ required: true, message: t('port is required') }]}
          >
            <InputNumber min={1} max={65535} />
          </Form.Item>

          <Form.Item
            label={t('user name')}
            name={['target_config', 'user']}
            rules={[
              {
                required: true,
                whitespace: true,
                message: t('user name is required'),
              },
            ]}
          >
            <Input placeholder="root" />
          </Form.Item>

          <Form.Item label={t('password')} name={['target_config', 'password']}>
            <Input type="password" />
          </Form.Item>
        </div>
        <div className="flex-1">
          <Form.Item
            label="SSL CA"
            name={['target_config', 'security', 'ssl_ca_content']}
          >
            <Input.TextArea />
          </Form.Item>

          <Form.Item
            label="SSL CERT"
            name={['target_config', 'security', 'ssl_cert_content']}
          >
            <Input.TextArea />
          </Form.Item>

          <Form.Item
            label="SSL KEY"
            name={['target_config', 'security', 'ssl_key_content']}
          >
            <Input.TextArea />
          </Form.Item>
        </div>
      </div>
      <Form.Item>
        <Button className="mr-4" onClick={prev}>
          {t('previous')}
        </Button>
        <Button type="primary" htmlType="submit">
          {t('next')}
        </Button>
      </Form.Item>
    </Form>
  )
}

export default TargetInfo
