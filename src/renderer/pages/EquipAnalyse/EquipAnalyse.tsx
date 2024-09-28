import { Button, Image } from '@arco-design/web-react';

import { useMachine } from '@xstate/react';
import { createMachine } from 'xstate';
import { useEffect, useState } from 'react';
import { ParseNowPage } from '../../utils/api/localDetect';
import { HandlerAxiosErrPrefix } from '../../utils/api/help';
import { HeroTemplate, ParseResult } from '../../utils/const';
import { IntensifyAnalyse } from '../../utils/EquipAnalyse/IntensifyAnalyse/IntensifyAnalyse';
import { HeroTemplateList } from '../../utils/api/heroTemplate';

const ImageTypeBeiBao    = "beibao"
const ImageTypeIntensify = "intensify"
const toggleMachine = createMachine({
  id: 'toggle',
  initial: 'inactive',
  states: {
    inactive: {
      on: { TOGGLE: 'active' },
    },
    active: {
      on: { TOGGLE: 'inactive' },
    },
  },
});

// eslint-disable-next-line import/prefer-default-export
export function EquipAutoAnalyse() {
  const [timestamp, setTimestamp] = useState('');
  const [start, setStart] = useState(false);
  const [_, parseNowPage, refetch] = ParseNowPage();
  const [parseResult, setParseResult] = useState<ParseResult>();
  const [heroTemplateList, setHeroTemplateList] = useState<HeroTemplate[]>([])
  const [{ }, funcHeroTemplateList] = HeroTemplateList(false)

  useEffect(() => {
    funcHeroTemplateList().then((resp) => {
      setHeroTemplateList(resp.data.Data)
    }).catch((error) => {
      console.log(error)
      HandlerAxiosErrPrefix("读取角色模板", error)
    })
  }, []);

  console.log(parseResult);
  return (
    <div>
      <Button
        onClick={
          //() => setStart(!start)
          () => {
            parseNowPage()
              .then((resp) => {
                setParseResult(resp.data);
                return 200;
              })
              .catch((err) => {
                HandlerAxiosErrPrefix('读取角色模板', err);
              });
          }
        }
      >
        {start ? '停止检测' : '开始检测'}
      </Button>

      {parseResult?.Mode === ImageTypeBeiBao ? (
        <>背包</>
      ) : parseResult?.Mode === ImageTypeIntensify ? (
        <IntensifyAnalyse
          LevelUP={undefined}
          equip={parseResult?.Equips[0]}
          HeroTemplateList={heroTemplateList} />
      ) : (
        <>未检测到 背包/强化 界面</>
      )}
    </div>
  );
}
