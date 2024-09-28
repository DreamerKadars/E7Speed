import React from 'react';
import { Menu as ArcoMenu } from '@arco-design/web-react';
import {
  IconHome,
  IconCalendar,
  IconBook,
  IconDesktop,
  IconUserGroup,
  IconBulb,
} from '@arco-design/web-react/icon';
import '@arco-design/web-react/dist/css/arco.css';
import { SkipToUrl } from '../../utils/helper';
import {
  PathAdmin,
  PathEquipAutoAnalyse,
  PathEquipAutoIntensify,
  PathHeroTemplateManage,
} from '../../const';
const MenuItem = ArcoMenu.Item;
// const SubMenu = ArcoMenu.SubMenu;

interface MenuProps {
  setNewKey: (str:string) => any;
}

export default function Menu(props: MenuProps) {
  return (
    <ArcoMenu
      selectedKeys={[window.location.href]}
      onClickMenuItem={(key) => {
        props.setNewKey(key);
      }}
      style={{ width: '100%', height: '100%', minHeight:650 }}
    >
      <MenuItem key={PathEquipAutoAnalyse} style={{}}>
        <IconBulb />
        装备分析
      </MenuItem>
      <MenuItem key={PathEquipAutoIntensify} style={{}}>
        <IconHome />
        自动强化
      </MenuItem>
    </ArcoMenu>
  );
}
