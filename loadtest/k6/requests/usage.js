import http from 'k6/http';

export function request(data, config) {
  const url = `${config.baseUrl}/rewards/${data.id}/usage`;
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const body = {
    amount: data.amount,
  };

  return http.put(url, JSON.stringify(body), params);
}
