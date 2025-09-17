import ws from "k6/ws";
import { check } from "k6";

export default function () {
  const url = "ws://localhost:8080/ws/motion";
  const res = ws.connect(url, {}, function (socket) {
    socket.on("open", () => {
      socket.send("hello");
    });
    socket.on("message", (msg) => {
      // console.log(msg); // optional
    });
    socket.setTimeout(() => {
      socket.close();
    }, 10000);
  });

  check(res, { "status is 101": (r) => r && r.status === 101 });
}
