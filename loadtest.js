import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 20000,
  duration: '5m',     // test duration
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

export default function () {
  const body = JSON.stringify({
    model: "gpt-3.5-turbo",
    messages: [{ role: "user", content: "Hello, world!" }],
    stream: false,    // streaming disabled
  });

  const headers = { "Content-Type": "application/json" };

  let res = http.post(`${BASE_URL}/v1/chat/completions`, body, { headers });

  check(res, { "status 200": (r) => r.status === 200 });
}

