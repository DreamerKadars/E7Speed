import React, { useState } from 'react';
import { Layout } from '@arco-design/web-react';
import '@arco-design/web-react/dist/css/arco.css';
import Menu from './common/Menu/Menu';
import Header from './common/Header/Header';
import Footer from './common/Footer/Footer';
import { PathEquipAutoAnalyse, PathEquipAutoIntensify } from './const';
import { EquipAutoIntensify } from './pages/AutoIntensify/AutoIntensify';
import { EquipAutoAnalyse } from './pages/EquipAnalyse/EquipAnalyse';
// import icon from '../../assets/icon.svg';
// import './App.css';

// eslint-disable-next-line prefer-destructuring
const Sider = Layout.Sider;
// eslint-disable-next-line prefer-destructuring
const Content = Layout.Content;

export default function App() {
  const [RouteName, setRouteName] = useState(PathEquipAutoAnalyse);

  return (
    <div>
      <Header />
      <Layout style={{ width: 1000 }}>
        <Layout>
          {/* <Sider style={{ width: 160 }}>
            <div className="logo" />
            <Menu setNewKey={setRouteName} />
          </Sider> */}
          <Content style={{width: '96%', height: '100%', minHeight:560, background: 'rgb(240,255,255)', padding: '30px' }}>
            {RouteName === PathEquipAutoIntensify ? (
              <EquipAutoIntensify />
            ) : (
              <EquipAutoAnalyse />
            )}
          </Content>
        </Layout>
      </Layout>
      {/* <Footer /> */}
    </div>
  );
}
