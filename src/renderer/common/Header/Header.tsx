import React, { useState } from 'react';
import { Button, Grid, Input, InputNumber, Menu, Message } from '@arco-design/web-react';
import '@arco-design/web-react/dist/css/arco.css';
import { AdbConnect } from '../../utils/api/localDetect';
import { HandlerAxiosErr } from '../../const';
const MenuItem = Menu.Item;

function ConnectPortConfig() {
  const [connectPort, setConnectPort] = useState(16384);
  const [connectState, setConnectState] = useState('');

  const [{ data, loading, error }, adbConnect, refetch] = AdbConnect();

  return (
    <div>
      <Grid.Row>
        <Grid.Col span={16}>
          <InputNumber
            prefix={"模拟器端口"}
            value={connectPort}
            hideControl
            onChange={setConnectPort}
          ></InputNumber>
        </Grid.Col>
        <Grid.Col span={4}>
          {connectState === 'success' ? (
            <>
              <Button status="success">连接成功</Button>
            </>
          ) : (
            <>
              <Button
                type="primary"
                onClick={() => {
                   adbConnect({ params: { Port: connectPort } }).then(() => {
                     Message.success('连接模拟器成功');
                     setConnectState('success');
                   }).catch(err=>{
                    HandlerAxiosErr(err)
                   });
                }}
              >
                连接
              </Button>
            </>
          )}
        </Grid.Col>
      </Grid.Row>
    </div>
  );
}

export default function Header() {
  return (
    <div className="menu-demo">
      <Grid.Row>
        <Grid.Col span={18}>
          <Menu
            mode="horizontal"
            defaultSelectedKeys={['1']}
            className="menu-header"
          >
            <MenuItem key="0" style={{ padding: 0, marginRight: 38 }} disabled>
              <div
                style={{
                  width: 80,
                  height: 30,
                  borderRadius: 2,
                  background: 'var(--color-fill-3)',
                  cursor: 'text',
                }}
              />
            </MenuItem>
          </Menu>
        </Grid.Col>
        <Grid.Col span={6}>
          <ConnectPortConfig />
        </Grid.Col>
      </Grid.Row>
    </div>
  );
}
