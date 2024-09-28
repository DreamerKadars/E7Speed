import useAxios from 'axios-hooks';

export function ParseNowPage() {
  return useAxios(
    {
      baseURL: 'http://127.0.0.1:7766',
      url: 'api/v1/local/ParseNowPage',
      method: 'GET',
    },
    {
      manual: true,
    },
  );
}

export function AdbConnect() {
  return useAxios(
    {
      baseURL: 'http://127.0.0.1:7766',
      url: 'api/v1/local/AdbConnect',
      method: 'GET',
    },
    {
      manual: true,
    },
  );
}
