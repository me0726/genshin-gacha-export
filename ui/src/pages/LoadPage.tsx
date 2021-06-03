/** @jsxImportSource @emotion/react */
import { css } from '@emotion/react';
import { FC, useCallback, useState } from 'react';
import { Button, Upload, Alert, Spin } from 'antd';
import InboxOutlined from '@ant-design/icons/InboxOutlined';
import { RcFile } from 'antd/lib/upload';
import { useGlobalContext } from 'context/GlobalContext';
import XLSXNameSpace from 'xlsx/types';
import { FriendLinks } from 'components/FriendLinks';
import parseExcel from 'utils/parseExcel';

const { Dragger } = Upload;
type LoadPageProps = {
  onLoad?: () => void;
};
// 预加载
// @ts-ignore
import('xlsx/dist/xlsx.mini.min.js');

export const LoadPage: FC<LoadPageProps> = function ({ onLoad }) {
  const [loading, setLoading] = useState<Boolean>(false);
  const [errorMessage, setErrorMessage] = useState<String | null>(null);
  const [loadingTip, setLoadingTip] = useState('加载中...');
  const { updateParsedData, updatePage } = useGlobalContext();

  function beginParse(blob: Blob) {
    const reader = new FileReader();
    reader.onload = function (e: ProgressEvent<FileReader>) {
      setLoadingTip('xlsx解析文件加载中...');
      // @ts-ignore
      import('xlsx/dist/xlsx.mini.min.js')
        .then((module) => {
          try {
            const XLSX: typeof XLSXNameSpace = module;
            const data = new Uint8Array((e.target as FileReader).result as any);
            const workbook = XLSX.read(data, { type: 'array' });
            updateParsedData(parseExcel(XLSX, workbook));
          } catch (e) {
            setErrorMessage(e.message);
          }
        })
        .catch(() => {
          setErrorMessage('XLSX解析文件加载失败，请重新上传');
        });
    };
    reader.onerror = function (e: ProgressEvent<FileReader>) {
      setErrorMessage('解析文件失败, 请重新上传');
      setLoading(false);
    };
    reader.readAsArrayBuffer(blob);
  }

  const onClick = async () => {
    try {
      setLoading(true);
      var blob = await (await fetch('/static/data.xlsx')).blob()
      if (blob.type == 'text/plain' || blob.type == 'text/html') {
        throw new Error(await blob.text())
      }
      beginParse(blob);
    } catch (e) {
      setErrorMessage('发生错误: '+String(e));
      setLoading(false);
    }
  }
  const handleUpload = useCallback((file: RcFile) => {
    if (!file.name.endsWith('.xlsx')) {
      setErrorMessage('文件类型错误，请上传xlsx文件');
      return false;
    }
    setLoading(true);
    beginParse(file)
    return false;
  }, []);
  return (
    <div
      css={css`
        overflow: auto;
      `}
    >
      <div
        className='ant-alert-info'
        css={css`
          width: 100%;
          display: flex;
          justify-content: center;
          border: none;
          margin: 20px 0;
        `}
      >
        <Alert
          message={
            <>
              <div>
                直接加载此电脑的抽卡记录
                <Button type='link' onClick={onClick}>
                  请点击这里
                </Button>
              </div>
            </>
          }
          type='info'
          showIcon={false}
          banner
        />
      </div>
      <div css={css`
          width: 100%;
          max-width: 600px;
          margin: 0 auto;
          padding: 40px;
        `}>
        {errorMessage && <Alert message={errorMessage} type='error' />}
      </div>
      <Dragger
        name='file'
        accept='application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, .xlsx'
        multiple={false}
        beforeUpload={handleUpload}
        showUploadList={false}
        css={css`
          width: 100%;
          max-width: 600px;
          margin: 0 auto;
          padding: 40px;
        `}
      >

        <p className='ant-upload-drag-icon'>
          {loading ? <Spin tip={loadingTip} /> : <InboxOutlined />}
        </p>
        <p className='ant-upload-text'>点击选择抽卡记录导出文件或将文件拖拽到此区域</p>
        <p className='ant-upload-text'>( 注：文件的后缀应为.xlsx )</p>
      </Dragger>
      <FriendLinks mode='bottom' visible />
    </div>
  );
};
