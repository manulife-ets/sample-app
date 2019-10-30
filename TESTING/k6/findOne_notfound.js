import http from "k6/http";
import { check, sleep } from "k6";

const baseURL = "https://go-dev.manulife.com/api/v2";

export default function() {
  let res = http.get(baseURL + "/somethingReallyFake");
  check(res, {
    "status was 404": (r) => r.status == 404,
    "transaction time OK": (r) => r.timings.duration < 200,
    "longURL was always missing OK": (r) => r.json().longurl === undefined,
    "message was Invalid OK": (r) => r.json().message === "Invalid ID",
  });
  sleep(1);
};
