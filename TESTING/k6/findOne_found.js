import http from "k6/http";
import { check, sleep } from "k6";

const baseURL = "https://go-dev.manulife.com/api/v2";

export default function() {
  let res = http.get(baseURL + "/google");
  check(res, {
    "status was 200": (r) => r.status == 200,
    "transaction time OK": (r) => r.timings.duration < 200,
    "longURL was OK": (r) => r.json().longurl === "https://google.com"
  });
  sleep(1);
};
